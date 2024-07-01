package command

import (
	"fmt"
	"os"

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
		return err
	}

	if changeSet == nil && os.Getenv("FORCE_DEPLOY") == "true" {
		logger.Title("Could not find existing change set. FORCE_DEPLOY is true, so creating a new change set.")

		_, changeSet, err = Stage(ctx, client, stack)
		if err != nil {
			return err
		}
	}

	if changeSet == nil {
		logger.Title("Could not find existing change set, exiting. To force a deployment with a new change set, set FORCE_DEPLOY=true and retry.")
		return fmt.Errorf("could not find existing change set")
	}

	if stratus.IsNoopChangeSet(changeSet) {
		logger.Title("No changes to execute.")
	} else {
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
