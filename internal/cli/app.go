package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"golang.org/x/sync/errgroup"

	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/context"
	"github.com/72636c/stratus/internal/stratus"
)

const (
	usageFormat = `usage: stratus [options] %[1]s

[options]
--file path%[2]cto%[2]cstratus.json|yaml (default .%[2]cstratus.yaml)
--output %[3]s (default plain)
`
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			usageFormat,
			commandNames,
			os.PathSeparator,
			loggerNames,
		)
	}
}

type App struct {
	cfg     *config.Config
	client  *stratus.Client
	command Command
	logger  Logger
}

func New() (_ *App, err error) {
	defer func() {
		if err != nil {
			flag.Usage()
		}
	}()

	cfgPath := flag.String("file", "stratus.yaml", "config file")
	loggerName := flag.String("output", "plain", "output format")

	flag.Parse()

	logger, ok := nameToLogger[*loggerName]
	if !ok {
		return nil, fmt.Errorf("output '%s' not recognised", *loggerName)
	}

	commandName := flag.Arg(0)

	command, ok := nameToCommand[commandName]
	if !ok {
		return nil, fmt.Errorf("command '%s' not recognised", commandName)
	}

	cfg, err := config.FromPath(*cfgPath)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}

	awsConfig := aws.NewConfig().WithHTTPClient(httpClient)

	provider, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	cloudFormation := cloudformation.New(provider)

	client := stratus.NewClient(cloudFormation)

	app := &App{
		cfg:     cfg,
		client:  client,
		command: command,
		logger:  logger,
	}

	return app, nil
}

func (app *App) Do(ctx context.Context) error {
	for index := 0; index < len(app.cfg.Stacks); index++ {
		err := app.do(ctx, index)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) do(ctx context.Context, index int) error {
	stack := app.cfg.Stacks[index]

	group, ctx := errgroup.WithContext(ctx)

	output := make(chan interface{})

	group.Go(func() (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = fmt.Errorf("recovered from panic: %+v", recovered)
			}
		}()

		defer close(output)

		output <- fmt.Sprintf("Stratus.[%d].StackConfig", index)

		output <- stack

		ctx = context.WithOutput(ctx, output)

		return app.command(ctx, app.client, stack)
	})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	group.Go(func() error {
		for {
			select {
			case <-ticker.C:
				fmt.Printf(".")

			case message, ok := <-output:
				app.logger(message, ok)

				if !ok {
					return nil
				}
			}
		}
	})

	return group.Wait()
}
