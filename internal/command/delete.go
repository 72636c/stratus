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
	logger := context.Logger(ctx)

	logger.Title("Delete stack")

	return client.DeleteStack(ctx, stack)
}
