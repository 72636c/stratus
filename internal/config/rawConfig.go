package config

type RawConfig struct {
	Stacks []*RawStack `json:"stacks"`
}

type RawStack struct {
	Name *string `json:"name"`

	Capabilities          []string           `json:"capabilities"`
	Parameters            RawStackParameters `json:"parameters"`
	Tags                  RawStackTags       `json:"tags"`
	TerminationProtection *bool              `json:"terminationProtection" yaml:"terminationProtection"`

	PolicyFile   *string `json:"policyFile" yaml:"policyFile"`
	TemplateFile *string `json:"templateFile" yaml:"templateFile"`
}

type RawStackParameters []*RawStackParameter

type RawStackParameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RawStackTags []*RawStackTag

type RawStackTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
