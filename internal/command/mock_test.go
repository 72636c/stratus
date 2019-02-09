package command_test

import (
	"fmt"
)

const (
	mockChecksum = "1000000000200000000030000000004000000000500000000060000000007000"

	mockArtefactBucket   = "test-bucket-name"
	mockStackPolicyKey   = "test-policy-key"
	mockStackTemplateKey = "test-template-key"

	mockStackName     = "test-stack-name"
	mockStackPolicy   = `"test-stack-policy"`
	mockStackTemplate = "test-stack-template"
)

var (
	mockChangeSetCreateName = fmt.Sprintf("stratus-create-%s", mockChecksum)
	mockChangeSetUpdateName = fmt.Sprintf("stratus-update-%s", mockChecksum)

	mockStackPolicyURL   = fmt.Sprintf("https://s3.amazonaws.com/%s/%s", mockArtefactBucket, mockStackPolicyKey)
	mockStackTemplateURL = fmt.Sprintf("https://s3.amazonaws.com/%s/%s", mockArtefactBucket, mockStackTemplateKey)
)
