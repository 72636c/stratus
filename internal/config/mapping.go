package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func fromRawConfig(
	raw *RawConfig,
	checksum string,
	relativePath string,
) (*Config, error) {
	stacks, err := fromRawStacks(raw.Stacks, checksum, relativePath)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Stacks: stacks,
	}

	return config, nil
}

func fromRawStacks(
	raw []*RawStack,
	checksum string,
	relativePath string,
) ([]*Stack, error) {
	slice := make([]*Stack, len(raw))

	for index, stack := range raw {
		var err error

		slice[index], err = fromRawStack(stack, checksum, relativePath)
		if err != nil {
			return nil, err
		}
	}

	return slice, nil
}

func fromRawStack(
	raw *RawStack,
	checksum string,
	relativePath string,
) (*Stack, error) {
	policyPath := filepath.Join(filepath.Dir(relativePath), *raw.PolicyFile)

	policyData, err := ioutil.ReadFile(policyPath)
	if err != nil {
		return nil, err
	}

	var policy interface{}

	err = json.Unmarshal(policyData, &policy)
	if err != nil {
		return nil, err
	}

	templatePath := filepath.Join(filepath.Dir(relativePath), *raw.TemplateFile)

	template, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	stack := &Stack{
		Name: *raw.Name,

		Capabilities:          raw.Capabilities,
		Parameters:            fromRawStackParameters(raw.Parameters),
		Tags:                  fromRawStackTags(raw.Tags),
		TerminationProtection: *raw.TerminationProtection,

		Policy:   policy,
		Template: template,

		Checksum: checksum,
	}

	return stack, nil
}

func fromRawStackParameters(
	raw RawStackParameters,
) StackParameters {
	slice := make(StackParameters, len(raw))

	for index, rawParameter := range raw {
		parameter := StackParameter(*rawParameter)
		slice[index] = &parameter
	}

	return slice
}

func fromRawStackTags(
	raw RawStackTags,
) StackTags {
	slice := make(StackTags, len(raw))

	for index, rawTag := range raw {
		tag := StackTag(*rawTag)
		slice[index] = &tag
	}

	return slice
}
