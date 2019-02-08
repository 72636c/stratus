package config

type RawConfig struct {
	Defaults RawDefaults `json:"defaults"`
	Stacks   []*RawStack `json:"stacks"`
}

type RawDefaults struct {
	UploadArtefacts String `json:"uploadArtefacts" yaml:"uploadArtefacts"`
}

type RawStack struct {
	Name String `json:"name"`

	Capabilities          RawStackCapabilities `json:"capabilities"`
	Parameters            RawStackParameters   `json:"parameters"`
	Tags                  RawStackTags         `json:"tags"`
	TerminationProtection Bool                 `json:"terminationProtection" yaml:"terminationProtection"`

	PolicyFile   String `json:"policyFile" yaml:"policyFile"`
	TemplateFile String `json:"templateFile" yaml:"templateFile"`
}

type RawStackCapabilities []String

type RawStackParameters []*RawStackParameter

type RawStackParameter struct {
	Key   String `json:"key"`
	Value String `json:"value"`
}

type RawStackTags []*RawStackTag

type RawStackTag struct {
	Key   String `json:"key"`
	Value String `json:"value"`
}
