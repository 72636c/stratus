package command

import (
	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/context"
	"github.com/72636c/stratus/internal/stratus"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func Stage(
	ctx context.Context,
	client *stratus.Client,
	stack *config.Stack,
) (*stratus.Diff, *cloudformation.DescribeChangeSetOutput, error) {
	logger := context.Logger(ctx)

	logger.Title("Validate template")

	validateOutput, err := client.ValidateTemplate(ctx, stack)
	if err != nil {
		return nil, nil, err
	}

	logger.Data(validateOutput)

	if stack.ShouldUpload() {
		logger.Title("Upload artefacts")

		err = client.UploadArtefacts(ctx, stack)
		if err != nil {
			return nil, nil, err
		}
	}

	logger.Title("Create change set")

	describeOutput, err := client.CreateChangeSet(ctx, stack)
	if err != nil {
		return nil, nil, err
	}

	logger.Title("Diff stack")

	diffOutput, err := client.Diff(ctx, stack, describeOutput)
	if err != nil {
		return nil, nil, err
	}

	logger.Data(diffOutput)

	return diffOutput, describeOutput, nil
}
