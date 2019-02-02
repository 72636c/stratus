package stratus

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"golang.org/x/sync/errgroup"

	"github.com/72636c/stratus/internal/config"
)

var (
	delay = request.WithWaiterDelay(request.ConstantWaiterDelay(1 * time.Second))
)

type Client struct {
	cfn CloudFormation
}

func NewClient(cfn CloudFormation) *Client {
	return &Client{
		cfn: cfn,
	}
}

func (client *Client) CreateChangeSet(
	ctx context.Context,
	stack *config.Stack,
) (_ *cloudformation.DescribeChangeSetOutput, err error) {
	name := newChangeSetName(stack.Checksum, ChangeSetTypeUpdate)

	input := &cloudformation.CreateChangeSetInput{
		Capabilities:          aws.StringSlice(stack.Capabilities),
		ChangeSetName:         aws.String(name),
		ChangeSetType:         aws.String(ChangeSetTypeUpdate.String()),
		ClientToken:           nil,
		Description:           nil,
		NotificationARNs:      nil,
		Parameters:            toCloudFormationParameters(stack.Parameters),
		ResourceTypes:         nil,
		RoleARN:               nil,
		RollbackConfiguration: nil,
		StackName:             aws.String(stack.Name),
		Tags:                  toCloudFormationTags(stack.Tags),
		TemplateBody:          aws.String(string(stack.Template)),
		TemplateURL:           nil,
		UsePreviousTemplate:   aws.Bool(false),
	}

	_, err = client.cfn.CreateChangeSetWithContext(ctx, input)
	if isStackDoesNotExistError(err) {
		name = newChangeSetName(stack.Checksum, ChangeSetTypeCreate)
		input.SetChangeSetName(name)
		input.SetChangeSetType(ChangeSetTypeCreate.String())
		_, err = client.cfn.CreateChangeSetWithContext(ctx, input)
	}
	if err != nil {
		return nil, err
	}

	err = client.waitUntilChangeSetCreateComplete(ctx, stack, name)
	if err != nil {
		return nil, client.handleCreateChangeSetError(ctx, stack, name, err)
	}

	return client.DescribeChangeSet(ctx, stack, name)
}

func (client *Client) DeleteStack(
	ctx context.Context,
	stack *config.Stack,
) error {
	input := &cloudformation.DeleteStackInput{
		ClientRequestToken: nil,
		RetainResources:    nil,
		RoleARN:            nil,
		StackName:          aws.String(stack.Name),
	}

	_, err := client.cfn.DeleteStackWithContext(ctx, input)
	if err != nil {
		return err
	}

	return client.waitUntilStackDeleteComplete(ctx, stack)
}

func (client *Client) DescribeChangeSet(
	ctx context.Context,
	stack *config.Stack,
	name string,
) (*cloudformation.DescribeChangeSetOutput, error) {
	input := &cloudformation.DescribeChangeSetInput{
		ChangeSetName: aws.String(name),
		// TODO: handle pagination
		NextToken: nil,
		StackName: aws.String(stack.Name),
	}

	return client.cfn.DescribeChangeSetWithContext(ctx, input)
}

func (client *Client) Diff(
	ctx context.Context,
	stack *config.Stack,
	describeOutput *cloudformation.DescribeChangeSetOutput,
) (*Diff, error) {
	group, ctx := errgroup.WithContext(ctx)

	var description *cloudformation.Stack
	var policyOutput *cloudformation.GetStackPolicyOutput

	group.Go(func() (err error) {
		description, err = client.describeStack(ctx, stack)
		return
	})

	group.Go(func() (err error) {
		policyOutput, err = client.getStackPolicy(ctx, stack)
		return
	})

	err := group.Wait()
	if err != nil {
		return nil, err
	}

	var policy interface{}

	if policyOutput != nil && policyOutput.StackPolicyBody != nil {
		err = json.Unmarshal([]byte(*policyOutput.StackPolicyBody), &policy)
		if err != nil {
			return nil, err
		}
	}

	diff := &Diff{
		ChangeSet: describeOutput,
		New: &StackState{
			StackPolicy:           stack.Policy,
			TerminationProtection: aws.Bool(stack.TerminationProtection),
		},
		Old: &StackState{
			StackPolicy:           policy,
			TerminationProtection: description.EnableTerminationProtection,
		},
	}

	return diff, nil
}

func (client *Client) ExecuteChangeSet(
	ctx context.Context,
	stack *config.Stack,
	name string,
) error {
	executeInput := &cloudformation.ExecuteChangeSetInput{
		ChangeSetName:      aws.String(name),
		ClientRequestToken: nil,
		StackName:          aws.String(stack.Name),
	}

	describeInput := &cloudformation.DescribeStacksInput{
		NextToken: nil,
		StackName: aws.String(stack.Name),
	}

	waiter, err := client.newChangeSetExecuteCompleteWaiter(name)
	if err != nil {
		return err
	}

	_, err = client.cfn.ExecuteChangeSetWithContext(ctx, executeInput)
	if err != nil {
		return err
	}

	return waiter(ctx, describeInput)
}

