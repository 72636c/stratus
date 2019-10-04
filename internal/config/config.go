package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awsutil"
)

type Config struct {
	Stacks Stacks
}

func FromPath(path string) (*Config, error) {
	extension := strings.ToLower(filepath.Ext(path))

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw *RawConfig

	err = Unmarshal(extension, data, &raw)
	if err != nil {
		return nil, err
	}

	// TODO: validate config

	return fromRawConfig(raw, path)
}

type Stacks []*Stack

func (stacks Stacks) Find(stackName string) (*Stack, bool) {
	for _, stack := range stacks {
		if stack.Name == stackName {
			return stack, true
		}
	}

	return nil, false
}

type Stack struct {
	Name string

	Capabilities          []string
	Parameters            StackParameters
	Region                *string
	Tags                  StackTags
	TerminationProtection bool

	Policy   []byte `json:"-"`
	Template []byte `json:"-"`

	ArtefactBucket string `json:",omitempty"`
	PolicyKey      string `json:",omitempty"`
	TemplateKey    string `json:",omitempty"`

	Checksum string
}

func (stack *Stack) Hashable() interface{} {
	if stack == nil {
		return stack
	}

	return struct {
		Name string

		Capabilities          []string
		Parameters            StackParameters
		Region                *string `json:"-"`
		Tags                  StackTags
		TerminationProtection bool

		Policy   []byte
		Template []byte

		ArtefactBucket string
		PolicyKey      string `json:"-"`
		TemplateKey    string `json:"-"`

		Checksum string `json:"-"`
	}(*stack)
}

func (stack *Stack) ShouldUpload() bool {
	return stack.ArtefactBucket != ""
}

func (stack *Stack) String() string {
	return awsutil.Prettify(stack)
}

type StackParameters []*StackParameter

func (parameters StackParameters) Contains(key, value string) bool {
	for _, parameter := range parameters {
		if parameter.Key == key && parameter.Value == value {
			return true
		}
	}

	return false
}

type StackParameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StackTags []*StackTag

type StackTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
