package cli

import (
	"strings"

	"github.com/72636c/stratus/internal/command"
	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/context"
	"github.com/72636c/stratus/internal/stratus"
)

var (
	nameToCommand = map[string]Command{
		"delete": command.Delete,
		"deploy": command.Deploy,
		"stage":  stageAdapter,
	}

	commandNames = func() string {
		names := make([]string, 0)

		for name := range nameToCommand {
			names = append(names, name)
		}

		return strings.Join(names, "|")
	}()
)

type Command func(context.Context, *stratus.Client, *config.Stack) error

func stageAdapter(
	ctx context.Context,
	client *stratus.Client,
	stack *config.Stack,
) (err error) {
	_, err = command.Stage(ctx, client, stack)
	return
}
