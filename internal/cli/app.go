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
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/context"
	"github.com/72636c/stratus/internal/log"
	"github.com/72636c/stratus/internal/stratus"
)

const (
	usageFormat = `usage: stratus [options] %[1]s

[options]
--file path%[2]cto%[2]cstratus.json|yaml (default .%[2]cstratus.yaml)
--name select specific stack (default select all stacks)
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
	cfg       *config.Config
	command   Command
	logger    log.Logger
	stackName string

	newClient func(region *string) *stratus.Client
}

func New() (_ *App, err error) {
	defer func() {
		if err != nil {
			flag.Usage()
		}
	}()

	cfgPath := flag.String("file", "stratus.yaml", "config file")
	rawStackName := flag.String("name", "", "stack name")
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

	httpClient := &http.Client{Timeout: 10 * time.Second}

	httpConfig := aws.NewConfig().WithHTTPClient(httpClient)

	provider, err := session.NewSession(httpConfig)
	if err != nil {
		return nil, err
	}

	// TODO: memoise
	newClient := func(region *string) *stratus.Client {
		regionConfig := aws.NewConfig()

		if region != nil {
			regionConfig = regionConfig.WithRegion(*region)
		}

		cfnClient := cloudformation.New(provider, regionConfig)
		s3Client := s3.New(provider, regionConfig)

		return stratus.NewClient(cfnClient, s3Client)
	}

	config.Init(provider)

	cfg, err := config.FromPath(*cfgPath)
	if err != nil {
		return nil, err
	}

	stackName, err := config.Resolve(*rawStackName)
	if err != nil {
		return nil, err
	}

	app := &App{
		cfg:       cfg,
		command:   command,
		logger:    logger,
		stackName: stackName,

		newClient: newClient,
	}

	return app, nil
}

func (app *App) Do(ctx context.Context) error {
	ctx = context.WithLogger(ctx, app.logger)

	if app.stackName == "" {
		return app.doAll(ctx)
	}

	stack, ok := app.cfg.Stacks.Find(app.stackName)
	if !ok {
		return fmt.Errorf("stack '%s' not found in config", app.stackName)
	}

	app.logger.Title("Load config")
	app.logger.Data(stack)

	client := app.newClient(stack.Region)

	return app.command(context.WithLogger(ctx, app.logger), client, stack)
}

func (app *App) doAll(ctx context.Context) error {
	for index := 0; index < len(app.cfg.Stacks); index++ {
		stack := app.cfg.Stacks[index]

		app.logger.Title("Load config %d", index)
		app.logger.Data(stack)

		client := app.newClient(stack.Region)

		err := app.command(context.WithLogger(ctx, app.logger), client, stack)
		if err != nil {
			return err
		}
	}

	return nil
}
