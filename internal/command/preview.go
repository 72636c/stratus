package command

import (
	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/context"
	"github.com/72636c/stratus/internal/stratus"
)

func Preview(
	ctx context.Context,
	client *stratus.Client,
	stack *config.Stack,
) (*stratus.Diff, error) {
	output := context.Output(ctx)

	output <- "Validate template"

	validateOutput, err := client.ValidateTemplate(ctx, stack)
	if err != nil {
		return nil, err
	}

	output <- validateOutput

	if stack.ShouldUpload() {
		output <- "Upload artefacts"

		err = client.UploadArtefacts(ctx, stack)
		if err != nil {
			return nil, err
		}
	}

	output <- "Create change set"

	describeOutput, err := client.CreateChangeSet(ctx, stack)
	if err != nil {
		return nil, err
	}

	output <- "Diff stack"

	diffOutput, err := client.Diff(ctx, stack, describeOutput)
	if err != nil {
		return nil, err
	}

	output <- diffOutput

	return diffOutput, nil
}
