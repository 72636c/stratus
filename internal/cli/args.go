package cli

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/stratus"
)

var (
	usage = fmt.Sprintf(
		"usage: stratus [-output %s] %s [stratus.json]",
		loggerNames,
		commandNames,
	)
)

type Args struct {
	cfg     *config.Config
	client  *stratus.Client
	command Command
	logger  Logger
}

func FromCommandLine() (*Args, error) {
	if len(os.Args) < 2 {
		return nil, errors.New(usage)
	}

	output := flag.String("output", "plain", "output format")

	flag.Parse()

	logger, ok := loggers[*output]
	if !ok {
		return nil, fmt.Errorf("output '%s' not recognised", *output)
	}

	commandString := flag.Arg(0)

	cfgPath := flag.Arg(1)
	if cfgPath == "" {
		cfgPath = "stratus.json"
	}

	command, ok := commands[commandString]
	if !ok {
		return nil, fmt.Errorf("command '%s' not recognised", commandString)
	}

	cfg, err := config.FromPath(cfgPath)
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

	args := &Args{
		cfg:     cfg,
		client:  client,
		command: command,
		logger:  logger,
	}

	return args, nil
}
