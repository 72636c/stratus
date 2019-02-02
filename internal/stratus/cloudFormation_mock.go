package stratus

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/stretchr/testify/mock"
)

type CloudFormationMock struct {
	mock.Mock
}

func NewCloudFormationMock() *CloudFormationMock {
	return new(CloudFormationMock)
}

func (client *CloudFormationMock) CreateChangeSetWithContext(
	_ aws.Context,
	input *cloudformation.CreateChangeSetInput,
	_ ...request.Option,
) (*cloudformation.CreateChangeSetOutput, error) {
	args := client.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudformation.CreateChangeSetOutput), args.Error(1)
}

func (client *CloudFormationMock) DeleteStackWithContext(
	_ aws.Context,
	input *cloudformation.DeleteStackInput,
	_ ...request.Option,
) (*cloudformation.DeleteStackOutput, error) {
	args := client.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudformation.DeleteStackOutput), args.Error(1)
}

func (client *CloudFormationMock) DescribeChangeSetWithContext(
	_ aws.Context,
	input *cloudformation.DescribeChangeSetInput,
	_ ...request.Option,
) (*cloudformation.DescribeChangeSetOutput, error) {
	args := client.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudformation.DescribeChangeSetOutput), args.Error(1)
}

func (client *CloudFormationMock) DescribeStacksWithContext(
	_ aws.Context,
	input *cloudformation.DescribeStacksInput,
	_ ...request.Option,
) (*cloudformation.DescribeStacksOutput, error) {
	args := client.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudformation.DescribeStacksOutput), args.Error(1)
}

func (client *CloudFormationMock) ExecuteChangeSetWithContext(
	_ aws.Context,
	input *cloudformation.ExecuteChangeSetInput,
	_ ...request.Option,
) (*cloudformation.ExecuteChangeSetOutput, error) {
	args := client.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudformation.ExecuteChangeSetOutput), args.Error(1)
}

func (client *CloudFormationMock) GetStackPolicyWithContext(
	_ aws.Context,
	input *cloudformation.GetStackPolicyInput,
	_ ...request.Option,
) (*cloudformation.GetStackPolicyOutput, error) {
	args := client.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudformation.GetStackPolicyOutput), args.Error(1)
}

func (client *CloudFormationMock) ListChangeSetsWithContext(
	_ aws.Context,
	input *cloudformation.ListChangeSetsInput,
	_ ...request.Option,
) (*cloudformation.ListChangeSetsOutput, error) {
	args := client.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudformation.ListChangeSetsOutput), args.Error(1)
}

func (client *CloudFormationMock) SetStackPolicyWithContext(
	_ aws.Context,
	input *cloudformation.SetStackPolicyInput,
	_ ...request.Option,
) (*cloudformation.SetStackPolicyOutput, error) {
	args := client.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudformation.SetStackPolicyOutput), args.Error(1)
}

func (client *CloudFormationMock) UpdateTerminationProtectionWithContext(
	_ aws.Context,
	input *cloudformation.UpdateTerminationProtectionInput,
	_ ...request.Option,
) (*cloudformation.UpdateTerminationProtectionOutput, error) {
	args := client.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudformation.UpdateTerminationProtectionOutput), args.Error(1)
}

func (client *CloudFormationMock) WaitUntilChangeSetCreateCompleteWithContext(
	_ aws.Context,
	input *cloudformation.DescribeChangeSetInput,
	_ ...request.WaiterOption,
) error {
	args := client.Called(input)
	return args.Error(0)
}

func (client *CloudFormationMock) WaitUntilStackCreateCompleteWithContext(
	_ aws.Context,
	input *cloudformation.DescribeStacksInput,
	_ ...request.WaiterOption,
) error {
	args := client.Called(input)
	return args.Error(0)
}

func (client *CloudFormationMock) WaitUntilStackDeleteCompleteWithContext(
	_ aws.Context,
	input *cloudformation.DescribeStacksInput,
	_ ...request.WaiterOption,
) error {
	args := client.Called(input)
	return args.Error(0)
}

func (client *CloudFormationMock) WaitUntilStackUpdateCompleteWithContext(
	_ aws.Context,
	input *cloudformation.DescribeStacksInput,
	_ ...request.WaiterOption,
) error {
	args := client.Called(input)
	return args.Error(0)
}

func (client *CloudFormationMock) ValidateTemplateWithContext(
	_ aws.Context,
	input *cloudformation.ValidateTemplateInput,
	_ ...request.Option,
) (*cloudformation.ValidateTemplateOutput, error) {
	args := client.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudformation.ValidateTemplateOutput), args.Error(1)
}
