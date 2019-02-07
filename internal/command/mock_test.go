package command_test

import (
	"encoding/json"
	"fmt"
)

const (
	mockChecksum = "1000000000200000000030000000004000000000500000000060000000007000"

	mockStackName     = "test-stack-name"
	mockStackPolicy   = "test-stack-policy"
	mockStackTemplate = "test-stack-template"
)

var (
	mockChangeSetCreateName = fmt.Sprintf("stratus-create-%s", mockChecksum)
	mockChangeSetUpdateName = fmt.Sprintf("stratus-update-%s", mockChecksum)

	mockStackPolicyBody = mustMarshal(mockStackPolicy)
)

func mustMarshal(model interface{}) []byte {
	data, err := json.Marshal(model)
	if err != nil {
		panic(fmt.Errorf("mustMarshal: %+v", err))
	}

	return data
}