func (client *Client) FindExistingChangeSet(
	ctx context.Context,
	stack *config.Stack,
) (*cloudformation.DescribeChangeSetOutput, error) {
	output, err := client.listChangeSets(ctx, stack)
	if isStackDoesNotExistError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	for _, summary := range output.Summaries {
		if matchesChangeSetSummary(stack, summary) {
			return client.DescribeChangeSet(ctx, stack, *summary.ChangeSetName)
		}
	}

	return nil, nil
}

func (client *Client) SetStackPolicy(
	ctx context.Context,
	stack *config.Stack,
) error {
	policyData, err := json.Marshal(stack.Policy)
	if err != nil {
		return err
	}

	input := &cloudformation.SetStackPolicyInput{
		StackName:       aws.String(stack.Name),
		StackPolicyBody: aws.String(string(policyData)),
		StackPolicyURL:  nil,
	}

	_, err = client.cfn.SetStackPolicyWithContext(ctx, input)
	return err
}

func (client *Client) UpdateTerminationProtection(
	ctx context.Context,
	stack *config.Stack,
) error {
	input := &cloudformation.UpdateTerminationProtectionInput{
		EnableTerminationProtection: aws.Bool(stack.TerminationProtection),
		StackName:                   aws.String(stack.Name),
	}

	_, err := client.cfn.UpdateTerminationProtectionWithContext(ctx, input)
	return err
}

func (client *Client) ValidateTemplate(
	ctx context.Context,
	stack *config.Stack,
) (*cloudformation.ValidateTemplateOutput, error) {
	input := &cloudformation.ValidateTemplateInput{
		TemplateBody: aws.String(string(stack.Template)),
		TemplateURL:  nil,
	}

	// TODO: check for insufficient capabilities

	return client.cfn.ValidateTemplateWithContext(ctx, input)
}

func (client *Client) describeStack(
	ctx context.Context,
	stack *config.Stack,
) (*cloudformation.Stack, error) {
	input := &cloudformation.DescribeStacksInput{
		NextToken: nil,
		StackName: aws.String(stack.Name),
	}

	describeOutput, err := client.cfn.DescribeStacksWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	if describeOutput == nil || len(describeOutput.Stacks) != 1 {
		return nil, errors.New("cloudformation.DescribeStacks: invalid output")
	}

	return describeOutput.Stacks[0], nil
}

func (client *Client) getStackPolicy(
	ctx context.Context,
	stack *config.Stack,
) (*cloudformation.GetStackPolicyOutput, error) {
	input := &cloudformation.GetStackPolicyInput{
		StackName: aws.String(stack.Name),
	}

	return client.cfn.GetStackPolicyWithContext(ctx, input)
}

func (client *Client) handleCreateChangeSetError(
	ctx context.Context,
	stack *config.Stack,
	name string,
	err error,
) error {
	if !isResourceNotReadyError(err) {
		return err
	}

	describeOutput, describeError := client.DescribeChangeSet(ctx, stack, name)
	if describeError != nil {
		return err
	}

	if !isNoopChangeSet(describeOutput) {
		return err
	}

	return nil
}

func (client *Client) listChangeSets(
	ctx context.Context,
	stack *config.Stack,
) (*cloudformation.ListChangeSetsOutput, error) {
	input := &cloudformation.ListChangeSetsInput{
		// TODO: handle pagination
		NextToken: nil,
		StackName: aws.String(stack.Name),
	}

	return client.cfn.ListChangeSetsWithContext(ctx, input)
}

func (client *Client) newChangeSetExecuteCompleteWaiter(
	name string,
) (StackWaiter, error) {
	changeSetType, err := getChangeSetType(name)

	switch changeSetType {
	case ChangeSetTypeCreate:
		return client.cfn.WaitUntilStackCreateCompleteWithContext, nil

	case ChangeSetTypeUpdate:
		return client.cfn.WaitUntilStackUpdateCompleteWithContext, nil

	default:
		return nil, err
	}
}

func (client *Client) waitUntilChangeSetCreateComplete(
	ctx context.Context,
	stack *config.Stack,
	name string,
) error {
	input := &cloudformation.DescribeChangeSetInput{
		ChangeSetName: aws.String(name),
		NextToken:     nil,
		StackName:     aws.String(stack.Name),
	}

	return client.cfn.
		WaitUntilChangeSetCreateCompleteWithContext(ctx, input, delay)
}

func (client *Client) waitUntilStackDeleteComplete(
	ctx context.Context,
	stack *config.Stack,
) error {
	input := &cloudformation.DescribeStacksInput{
		NextToken: nil,
		StackName: aws.String(stack.Name),
	}

	return client.cfn.WaitUntilStackDeleteCompleteWithContext(ctx, input, delay)
}
