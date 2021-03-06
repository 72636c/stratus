package stratus_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/stretchr/testify/assert"

	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/stratus"
)

func stringPointer(value string) *string {
	return &value
}

func Test_MatchesChangeSetSummary(t *testing.T) {
	expectedChecksum := "1000000000200000000030000000004000000000500000000060000000007000"
	unexpectedChecksum := "1000000000200000000030000000004000000000500000000060000000007001"

	stack := &config.Stack{
		Checksum: expectedChecksum,
	}

	testCases := []struct {
		description string
		summary     *cloudformation.ChangeSetSummary
		expected    bool
	}{
		{
			description: "change set name for create",
			summary: &cloudformation.ChangeSetSummary{
				ChangeSetName:   aws.String(fmt.Sprintf("stratus-create-%s", expectedChecksum)),
				ExecutionStatus: aws.String(cloudformation.ExecutionStatusAvailable),
			},
			expected: true,
		},
		{
			description: "change set name for update",
			summary: &cloudformation.ChangeSetSummary{
				ChangeSetName:   aws.String(fmt.Sprintf("stratus-update-%s", expectedChecksum)),
				ExecutionStatus: aws.String(cloudformation.ExecutionStatusAvailable),
			},
			expected: true,
		},
		{
			description: "change set ID",
			summary: &cloudformation.ChangeSetSummary{
				ChangeSetName:   aws.String(fmt.Sprintf("arn:aws:cloudformation:ap-southeast-2:000000000000:changeSet/stratus-create-%s/00000000-0000-4000-8000-000000000000", expectedChecksum)),
				ExecutionStatus: aws.String(cloudformation.ExecutionStatusAvailable),
			},
			expected: true,
		},
		{
			description: "noop change set",
			summary: &cloudformation.ChangeSetSummary{
				ChangeSetName:   aws.String(fmt.Sprintf("stratus-create-%s", expectedChecksum)),
				ExecutionStatus: aws.String(cloudformation.ExecutionStatusUnavailable),
				Status:          aws.String(cloudformation.ChangeSetStatusFailed),
				StatusReason:    aws.String("The submitted information didn't contain changes. Submit different information to create a change set."),
			},
			expected: true,
		},
		{
			description: "unexpected execution status",
			summary: &cloudformation.ChangeSetSummary{
				ChangeSetName:   aws.String(fmt.Sprintf("stratus-create-%s", expectedChecksum)),
				ExecutionStatus: aws.String(cloudformation.ExecutionStatusObsolete),
			},
			expected: false,
		},
		{
			description: "non-matching checksum",
			summary: &cloudformation.ChangeSetSummary{
				ChangeSetName:   aws.String(fmt.Sprintf("stratus-create-%s", unexpectedChecksum)),
				ExecutionStatus: aws.String(cloudformation.ExecutionStatusAvailable),
			},
			expected: false,
		},
		{
			description: "capitalised type",
			summary: &cloudformation.ChangeSetSummary{
				ChangeSetName:   aws.String(fmt.Sprintf("stratus-CREATE-%s", expectedChecksum)),
				ExecutionStatus: aws.String(cloudformation.ExecutionStatusAvailable),
			},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			assert := assert.New(t)

			actual := stratus.MatchesChangeSetSummary(stack, testCase.summary)

			assert.Equal(testCase.expected, actual)
		})
	}
}

func Test_MatchesChangeSetContents(t *testing.T) {
	stack := &config.Stack{
		Capabilities: make([]string, 0),
		Parameters: config.StackParameters{
			{
				Key:   "1",
				Value: "a",
			},
			{
				Key:   "2",
				Value: "b",
			},
		},
		Template: []byte("{}"),
	}

	testCases := []struct {
		description string

		changeSet *cloudformation.DescribeChangeSetOutput
		template  *cloudformation.GetTemplateOutput

		expected bool
	}{
		{
			description: "change set parameters are equal",
			changeSet: &cloudformation.DescribeChangeSetOutput{
				Capabilities: make([]*string, 0),
				Parameters: []*cloudformation.Parameter{
					{
						ParameterKey:   stringPointer("2"),
						ParameterValue: stringPointer("b"),
					},
					{
						ParameterKey:   stringPointer("1"),
						ParameterValue: stringPointer("a"),
					},
				},
			},
			template: &cloudformation.GetTemplateOutput{
				TemplateBody: stringPointer("{}"),
			},
			expected: true,
		},
		{
			description: "change set parameters are a superset",
			changeSet: &cloudformation.DescribeChangeSetOutput{
				Capabilities: make([]*string, 0),
				Parameters: []*cloudformation.Parameter{
					{
						ParameterKey:   stringPointer("2"),
						ParameterValue: stringPointer("b"),
					},
					{
						ParameterKey:   stringPointer("1"),
						ParameterValue: stringPointer("a"),
					},
					{
						ParameterKey:   stringPointer("3"),
						ParameterValue: stringPointer("c"),
					},
				},
			},
			template: &cloudformation.GetTemplateOutput{
				TemplateBody: stringPointer("{}"),
			},
			expected: true,
		},
		{
			description: "change set parameters are a subset",
			changeSet: &cloudformation.DescribeChangeSetOutput{
				Capabilities: make([]*string, 0),
				Parameters: []*cloudformation.Parameter{
					{
						ParameterKey:   stringPointer("1"),
						ParameterValue: stringPointer("a"),
					},
				},
			},
			template: &cloudformation.GetTemplateOutput{
				TemplateBody: stringPointer("{}"),
			},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			assert := assert.New(t)

			actual := stratus.MatchesChangeSetContents(
				stack,
				testCase.changeSet,
				testCase.template,
			)

			assert.Equal(testCase.expected, actual)
		})
	}
}
