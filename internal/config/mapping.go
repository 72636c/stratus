package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func fromRawConfig(rawConfig *RawConfig, path string) (*Config, error) {
	stacks, err := fromRawStacks(rawConfig, path)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Stacks: stacks,
	}

	return config, nil
}

func fromRawStacks(rawConfig *RawConfig, path string) (Stacks, error) {
	slice := make(Stacks, len(rawConfig.Stacks))

	for index, stack := range rawConfig.Stacks {
		var err error

		slice[index], err = fromRawStack(rawConfig, stack, path)
		if err != nil {
			return nil, err
		}
	}

	return slice, nil
}

func fromRawStack(
	rawConfig *RawConfig,
	rawStack *RawStack,
	path string,
) (*Stack, error) {
	policyPath := filepath.Join(
		filepath.Dir(path),
		rawStack.PolicyFile.String(),
	)

	policy, err := ioutil.ReadFile(policyPath)
	if err != nil {
		return nil, err
	}

	templatePath := filepath.Join(
		filepath.Dir(path),
		rawStack.TemplateFile.String(),
	)

	template, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	stack := &Stack{
		Name: rawStack.Name.String(),

		Capabilities:          fromRawStackCapabilities(rawStack.Capabilities),
		Parameters:            fromRawStackParameters(rawStack.Parameters),
		Tags:                  fromRawStackTags(rawStack.Tags),
		TerminationProtection: rawStack.TerminationProtection.Bool(),

		Policy:   policy,
		Template: template,

		ArtefactBucket: rawConfig.Defaults.ArtefactBucket.String(),
	}

	checksum, err := CalculateChecksum(stack.Hashable())
	if err != nil {
		return nil, err
	}

	stack.Checksum = checksum

	if stack.ArtefactBucket != "" {
		keyFormat := fmt.Sprintf("stratus/%s/%s/%%s%%s", stack.Name, stack.Checksum)

		policyExtension := filepath.Ext(rawStack.PolicyFile.String())
		templateExtension := filepath.Ext(rawStack.TemplateFile.String())

		stack.PolicyKey = fmt.Sprintf(keyFormat, "policy", policyExtension)
		stack.TemplateKey = fmt.Sprintf(keyFormat, "template", templateExtension)
	}

	return stack, nil
}

func fromRawStackCapabilities(raw RawStackCapabilities) []string {
	slice := make([]string, len(raw))

	for index, rawCapability := range raw {
		slice[index] = rawCapability.String()
	}

	return slice
}

func fromRawStackParameters(raw RawStackParameters) StackParameters {
	slice := make(StackParameters, len(raw))

	for index, rawParameter := range raw {
		slice[index] = &StackParameter{
			Key:   rawParameter.Key.String(),
			Value: rawParameter.Value.String(),
		}
	}

	return slice
}

func fromRawStackTags(raw RawStackTags) StackTags {
	slice := make(StackTags, len(raw))

	for index, rawTag := range raw {
		slice[index] = &StackTag{
			Key:   rawTag.Key.String(),
			Value: rawTag.Value.String(),
		}
	}

	return slice
}
