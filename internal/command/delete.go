package command

import (
	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/context"
	"github.com/72636c/stratus/internal/stratus"
)

func Delete(
	ctx context.Context,
	client *stratus.Client,
	stack *config.Stack,
) error {
	output := context.Output(ctx)

	output <- "Delete stack"

	return client.DeleteStack(ctx, stack)
}
