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
) error {
	output := context.Output(ctx)

	diff, err := Preview(ctx, client, stack)
	if err != nil {
		return err
	}

	if diff.HasChangeSet() {
		output <- "Execute change set"

		name := *diff.ChangeSet.ChangeSetName

		err := client.ExecuteChangeSet(ctx, stack, name)
		if err != nil {
			return err
		}
	}

	output <- "Set stack policy"

	err = client.SetStackPolicy(ctx, stack)
	if err != nil {
		return err
	}

	output <- "Update termination protection"

	return client.UpdateTerminationProtection(ctx, stack)
}
