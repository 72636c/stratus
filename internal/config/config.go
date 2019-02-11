package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awsutil"
)

type Config struct {
	Stacks []*Stack
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

type Stack struct {
	Name string

	Capabilities          []string
	Parameters            StackParameters
	Tags                  StackTags
	TerminationProtection bool

	Policy   []byte
	Template []byte

	ArtefactBucket string
	PolicyKey      string `json:"-"`
	TemplateKey    string `json:"-"`

	Checksum string `json:"-"`
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
