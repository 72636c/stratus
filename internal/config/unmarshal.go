package config

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

var (
	extensionToUnmarshal = map[string]func([]byte, interface{}) error{
		".json": json.Unmarshal,
		".yaml": yaml.UnmarshalStrict,
		".yml":  yaml.UnmarshalStrict,
	}

	prefixToMapper = map[string]Mapper{}
)

type Mapper func(placeholder string, dest interface{}) error

func Unmarshal(extension string, data []byte, dest interface{}) error {
	unmarshal, ok := extensionToUnmarshal[extension]
	if !ok {
		return fmt.Errorf("unrecognised file extension '%s'", extension)
	}

	// TODO: initialise prefixToMapper?

	return unmarshal(data, dest)
}
