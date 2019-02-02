package command

import (
	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/context"
	"github.com/72636c/stratus/internal/stratus"
)

func Deploy(
	ctx context.Context,
	client *stratus.Client,
	stack *config.Stack,
	diff *stratus.Diff,
) error {
	output := context.Output(ctx)

	if diff.HasChangeSet() {
		output <- "Execute change set"

		err := client.ExecuteChangeSet(ctx, stack, *diff.ChangeSet.ChangeSetName)
		if err != nil {
			return err
		}
	}

	output <- "Set stack policy"

	err := client.SetStackPolicy(ctx, stack)
	if err != nil {
		return err
	}

	output <- "Update termination protection"

	return client.UpdateTerminationProtection(ctx, stack)
}
