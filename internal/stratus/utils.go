package stratus

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/72636c/stratus/internal/config"
)

const (
	noopChangeSetStatusReason = "The submitted information didn't contain changes. Submit different information to create a change set."
)

var (
	changeSetRegexp = regexp.MustCompile(`stratus-(create|update)-[0-9a-f]{64}`)

	extensionToContentType = map[string]string{
		".json": "application/json; charset=utf-8",
		".yaml": "application/x-yaml; charset=utf-8",
		".yml":  "application/x-yaml; charset=utf-8",
	}
)

func MatchesChangeSetContents(
	stack *config.Stack,
	changeSet *cloudformation.DescribeChangeSetOutput,
	template *cloudformation.GetTemplateOutput,
) bool {
	return string(stack.Template) == *template.TemplateBody &&
		matchesChangeSetCapabilities(stack.Capabilities, changeSet.Capabilities) &&
		matchesChangeSetParameters(stack.Parameters, changeSet.Parameters)
}

func MatchesChangeSetSummary(
	stack *config.Stack,
	summary *cloudformation.ChangeSetSummary,
) bool {
	return summary != nil &&
		summary.ChangeSetName != nil &&
		isAcceptableChangeSetStatus(summary) &&
		matchesChangeSetName(stack.Checksum, *summary.ChangeSetName)
}

func getChangeSetType(name string) (ChangeSetType, error) {
	raw := changeSetRegexp.FindStringSubmatch(name)
	if len(raw) != 2 {
		return 0, fmt.Errorf("unrecognised change set name format '%s'", name)
	}

	changeSetType, ok := ParseChangeSetType(raw[1])
	if !ok {
		return 0, fmt.Errorf("unrecognised change set type '%s", raw[1])
	}

	return changeSetType, nil
}

func isAcceptableChangeSetStatus(
	summary *cloudformation.ChangeSetSummary,
) bool {
	if summary.ExecutionStatus == nil {
		return false
	}

	if *summary.ExecutionStatus == cloudformation.ExecutionStatusAvailable {
		return true
	}

	if summary.Status == nil || summary.StatusReason == nil {
		return false
	}

	if *summary.ExecutionStatus == cloudformation.ExecutionStatusUnavailable &&
		*summary.Status == cloudformation.ChangeSetStatusFailed &&
		*summary.StatusReason == noopChangeSetStatusReason {
		return true
	}

	return false
}

func isNoopChangeSet(output *cloudformation.DescribeChangeSetOutput) bool {
	return output != nil &&
		output.Status != nil &&
		output.StatusReason != nil &&
		*output.Status == cloudformation.ChangeSetStatusFailed &&
		*output.StatusReason == noopChangeSetStatusReason
}

func isResourceNotReadyError(err error) bool {
	if err == nil {
		return false
	}

	awsError, ok := err.(awserr.Error)
	if !ok {
		return false
	}

	return awsError.Code() == request.WaiterResourceNotReadyErrorCode
}

func isStackDoesNotExistError(err error) bool {
	if err == nil {
		return false
	}

	awsError, ok := err.(awserr.Error)
	if !ok {
		return false
	}

	return awsError.Code() == "ValidationError" &&
		strings.Contains(awsError.Message(), "does not exist")
}

func matchesChangeSetCapabilities(expected []string, actual []*string) bool {
	return reflect.DeepEqual(
		sort.StringSlice(expected),
		sort.StringSlice(toStringList(actual)),
	)
}

func matchesChangeSetName(checksum, changeSetName string) bool {
	str := fmt.Sprintf("stratus-(create|update)-%s", regexp.QuoteMeta(checksum))

	return regexp.MustCompile(str).MatchString(changeSetName)
}

func matchesChangeSetParameters(
	expected config.StackParameters,
	actual []*cloudformation.Parameter,
) bool {
	actualValues := make(config.StackParameters, 0)

	for _, parameter := range actual {
		if parameter != nil {
			value := &config.StackParameter{
				Key:   *parameter.ParameterKey,
				Value: *parameter.ParameterValue,
			}

			actualValues = append(actualValues, value)
		}
	}

	if len(expected) != len(actualValues) {
		return false
	}

	for _, value := range actualValues {
		if !expected.Contains(value.Key, value.Value) {
			return false
		}
	}

	return true
}

func newChangeSetName(checksum string, changeSetType ChangeSetType) string {
	return fmt.Sprintf(
		"stratus-%s-%s",
		strings.ToLower(changeSetType.String()),
		checksum,
	)
}

func toCloudFormationParameters(
	parameters config.StackParameters,
) []*cloudformation.Parameter {
	slice := make([]*cloudformation.Parameter, len(parameters))

	for index, parameter := range parameters {
		slice[index] = toCloudFormationParameter(parameter)
	}

	return slice
}

func toCloudFormationParameter(
	parameter *config.StackParameter,
) *cloudformation.Parameter {
	return &cloudformation.Parameter{
		ParameterKey:     aws.String(parameter.Key),
		ParameterValue:   aws.String(parameter.Value),
		ResolvedValue:    nil,
		UsePreviousValue: nil,
	}
}

func toCloudFormationTags(tags config.StackTags) []*cloudformation.Tag {
	slice := make([]*cloudformation.Tag, len(tags))

	for index, tag := range tags {
		slice[index] = toCloudFormationTag(tag)
	}

	return slice
}

func toCloudFormationTag(tag *config.StackTag) *cloudformation.Tag {
	return &cloudformation.Tag{
		Key:   aws.String(tag.Key),
		Value: aws.String(tag.Value),
	}
}

func toS3URL(bucket, key string) string {
	return fmt.Sprintf("https://s3.amazonaws.com/%s/%s", bucket, key)
}

func toStringList(xs []*string) []string {
	slice := make([]string, 0)

	for _, x := range xs {
		if x != nil {
			slice = append(slice, *x)
		}
	}

	return slice
}
