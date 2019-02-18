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
	changeSetTypeToString = map[ChangeSetType]string{
		ChangeSetTypeCreate: cloudformation.ChangeSetTypeCreate,
		ChangeSetTypeUpdate: cloudformation.ChangeSetTypeUpdate,
	}

	stringToChangeSetType = map[string]ChangeSetType{
		cloudformation.ChangeSetTypeCreate: ChangeSetTypeCreate,
		cloudformation.ChangeSetTypeUpdate: ChangeSetTypeUpdate,
	}
)

type ChangeSetType int

func ParseChangeSetType(raw string) (changeSetType ChangeSetType, ok bool) {
	changeSetType, ok = stringToChangeSetType[strings.ToUpper(raw)]
	return
}

func (changeSetType ChangeSetType) String() string {
	return changeSetTypeToString[changeSetType]
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

type StackEventCache struct {
	ids []string
}

func NewStackEventCache(events []*cloudformation.StackEvent) *StackEventCache {
	slice := make([]string, len(events))

	for index, event := range events {
		slice[index] = *event.EventId
	}

	return &StackEventCache{
		ids: slice,
	}
}

func (cache *StackEventCache) Contains(event *cloudformation.StackEvent) bool {
	if event == nil || event.EventId == nil {
		return false
	}

	for _, id := range cache.ids {
		if id == *event.EventId {
			return true
		}
	}

	return false
}

func (cache *StackEventCache) Diff(
	events []*cloudformation.StackEvent,
) []*cloudformation.StackEvent {
	slice := make([]*cloudformation.StackEvent, 0)

	for _, event := range events {
		if !cache.Contains(event) {
			cache.ids = append(cache.ids, *event.EventId)
			slice = append(slice, event)
		}
	}

	return slice
}

type StackState struct {
	StackPolicy           interface{}
	TerminationProtection *bool
}

type stackWaiter func(
	aws.Context,
	*cloudformation.DescribeStacksInput,
	...request.WaiterOption,
) error
