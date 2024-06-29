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

func Test_Deploy_Happy_ExistingChangeSet_Create(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,

		Capabilities: []string{
			cloudformation.CapabilityCapabilityAutoExpand,
			cloudformation.CapabilityCapabilityNamedIam,
		},
		Parameters: config.StackParameters{
			&config.StackParameter{
				Key:   "test-parameter-key-b",
				Value: "test-parameter-value-b",
			},
			&config.StackParameter{
				Key:   "test-parameter-key-a",
				Value: "test-parameter-value-a",
			},
		},
		TerminationProtection: true,

		Policy:   []byte(mockStackPolicy),
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
			"GetTemplateWithContext",
			&cloudformation.GetTemplateInput{
				ChangeSetName: aws.String(mockChangeSetCreateName),
				StackName:     aws.String(stack.Name),
				TemplateStage: aws.String(cloudformation.TemplateStageOriginal),
			},
		).
		Return(
			&cloudformation.GetTemplateOutput{
				TemplateBody: aws.String(mockStackTemplate),
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
		Return(
			&cloudformation.DescribeChangeSetOutput{
				ChangeSetName: aws.String(mockChangeSetCreateName),
				Capabilities: []*string{
					aws.String(cloudformation.CapabilityCapabilityAutoExpand),
					aws.String(cloudformation.CapabilityCapabilityNamedIam),
				},
				Parameters: []*cloudformation.Parameter{
					&cloudformation.Parameter{
						ParameterKey:   aws.String("test-parameter-key-a"),
						ParameterValue: aws.String("test-parameter-value-a"),
					},
					&cloudformation.Parameter{
						ParameterKey:   aws.String("test-parameter-key-b"),
						ParameterValue: aws.String("test-parameter-value-b"),
					},
				},
			},
			nil,
		).
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
			"DescribeStackEventsWithContext",
			&cloudformation.DescribeStackEventsInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.DescribeStackEventsOutput{
				StackEvents: make([]*cloudformation.StackEvent, 0),
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
						Outputs: make([]*cloudformation.Output, 0),
					},
				},
			},
			nil,
		).
		On(
			"SetStackPolicyWithContext",
			&cloudformation.SetStackPolicyInput{
				StackName:       aws.String(mockStackName),
				StackPolicyBody: aws.String(mockStackPolicy),
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

	client := stratus.NewClient(cfn, nil)

	err := command.Deploy(context.Background(), client, stack)
	assert.NoError(err)
}

func Test_Deploy_Happy_ExistingChangeSet_Update(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,

		Capabilities:          make([]string, 0),
		Parameters:            make(config.StackParameters, 0),
		TerminationProtection: true,

		Policy:   []byte(mockStackPolicy),
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
		Return(
			&cloudformation.ListChangeSetsOutput{
				Summaries: []*cloudformation.ChangeSetSummary{
					&cloudformation.ChangeSetSummary{
						ChangeSetName:   aws.String(mockChangeSetUpdateName),
						ExecutionStatus: aws.String(cloudformation.ExecutionStatusAvailable),
					},
				},
			},
			nil,
		).
		On(
			"GetTemplateWithContext",
			&cloudformation.GetTemplateInput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				StackName:     aws.String(stack.Name),
				TemplateStage: aws.String(cloudformation.TemplateStageOriginal),
			},
		).
		Return(
			&cloudformation.GetTemplateOutput{
				TemplateBody: aws.String(mockStackTemplate),
			},
			nil,
		).
		On(
			"DescribeChangeSetWithContext",
			&cloudformation.DescribeChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				StackName:     aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.DescribeChangeSetOutput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				Capabilities:  make([]*string, 0),
				Parameters:    make([]*cloudformation.Parameter, 0),
			},
			nil,
		).
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
			"DescribeStackEventsWithContext",
			&cloudformation.DescribeStackEventsInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.DescribeStackEventsOutput{
				StackEvents: make([]*cloudformation.StackEvent, 0),
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
						Outputs: make([]*cloudformation.Output, 0),
					},
				},
			},
			nil,
		).
		On(
			"SetStackPolicyWithContext",
			&cloudformation.SetStackPolicyInput{
				StackName:       aws.String(mockStackName),
				StackPolicyBody: aws.String(mockStackPolicy),
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

	client := stratus.NewClient(cfn, nil)

	err := command.Deploy(context.Background(), client, stack)
	assert.NoError(err)
}

func Test_Deploy_Happy_NoopChangeSet(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,

		Capabilities:          make([]string, 0),
		Parameters:            make(config.StackParameters, 0),
		TerminationProtection: true,

		Policy:   []byte(mockStackPolicy),
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
		Return(
			&cloudformation.ListChangeSetsOutput{
				Summaries: []*cloudformation.ChangeSetSummary{
					&cloudformation.ChangeSetSummary{
						ChangeSetName:   aws.String(mockChangeSetUpdateName),
						ExecutionStatus: aws.String(cloudformation.ExecutionStatusUnavailable),
						Status:          aws.String(cloudformation.ChangeSetStatusFailed),
						StatusReason:    aws.String("The submitted information didn't contain changes. Submit different information to create a change set."),
					},
				},
			},
			nil,
		).
		On(
			"GetTemplateWithContext",
			&cloudformation.GetTemplateInput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				StackName:     aws.String(stack.Name),
				TemplateStage: aws.String(cloudformation.TemplateStageOriginal),
			},
		).
		Return(
			&cloudformation.GetTemplateOutput{
				TemplateBody: aws.String(mockStackTemplate),
			},
			nil,
		).
		On(
			"DescribeChangeSetWithContext",
			&cloudformation.DescribeChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				StackName:     aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.DescribeChangeSetOutput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				Capabilities:  make([]*string, 0),
				Parameters:    make([]*cloudformation.Parameter, 0),
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
						Outputs: make([]*cloudformation.Output, 0),
					},
				},
			},
			nil,
		).
		On(
			"SetStackPolicyWithContext",
			&cloudformation.SetStackPolicyInput{
				StackName:       aws.String(mockStackName),
				StackPolicyBody: aws.String(mockStackPolicy),
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

	client := stratus.NewClient(cfn, nil)

	err := command.Deploy(context.Background(), client, stack)
	assert.NoError(err)
}

func Test_Deploy_Happy_NoopChangeSet_UploadArtefacts(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,

		Capabilities:          make([]string, 0),
		Parameters:            make(config.StackParameters, 0),
		TerminationProtection: true,

		Policy:   []byte(mockStackPolicy),
		Template: []byte(mockStackTemplate),

		ArtefactBucket: mockArtefactBucket,
		PolicyKey:      mockStackPolicyKey,

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
						ChangeSetName:   aws.String(mockChangeSetUpdateName),
						ExecutionStatus: aws.String(cloudformation.ExecutionStatusUnavailable),
						Status:          aws.String(cloudformation.ChangeSetStatusFailed),
						StatusReason:    aws.String("The submitted information didn't contain changes. Submit different information to create a change set."),
					},
				},
			},
			nil,
		).
		On(
			"GetTemplateWithContext",
			&cloudformation.GetTemplateInput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				StackName:     aws.String(stack.Name),
				TemplateStage: aws.String(cloudformation.TemplateStageOriginal),
			},
		).
		Return(
			&cloudformation.GetTemplateOutput{
				TemplateBody: aws.String(mockStackTemplate),
			},
			nil,
		).
		On(
			"DescribeChangeSetWithContext",
			&cloudformation.DescribeChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				StackName:     aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.DescribeChangeSetOutput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				Capabilities:  make([]*string, 0),
				Parameters:    make([]*cloudformation.Parameter, 0),
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
						Outputs: make([]*cloudformation.Output, 0),
					},
				},
			},
			nil,
		).
		On(
			"SetStackPolicyWithContext",
			&cloudformation.SetStackPolicyInput{
				StackName:      aws.String(mockStackName),
				StackPolicyURL: aws.String(mockStackPolicyURL),
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

	client := stratus.NewClient(cfn, nil)

	err := command.Deploy(context.Background(), client, stack)
	assert.NoError(err)
}

func Test_Deploy_Happy_ImplicitStage(t *testing.T) {
	assert := assert.New(t)

	stack := &config.Stack{
		Name: mockStackName,

		Capabilities:          make([]string, 0),
		Parameters:            make(config.StackParameters, 0),
		TerminationProtection: true,

		Policy:   []byte(mockStackPolicy),
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
		Return(
			&cloudformation.ListChangeSetsOutput{
				Summaries: []*cloudformation.ChangeSetSummary{},
			},
			nil,
		).Once().
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
		Return(
			&cloudformation.DescribeChangeSetOutput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				Capabilities:  make([]*string, 0),
				Parameters:    make([]*cloudformation.Parameter, 0),
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
					{
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
		Return(nil, nil).
		On(
			"DescribeStackEventsWithContext",
			&cloudformation.DescribeStackEventsInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.DescribeStackEventsOutput{
				StackEvents: make([]*cloudformation.StackEvent, 0),
			},
			nil,
		).
		On(
			"ExecuteChangeSetWithContext",
			&cloudformation.ExecuteChangeSetInput{
				ChangeSetName: aws.String(mockChangeSetUpdateName),
				StackName:     aws.String(mockStackName),
			},
		).
		Return(nil, nil).
		On(
			"DescribeStackEventsWithContext",
			&cloudformation.DescribeStackEventsInput{
				StackName: aws.String(stack.Name),
			},
		).
		Return(
			&cloudformation.DescribeStackEventsOutput{
				StackEvents: make([]*cloudformation.StackEvent, 0),
			},
			nil,
		).
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
				StackPolicyBody: aws.String(mockStackPolicy),
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

	client := stratus.NewClient(cfn, nil)

	err := command.Deploy(context.Background(), client, stack)
	assert.NoError(err)
}
