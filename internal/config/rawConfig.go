package config

type RawConfig struct {
	Stacks []*RawStack `json:"stacks"`
}

type RawStack struct {
	Name *string `json:"name"`

	Capabilities          []string        `json:"capabilities"`
	Parameters            StackParameters `json:"parameters"`
	Tags                  StackTags       `json:"tags"`
	TerminationProtection *bool           `json:"terminationProtection" yaml:"terminationProtection"`

	PolicyFile   *string `json:"policyFile" yaml:"policyFile"`
	TemplateFile *string `json:"templateFile" yaml:"templateFile"`
}
