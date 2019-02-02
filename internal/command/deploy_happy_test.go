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

func Test_Deploy_Happy_ChangeSet_Create(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,

		TerminationProtection: true,

		Policy: interface{}(mockStackPolicy),

		Checksum: mockChecksum,
	}

	diff := &stratus.Diff{
		ChangeSet: &cloudformation.DescribeChangeSetOutput{
			ChangeSetName: aws.String(mockChangeSetCreateName),
		},
	}

	cfn := stratus.NewCloudFormationMock()
	defer cfn.AssertExpectations(t)
	cfn.
		On(
			"ExecuteChangeSetWithContext",
			&cloudformation.ExecuteChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetCreateName),
				StackName:     aws.String(mockStackName),
			},
		).
		Return(nil, nil).
		On(
			"WaitUntilStackCreateCompleteWithContext",
			&cloudformation.DescribeStacksInput{
				StackName: aws.String(mockStackName),
			},
		).
		Return(nil).
		On(
			"SetStackPolicyWithContext",
			&cloudformation.SetStackPolicyInput{
				StackName:       aws.String(mockStackName),
				StackPolicyBody: aws.String(string(mockStackPolicyBody)),
			},
		).
		Return(nil, nil).
		On(
			"UpdateTerminationProtectionWithContext",
			&cloudformation.UpdateTerminationProtectionInput{
				EnableTerminationProtection: aws.Bool(true),
				StackName:                   aws.String(mockStackName),
			},
		).
		Return(nil, nil)

	client := stratus.NewClient(cfn)

	err := command.Deploy(context.Background(), client, stack, diff)
	assert.NoError(err)
}

func Test_Deploy_Happy_ChangeSet_Update(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,

		TerminationProtection: true,

		Policy: interface{}(mockStackPolicy),

		Checksum: mockChecksum,
	}

	diff := &stratus.Diff{
		ChangeSet: &cloudformation.DescribeChangeSetOutput{
			ChangeSetName: aws.String(mockChangeSetUpdateName),
		},
	}

	cfn := stratus.NewCloudFormationMock()
	defer cfn.AssertExpectations(t)
	cfn.
		On(
			"ExecuteChangeSetWithContext",
			&cloudformation.ExecuteChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				StackName:     aws.String(mockStackName),
			},
		).
		Return(nil, nil).
		On(
			"WaitUntilStackUpdateCompleteWithContext",
			&cloudformation.DescribeStacksInput{
				StackName: aws.String(mockStackName),
			},
		).
		Return(nil).
		On(
			"SetStackPolicyWithContext",
			&cloudformation.SetStackPolicyInput{
				StackName:       aws.String(mockStackName),
				StackPolicyBody: aws.String(string(mockStackPolicyBody)),
			},
		).
		Return(nil, nil).
		On(
			"UpdateTerminationProtectionWithContext",
			&cloudformation.UpdateTerminationProtectionInput{
				EnableTerminationProtection: aws.Bool(true),
				StackName:                   aws.String(mockStackName),
			},
		).
		Return(nil, nil)

	client := stratus.NewClient(cfn)

	err := command.Deploy(context.Background(), client, stack, diff)
	assert.NoError(err)
}

func Test_Deploy_Happy_NoChangeSet(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,

		TerminationProtection: true,

		Policy: interface{}(mockStackPolicy),

		Checksum: mockChecksum,
	}

	diff := new(stratus.Diff)

	cfn := stratus.NewCloudFormationMock()
	defer cfn.AssertExpectations(t)
	cfn.
		On(
			"SetStackPolicyWithContext",
			&cloudformation.SetStackPolicyInput{
				StackName:       aws.String(mockStackName),
				StackPolicyBody: aws.String(string(mockStackPolicyBody)),
			},
		).
		Return(nil, nil).
		On(
			"UpdateTerminationProtectionWithContext",
			&cloudformation.UpdateTerminationProtectionInput{
				EnableTerminationProtection: aws.Bool(true),
				StackName:                   aws.String(mockStackName),
			},
		).
		Return(nil, nil)

	client := stratus.NewClient(cfn)

	err := command.Deploy(context.Background(), client, stack, diff)
	assert.NoError(err)
}
