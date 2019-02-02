package command_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/stretchr/testify/assert"

	"github.com/72636c/stratus/internal/command"
	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/context"
	"github.com/72636c/stratus/internal/stratus"
)

func Test_Preview_Happy_ExistingChangeSet(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,

		Checksum: mockChecksum,
	}

	cfn := stratus.NewCloudFormationMock()
	defer cfn.AssertExpectations(t)
	cfn.
		On(
			"ListChangeSetsWithContext",
			&cloudformation.ListChangeSetsInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.ListChangeSetsOutput{
				Summaries: []*cloudformation.ChangeSetSummary{
					&cloudformation.ChangeSetSummary{
						ChangeSetName:   aws.String(mockChangeSetCreateName),
						ExecutionStatus: aws.String(cloudformation.ExecutionStatusAvailable),
					},
				},
			},
			nil,
		).
		On(
			"DescribeChangeSetWithContext",
			&cloudformation.DescribeChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetCreateName),
				StackName:     aws.String(stack.Name),
			},
		).
		Return(new(cloudformation.DescribeChangeSetOutput), nil).
		On(
			"DescribeStacksWithContext",
			&cloudformation.DescribeStacksInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.DescribeStacksOutput{
				Stacks: []*cloudformation.Stack{
					&cloudformation.Stack{
						EnableTerminationProtection: aws.Bool(false),
					},
				},
			},
			nil,
		).
		On(
			"GetStackPolicyWithContext",
			&cloudformation.GetStackPolicyInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(nil, nil)

	client := stratus.NewClient(cfn)

	_, err := command.Preview(context.Background(), client, stack)
	assert.NoError(err)
}

func Test_Preview_Happy_NewChangeSet_Create(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,

		Template: []byte(mockStackTemplate),

		Checksum: mockChecksum,
	}

	cfn := stratus.NewCloudFormationMock()
	defer cfn.AssertExpectations(t)
	cfn.
		On(
			"ListChangeSetsWithContext",
			&cloudformation.ListChangeSetsInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(nil, awserr.New("ValidationError", "does not exist", nil)).
		On(
			"ValidateTemplateWithContext",
			&cloudformation.ValidateTemplateInput{
				TemplateBody: aws.String(string(stack.Template)),
			},
		).
		Return(nil, nil).
		On(
			"CreateChangeSetWithContext",
			&cloudformation.CreateChangeSetInput{
				Capabilities:        make([]*string, 0),
				ChangeSetName:       aws.String(mockChangeSetUpdateName),
				ChangeSetType:       aws.String(cloudformation.ChangeSetTypeUpdate),
				StackName:           aws.String(stack.Name),
				Parameters:          make([]*cloudformation.Parameter, 0),
				Tags:                make([]*cloudformation.Tag, 0),
				TemplateBody:        aws.String(string(stack.Template)),
				UsePreviousTemplate: aws.Bool(false),
			},
		).
		Return(nil, awserr.New("ValidationError", "does not exist", nil)).
		On(
			"CreateChangeSetWithContext",
			&cloudformation.CreateChangeSetInput{
				Capabilities:        make([]*string, 0),
				ChangeSetName:       aws.String(mockChangeSetCreateName),
				ChangeSetType:       aws.String(cloudformation.ChangeSetTypeCreate),
				StackName:           aws.String(stack.Name),
				Parameters:          make([]*cloudformation.Parameter, 0),
				Tags:                make([]*cloudformation.Tag, 0),
				TemplateBody:        aws.String(string(stack.Template)),
				UsePreviousTemplate: aws.Bool(false),
			},
		).
		Return(nil, nil).
		On(
			"WaitUntilChangeSetCreateCompleteWithContext",
			&cloudformation.DescribeChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetCreateName),
				StackName:     aws.String(stack.Name),
			},
		).
		Return(nil).
		On(
			"DescribeChangeSetWithContext",
			&cloudformation.DescribeChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetCreateName),
				StackName:     aws.String(stack.Name),
			},
		).
		Return(new(cloudformation.DescribeChangeSetOutput), nil).
		On(
			"DescribeStacksWithContext",
			&cloudformation.DescribeStacksInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.DescribeStacksOutput{
				Stacks: []*cloudformation.Stack{
					&cloudformation.Stack{
						EnableTerminationProtection: aws.Bool(false),
					},
				},
			},
			nil,
		).
		On(
			"GetStackPolicyWithContext",
			&cloudformation.GetStackPolicyInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(nil, nil)

	client := stratus.NewClient(cfn)

	_, err := command.Preview(context.Background(), client, stack)
	assert.NoError(err)
}

