package config

type RawConfig struct {
	Stacks []*RawStack `json:"stacks"`
}

type RawStack struct {
	Name *string `json:"name"`

	Capabilities          []string        `json:"capabilities"`
	Parameters            StackParameters `json:"parameters"`
	Tags                  StackTags       `json:"tags"`
	TerminationProtection *bool           `json:"terminationProtection"`

	PolicyFile   *string `json:"policyFile"`
	TemplateFile *string `json:"templateFile"`
}
