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

func resolve(data string) (string, error) {
	// TODO: resolve placeholders

	return data, nil
}

type Bool bool

func (bit *Bool) Bool() bool {
	return bool(*bit)
}

func (bit *Bool) UnmarshalJSON(data []byte) error {
	type alias Bool

	resolved, err := resolve(string(data))
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(resolved), (*alias)(bit))
}

func (bit *Bool) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias Bool

	var data string
	err := unmarshal(&data)
	if err != nil {
		return err
	}

	resolved, err := resolve(data)
	if err != nil {
		return err
	}

	return yaml.UnmarshalStrict([]byte(resolved), (*alias)(bit))
}

type String string

func (str *String) String() string {
	return string(*str)
}

func (str *String) UnmarshalJSON(data []byte) error {
	type alias String

	resolved, err := resolve(string(data))
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(resolved), (*alias)(str))
}

func (str *String) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias String

	var data string
	err := unmarshal(&data)
	if err != nil {
		return err
	}

	resolved, err := resolve(data)
	if err != nil {
		return err
	}

	return yaml.UnmarshalStrict([]byte(resolved), (*alias)(str))
}
