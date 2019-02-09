package command_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/stretchr/testify/assert"

	"github.com/72636c/stratus/internal/command"
	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/context"
	"github.com/72636c/stratus/internal/stratus"
)

func Test_Delete_Happy(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,
	}

	cfn := stratus.NewCloudFormationMock()
	defer cfn.AssertExpectations(t)
	cfn.
		On(
			"DeleteStackWithContext",
			&cloudformation.DeleteStackInput{
				StackName: aws.String(mockStackName),
			},
		).
		Return(nil, nil).
		On(
			"WaitUntilStackDeleteCompleteWithContext",
			&cloudformation.DescribeStacksInput{
				StackName: aws.String(mockStackName),
			},
		).
		Return(nil)

	client := stratus.NewClient(cfn, nil)

	err := command.Delete(context.Background(), client, stack)
	assert.NoError(err)

}