func Test_Preview_Happy_NewChangeSet_Update(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,

		Template: []byte(mockStackTemplate),

		Checksum: mockChecksum,
	}

	cfn := stratus.NewCloudFormationMock()
	defer cfn.AssertExpectations(t)
	cfn.
		On(
			"ListChangeSetsWithContext",
			&cloudformation.ListChangeSetsInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(nil, awserr.New("ValidationError", "does not exist", nil)).
		On(
			"ValidateTemplateWithContext",
			&cloudformation.ValidateTemplateInput{
				TemplateBody: aws.String(string(stack.Template)),
			},
		).
		Return(nil, nil).
		On(
			"CreateChangeSetWithContext",
			&cloudformation.CreateChangeSetInput{
				Capabilities:        make([]*string, 0),
				ChangeSetName:       aws.String(mockChangeSetUpdateName),
				ChangeSetType:       aws.String(cloudformation.ChangeSetTypeUpdate),
				StackName:           aws.String(stack.Name),
				Parameters:          make([]*cloudformation.Parameter, 0),
				Tags:                make([]*cloudformation.Tag, 0),
				TemplateBody:        aws.String(string(stack.Template)),
				UsePreviousTemplate: aws.Bool(false),
			},
		).
		Return(nil, nil).
		On(
			"WaitUntilChangeSetCreateCompleteWithContext",
			&cloudformation.DescribeChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				StackName:     aws.String(stack.Name),
			},
		).
		Return(nil).
		On(
			"DescribeChangeSetWithContext",
			&cloudformation.DescribeChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				StackName:     aws.String(stack.Name),
			},
		).
		Return(new(cloudformation.DescribeChangeSetOutput), nil).
		On(
			"DescribeStacksWithContext",
			&cloudformation.DescribeStacksInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.DescribeStacksOutput{
				Stacks: []*cloudformation.Stack{
					&cloudformation.Stack{
						EnableTerminationProtection: aws.Bool(false),
					},
				},
			},
			nil,
		).
		On(
			"GetStackPolicyWithContext",
			&cloudformation.GetStackPolicyInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(nil, nil)

	client := stratus.NewClient(cfn)

	_, err := command.Preview(context.Background(), client, stack)
	assert.NoError(err)
}

func Test_Preview_Happy_NoopChangeSet(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,

		Template: []byte(mockStackTemplate),

		Checksum: mockChecksum,
	}

	cfn := stratus.NewCloudFormationMock()
	defer cfn.AssertExpectations(t)
	cfn.
		On(
			"ListChangeSetsWithContext",
			&cloudformation.ListChangeSetsInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(nil, awserr.New("ValidationError", "does not exist", nil)).
		On(
			"ValidateTemplateWithContext",
			&cloudformation.ValidateTemplateInput{
				TemplateBody: aws.String(string(stack.Template)),
			},
		).
		Return(nil, nil).
		On(
			"CreateChangeSetWithContext",
			&cloudformation.CreateChangeSetInput{
				Capabilities:        make([]*string, 0),
				ChangeSetName:       aws.String(mockChangeSetUpdateName),
				ChangeSetType:       aws.String(cloudformation.ChangeSetTypeUpdate),
				StackName:           aws.String(stack.Name),
				Parameters:          make([]*cloudformation.Parameter, 0),
				Tags:                make([]*cloudformation.Tag, 0),
				TemplateBody:        aws.String(string(stack.Template)),
				UsePreviousTemplate: aws.Bool(false),
			},
		).
		Return(nil, nil).
		On(
			"WaitUntilChangeSetCreateCompleteWithContext",
			&cloudformation.DescribeChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				StackName:     aws.String(stack.Name),
			},
		).
		Return(awserr.New(request.WaiterResourceNotReadyErrorCode, "", nil)).
		On(
			"DescribeChangeSetWithContext",
			&cloudformation.DescribeChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				StackName:     aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.DescribeChangeSetOutput{
				Status:       aws.String(cloudformation.ChangeSetStatusFailed),
				StatusReason: aws.String("The submitted information didn't contain changes. Submit different information to create a change set."),
			},
			nil,
		).
		On(
			"DescribeStacksWithContext",
			&cloudformation.DescribeStacksInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.DescribeStacksOutput{
				Stacks: []*cloudformation.Stack{
					&cloudformation.Stack{
						EnableTerminationProtection: aws.Bool(false),
					},
				},
			},
			nil,
		).
		On(
			"GetStackPolicyWithContext",
			&cloudformation.GetStackPolicyInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(nil, nil)

	client := stratus.NewClient(cfn)

	_, err := command.Preview(context.Background(), client, stack)
	assert.NoError(err)
}
