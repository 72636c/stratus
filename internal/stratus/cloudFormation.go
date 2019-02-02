package stratus

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

var (
	_ CloudFormation = new(cloudformation.CloudFormation)
	_ CloudFormation = new(CloudFormationMock)
)

type CloudFormation interface {
	CreateChangeSetWithContext(
		aws.Context,
		*cloudformation.CreateChangeSetInput,
		...request.Option,
	) (*cloudformation.CreateChangeSetOutput, error)

	DeleteStackWithContext(
		aws.Context,
		*cloudformation.DeleteStackInput,
		...request.Option,
	) (*cloudformation.DeleteStackOutput, error)

	DescribeChangeSetWithContext(
		aws.Context,
		*cloudformation.DescribeChangeSetInput,
		...request.Option,
	) (*cloudformation.DescribeChangeSetOutput, error)

	DescribeStacksWithContext(
		aws.Context,
		*cloudformation.DescribeStacksInput,
		...request.Option,
	) (*cloudformation.DescribeStacksOutput, error)

	ExecuteChangeSetWithContext(
		aws.Context,
		*cloudformation.ExecuteChangeSetInput,
		...request.Option,
	) (*cloudformation.ExecuteChangeSetOutput, error)

	GetStackPolicyWithContext(
		aws.Context,
		*cloudformation.GetStackPolicyInput,
		...request.Option,
	) (*cloudformation.GetStackPolicyOutput, error)

	ListChangeSetsWithContext(
		aws.Context,
		*cloudformation.ListChangeSetsInput,
		...request.Option,
	) (*cloudformation.ListChangeSetsOutput, error)

	SetStackPolicyWithContext(
		aws.Context,
		*cloudformation.SetStackPolicyInput,
		...request.Option,
	) (*cloudformation.SetStackPolicyOutput, error)

	UpdateTerminationProtectionWithContext(
		aws.Context,
		*cloudformation.UpdateTerminationProtectionInput,
		...request.Option,
	) (*cloudformation.UpdateTerminationProtectionOutput, error)

	WaitUntilChangeSetCreateCompleteWithContext(
		aws.Context,
		*cloudformation.DescribeChangeSetInput,
		...request.WaiterOption,
	) error

	WaitUntilStackCreateCompleteWithContext(
		aws.Context,
		*cloudformation.DescribeStacksInput,
		...request.WaiterOption,
	) error

	WaitUntilStackDeleteCompleteWithContext(
		aws.Context,
		*cloudformation.DescribeStacksInput,
		...request.WaiterOption,
	) error

	WaitUntilStackUpdateCompleteWithContext(
		aws.Context,
		*cloudformation.DescribeStacksInput,
		...request.WaiterOption,
	) error

	ValidateTemplateWithContext(
		aws.Context,
		*cloudformation.ValidateTemplateInput,
		...request.Option,
	) (*cloudformation.ValidateTemplateOutput, error)
}
