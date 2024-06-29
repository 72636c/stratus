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
	logger := context.Logger(ctx)

	logger.Title("Find existing change set")

	changeSet, err := client.FindExistingChangeSet(ctx, stack)

	if err != nil {
		logger.Title("Could not find existing change set. Creating new change set.")

		if _, changeSet, err = Stage(ctx, client, stack); err != nil {
			return err
		}
	}

	if changeSet != nil {
		logger.Title("Execute change set")

		err = client.ExecuteChangeSet(ctx, stack, *changeSet.ChangeSetName)
		if err != nil {
			return err
		}
	}

	logger.Title("Describe outputs")

	outputs, err := client.DescribeOutputs(ctx, stack)
	if err != nil {
		return err
	}

	logger.Data(outputs)

	logger.Title("Set stack policy")

	err = client.SetStackPolicy(ctx, stack)
	if err != nil {
		return err
	}

	logger.Title("Update termination protection")

	return client.UpdateTerminationProtection(ctx, stack)
}
