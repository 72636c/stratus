package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ssm"
	"gopkg.in/yaml.v2"
)

var (
	extensionToUnmarshal = map[string]func([]byte, interface{}) error{
		".json": json.Unmarshal,
		".yaml": yaml.UnmarshalStrict,
		".yml":  yaml.UnmarshalStrict,
	}

	mapperStore = NewMapperStore()
)

type MapperStore struct {
	sync.RWMutex
	fromPrefix map[string]Mapper
}

func NewMapperStore() *MapperStore {
	return &MapperStore{
		RWMutex:    sync.RWMutex{},
		fromPrefix: make(map[string]Mapper),
	}
}

func (store *MapperStore) Get(prefix string) (Mapper, bool) {
	store.RLock()
	mapper, ok := store.fromPrefix[prefix]
	store.RUnlock()

	return mapper, ok
}
func (store *MapperStore) Set(prefix string, mapper Mapper) {
	store.Lock()
	store.fromPrefix[prefix] = mapper
	store.Unlock()
}

type Mapper func(placeholder string) (string, error)

func envMapper(placeholder string) (string, error) {
	resolved, ok := os.LookupEnv(placeholder)
	if !ok {
		return "", fmt.Errorf("environment variable '%s' not set", placeholder)
	}

	return resolved, nil
}

func newAWSMapper(provider client.ConfigProvider) Mapper {
	client := ssm.New(provider)

	return func(placeholder string) (string, error) {
		if !strings.HasPrefix(placeholder, "ssm:parameter:") {
			return "", fmt.Errorf("unsupported AWS placeholder '%s'", placeholder)
		}

		name := strings.TrimPrefix(placeholder, "ssm:parameter:")

		input := &ssm.GetParameterInput{
			Name:           aws.String(name),
			WithDecryption: aws.Bool(false),
		}

		output, err := client.GetParameter(input)
		if err != nil {
			return "", err
		}

		return *output.Parameter.Value, nil
	}
}

func Init(provider client.ConfigProvider) {
	mapperStore.Set("env", envMapper)

	if provider != nil {
		mapperStore.Set("aws", newAWSMapper(provider))
	}
}

func Unmarshal(
	extension string,
	data []byte,
	model interface{},
) error {
	unmarshal, ok := extensionToUnmarshal[extension]
	if !ok {
		return fmt.Errorf("unsupported file extension '%s'", extension)
	}

	return unmarshal(data, model)
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

			mapper, ok := mapperStore.Get(prefix)
			if !ok {
				return "", fmt.Errorf("unsupported placeholder '%s'", str)
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
	type boolAlias Bool

	resolved, err := Resolve(string(data))
	if err != nil {
		return err
	}

	resolved = strings.Trim(resolved, `"`)

	return json.Unmarshal([]byte(resolved), (*boolAlias)(bit))
}

func (bit *Bool) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type boolAlias Bool

	var data string

	err := unmarshal(&data)
	if err != nil {
		return err
	}

	resolved, err := Resolve(data)
	if err != nil {
		return err
	}

	resolved = strings.Trim(resolved, `"`)

	return yaml.UnmarshalStrict([]byte(resolved), (*boolAlias)(bit))
}

type String string

func (str *String) String() string {
	return string(*str)
}

func (str *String) UnmarshalJSON(data []byte) error {
	type stringAlias String

	resolved, err := Resolve(string(data))
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(resolved), (*stringAlias)(str))
}

func (str *String) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type stringAlias String

	var data string

	err := unmarshal(&data)
	if err != nil {
		return err
	}

	resolved, err := Resolve(data)
	if err != nil {
		return err
	}

	return yaml.UnmarshalStrict([]byte(resolved), (*stringAlias)(str))
}
