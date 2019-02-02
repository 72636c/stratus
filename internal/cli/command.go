package cli

import (
	"strings"

	"github.com/72636c/stratus/internal/command"
	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/context"
	"github.com/72636c/stratus/internal/stratus"
)

var (
	commands = map[string]Command{
		"delete":  command.Delete,
		"deploy":  command.Deploy,
		"preview": previewAdapter,
	}

	commandNames = func() string {
		names := make([]string, 0)

		for name := range commands {
			names = append(names, name)
		}

		return strings.Join(names, "|")
	}()
)

type Command func(context.Context, *stratus.Client, *config.Stack) error

func previewAdapter(
	ctx context.Context,
	client *stratus.Client,
	stack *config.Stack,
) (err error) {
	_, err = command.Preview(ctx, client, stack)
	return
}
