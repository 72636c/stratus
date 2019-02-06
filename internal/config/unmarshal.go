package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	extensionToUnmarshal = map[string]func([]byte, interface{}) error{
		".json": json.Unmarshal,
		".yaml": yaml.UnmarshalStrict,
		".yml":  yaml.UnmarshalStrict,
	}

	prefixToMapper = map[string]Mapper{
		"env": envMapper,
	}
)

type Mapper func(placeholder string) (string, error)

func envMapper(placeholder string) (string, error) {
	resolved, ok := os.LookupEnv(placeholder)
	if !ok {
		return "", fmt.Errorf("environment variable '%s' not set", placeholder)
	}

	return resolved, nil
}

func Unmarshal(extension string, data []byte, dest interface{}) error {
	unmarshal, ok := extensionToUnmarshal[extension]
	if !ok {
		return fmt.Errorf("unrecognised file extension '%s'", extension)
	}

	return unmarshal(data, dest)
}

type ResolveStack []*strings.Builder

func NewResolveStack() ResolveStack {
	return ResolveStack{new(strings.Builder)}
}

func (stack ResolveStack) Peek() *strings.Builder {
	return stack[len(stack)-1]
}

func (stack ResolveStack) Pop() (ResolveStack, string) {
	item := stack.Peek().String()
	return stack[:len(stack)-1], item
}

func (stack ResolveStack) Push() ResolveStack {
	return append(stack, new(strings.Builder))
}

func (stack ResolveStack) String() string {
	var builder strings.Builder

	for len(stack) > 0 {
		var str string
		stack, str = stack.Pop()
		builder.WriteString(str)
	}

	return builder.String()
}

func Resolve(data string) (string, error) {
	skip := false
	stack := NewResolveStack()

	for index, token := range data {
		if skip {
			skip = false
			continue
		}

		switch token {
		case '{':
			if index+1 == len(data) || data[index+1] != '{' {
				stack.Peek().WriteRune('{')
				continue
			}

			stack = stack.Push()

			skip = true

		case '}':
			if index+1 == len(data) || data[index+1] != '}' {
				stack.Peek().WriteRune('}')
				continue
			}

			var str string

			stack, str = stack.Pop()

			slice := strings.SplitN(str, ":", 2)
			if len(slice) < 2 {
				return "", fmt.Errorf("malformed placeholder '%s'", str)
			}

			prefix := slice[0]
			suffix := slice[1]

			mapper, ok := prefixToMapper[prefix]
			if !ok {
				return "", fmt.Errorf("unrecognised placeholder '%s'", str)
			}

			resolved, err := mapper(suffix)
			if err != nil {
				return "", err
			}

			stack.Peek().WriteString(resolved)

			skip = true

		default:
			stack.Peek().WriteRune(token)
		}
	}

	return stack.String(), nil
}

type Bool bool

func (bit *Bool) Bool() bool {
	return bool(*bit)
}

func (bit *Bool) UnmarshalJSON(data []byte) error {
	type alias Bool

	resolved, err := Resolve(string(data))
	if err != nil {
		return err
	}

	resolved = strings.Replace(resolved, `"`, "", -1)

	return json.Unmarshal([]byte(resolved), (*alias)(bit))
}

func (bit *Bool) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias Bool

	var data string
	err := unmarshal(&data)
	if err != nil {
		return err
	}

	resolved, err := Resolve(data)
	if err != nil {
		return err
	}

	resolved = strings.Replace(resolved, `"`, "", -1)

	return yaml.UnmarshalStrict([]byte(resolved), (*alias)(bit))
}

type String string

func (str *String) String() string {
	return string(*str)
}

func (str *String) UnmarshalJSON(data []byte) error {
	type alias String

	resolved, err := Resolve(string(data))
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

	resolved, err := Resolve(data)
	if err != nil {
		return err
	}

	return yaml.UnmarshalStrict([]byte(resolved), (*alias)(str))
}
