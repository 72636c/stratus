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

	output <- "Find existing change set"

	changeSet, err := client.FindExistingChangeSet(ctx, stack)
	if err != nil {
		return err
	}

	if changeSet != nil {
		output <- "Execute change set"

		err = client.ExecuteChangeSet(ctx, stack, *changeSet.ChangeSetName)
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
