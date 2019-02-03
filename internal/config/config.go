package config

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"gopkg.in/yaml.v2"
)

var (
	unmarshals = map[string]func([]byte, interface{}) error{
		".json": json.Unmarshal,
		".yaml": yaml.UnmarshalStrict,
		".yml":  yaml.UnmarshalStrict,
	}
)

type Config struct {
	Stacks []*Stack
}

func FromPath(path string) (*Config, error) {
	extension := strings.ToLower(filepath.Ext(path))

	unmarshal, ok := unmarshals[extension]
	if !ok {
		return nil, fmt.Errorf("unrecognised config file extension '%s'", extension)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	checksum := fmt.Sprintf("%x", sha256.Sum256(data))

	var raw *RawConfig

	err = unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}

	// TODO: validate config

	return fromRawConfig(raw, checksum, path)
}

type Stack struct {
	Name string

	Capabilities          []string
	Parameters            StackParameters
	Tags                  StackTags
	TerminationProtection bool

	Policy   interface{}
	Template []byte

	Checksum string
}

func (stack *Stack) String() string {
	return awsutil.Prettify(stack)
}

type StackParameters []*StackParameter

type StackParameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StackTags []*StackTag

type StackTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
