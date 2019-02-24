package stratus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/72636c/stratus/internal/config"
	"github.com/72636c/stratus/internal/context"
	"github.com/72636c/stratus/internal/errgroup"
)

var (
	defaultOptions = []request.WaiterOption{
		request.WithWaiterDelay(request.ConstantWaiterDelay(1 * time.Second)),
		request.WithWaiterMaxAttempts(60 * 60),
	}
)

type Client struct {
	cfn CloudFormation
	s3  S3
}

func NewClient(cfn CloudFormation, s3 S3) *Client {
	return &Client{
		cfn: cfn,
		s3:  s3,
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
		TemplateBody:          nil,
		TemplateURL:           nil,
		UsePreviousTemplate:   aws.Bool(false),
	}

	if stack.TemplateKey == "" {
		input.SetTemplateBody(string(stack.Template))
	} else {
		input.SetTemplateURL(toS3URL(stack.ArtefactBucket, stack.TemplateKey))
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

	return client.describeChangeSet(ctx, stack, name)
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

	var newPolicy, oldPolicy interface{}

	err = json.Unmarshal(stack.Policy, &newPolicy)
	if err != nil {
		return nil, err
	}

	if policyOutput != nil && policyOutput.StackPolicyBody != nil {
		err := json.Unmarshal([]byte(*policyOutput.StackPolicyBody), &oldPolicy)
		if err != nil {
			return nil, err
		}
	}

	diff := &Diff{
		ChangeSet: describeOutput,
		New: &StackState{
			StackPolicy:           newPolicy,
			TerminationProtection: aws.Bool(stack.TerminationProtection),
		},
		Old: &StackState{
			StackPolicy:           oldPolicy,
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
	logger := context.Logger(ctx)

	executeInput := &cloudformation.ExecuteChangeSetInput{
		ChangeSetName:      aws.String(name),
		ClientRequestToken: nil,
		StackName:          aws.String(stack.Name),
	}

	eventsInput := &cloudformation.DescribeStackEventsInput{
		StackName: aws.String(stack.Name),
	}

	waitInput := &cloudformation.DescribeStacksInput{
		NextToken: nil,
		StackName: aws.String(stack.Name),
	}

	eventsOutput, err := client.cfn.DescribeStackEventsWithContext(
		ctx,
		eventsInput,
	)
	if err != nil {
		return err
	}

	eventCache := NewStackEventCache(eventsOutput.StackEvents)

	waiter, err := client.newChangeSetExecuteCompleteWaiter(name)
	if err != nil {
		return err
	}

	_, err = client.cfn.ExecuteChangeSetWithContext(ctx, executeInput)
	if err != nil {
		return err
	}

	option := func(req *request.Request) {
		eventsOutput, err := client.cfn.DescribeStackEventsWithContext(
			ctx,
			eventsInput,
		)
		if err != nil {
			// continue without failing request
			return
		}

		events := eventCache.Diff(eventsOutput.StackEvents)

		for index := len(events) - 1; index >= 0; index-- {
			logger.Data(formatStackEvent(events[index]))
		}
	}

	options := append(defaultOptions, request.WithWaiterRequestOptions(option))

	return waiter(ctx, waitInput, options...)

}

func (client *Client) FindExistingChangeSet(
	ctx context.Context,
	stack *config.Stack,
) (*cloudformation.DescribeChangeSetOutput, error) {
	listOutput, err := client.listChangeSets(ctx, stack)
	if isStackDoesNotExistError(err) {
		return nil, fmt.Errorf("stack '%s' does not exist", stack.Name)
	}
	if err != nil {
		return nil, err
	}

	for _, summary := range listOutput.Summaries {
		if MatchesChangeSetSummary(stack, summary) {
			name := *summary.ChangeSetName

			var (
				group *errgroup.Group

				changeSetOutput *cloudformation.DescribeChangeSetOutput
				templateOutput  *cloudformation.GetTemplateOutput
			)

			group, ctx = errgroup.WithContext(ctx)

			group.Go(func() (err error) {
				changeSetOutput, err = client.describeChangeSet(ctx, stack, name)
				return
			})

			group.Go(func() (err error) {
				templateOutput, err = client.getChangeSetTemplate(ctx, stack, name)
				return
			})

			err := group.Wait()
			if err != nil {
				return nil, err
			}

			if !MatchesChangeSetContents(stack, changeSetOutput, templateOutput) {
				return nil, fmt.Errorf(
					"change set '%s' has been modified",
					*summary.ChangeSetName,
				)
			}

			if *summary.ExecutionStatus == cloudformation.ExecutionStatusUnavailable {
				return nil, nil
			}

			return changeSetOutput, nil
		}
	}

	return nil, fmt.Errorf("change set '*%s*' does not exist", stack.Checksum)
}

func (client *Client) SetStackPolicy(
	ctx context.Context,
	stack *config.Stack,
) error {
	input := &cloudformation.SetStackPolicyInput{
		StackName:       aws.String(stack.Name),
		StackPolicyBody: nil,
		StackPolicyURL:  nil,
	}

	if stack.PolicyKey == "" {
		input.SetStackPolicyBody(string(stack.Policy))
	} else {
		input.SetStackPolicyURL(toS3URL(stack.ArtefactBucket, stack.PolicyKey))
	}

	_, err := client.cfn.SetStackPolicyWithContext(ctx, input)
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

func (client *Client) UploadArtefacts(
	ctx context.Context,
	stack *config.Stack,
) error {
	// TODO: object metadata and tagging

	policyExtension := filepath.Ext(stack.PolicyKey)
	policyContentType, ok := extensionToContentType[policyExtension]
	if !ok {
		return fmt.Errorf("unsupported policy extension '%s'", policyExtension)
	}

	templateExtension := filepath.Ext(stack.TemplateKey)
	templateContentType, ok := extensionToContentType[templateExtension]
	if !ok {
		return fmt.Errorf("unsupported template extension '%s'", templateExtension)
	}

	policyFilename := filepath.Base(stack.PolicyKey)
	policyContentDisposition := toContentDisposition(policyFilename)

	templateFilename := filepath.Base(stack.TemplateKey)
	templateContentDisposition := toContentDisposition(templateFilename)

	policyInput := &s3.PutObjectInput{
		Body:               bytes.NewReader(stack.Policy),
		Bucket:             aws.String(stack.ArtefactBucket),
		ContentDisposition: aws.String(policyContentDisposition),
		ContentType:        aws.String(policyContentType),
		Key:                aws.String(stack.PolicyKey),
	}

	templateInput := &s3.PutObjectInput{
		Body:               bytes.NewReader(stack.Template),
		Bucket:             aws.String(stack.ArtefactBucket),
		ContentDisposition: aws.String(templateContentDisposition),
		ContentType:        aws.String(templateContentType),
		Key:                aws.String(stack.TemplateKey),
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() (err error) {
		_, err = client.s3.PutObjectWithContext(ctx, policyInput)
		return
	})

	group.Go(func() (err error) {
		_, err = client.s3.PutObjectWithContext(ctx, templateInput)
		return
	})

	return group.Wait()
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

func (client *Client) describeChangeSet(
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

func (client *Client) getChangeSetTemplate(
	ctx context.Context,
	stack *config.Stack,
	name string,
) (*cloudformation.GetTemplateOutput, error) {
	input := &cloudformation.GetTemplateInput{
		ChangeSetName: aws.String(name),
		StackName:     aws.String(stack.Name),
		TemplateStage: aws.String(cloudformation.TemplateStageOriginal),
	}

	return client.cfn.GetTemplateWithContext(ctx, input)
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

	describeOutput, describeError := client.describeChangeSet(ctx, stack, name)
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
) (stackWaiter, error) {
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
		WaitUntilChangeSetCreateCompleteWithContext(ctx, input, defaultOptions...)
}

func (client *Client) waitUntilStackDeleteComplete(
	ctx context.Context,
	stack *config.Stack,
) error {
	input := &cloudformation.DescribeStacksInput{
		NextToken: nil,
		StackName: aws.String(stack.Name),
	}

	return client.cfn.
		WaitUntilStackDeleteCompleteWithContext(ctx, input, defaultOptions...)
}
