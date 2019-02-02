package stratus

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

const (
	_ ChangeSetType = iota
	ChangeSetTypeCreate
	ChangeSetTypeUpdate
)

var (
	changeSetTypeStrings = map[ChangeSetType]string{
		ChangeSetTypeCreate: cloudformation.ChangeSetTypeCreate,
		ChangeSetTypeUpdate: cloudformation.ChangeSetTypeUpdate,
	}

	changeSetTypeValues = map[string]ChangeSetType{
		cloudformation.ChangeSetTypeCreate: ChangeSetTypeCreate,
		cloudformation.ChangeSetTypeUpdate: ChangeSetTypeUpdate,
	}
)

type ChangeSetType int

func ParseChangeSetType(raw string) (changeSetType ChangeSetType, ok bool) {
	changeSetType, ok = changeSetTypeValues[strings.ToUpper(raw)]
	return
}

func (changeSetType ChangeSetType) String() string {
	return changeSetTypeStrings[changeSetType]
}

type Diff struct {
	ChangeSet *cloudformation.DescribeChangeSetOutput
	New       *StackState
	Old       *StackState
}

func (diff *Diff) HasChangeSet() bool {
	return diff.ChangeSet != nil
}

func (diff *Diff) String() string {
	return awsutil.Prettify(diff)
}

type StackState struct {
	StackPolicy           interface{}
	TerminationProtection *bool
}

type StackWaiter func(
	aws.Context,
	*cloudformation.DescribeStacksInput,
	...request.WaiterOption,
) error
