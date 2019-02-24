package config

import (
	"encoding/json"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v2"
)

type Result struct {
	Template  []byte
	LocalPath []byte
}

func TraverseNestedStacks(template []byte) ([]byte, error) {
	var magic interface{}
	json.Unmarshal(template, &magic)

	err := traverse(magic, make([]interface{}, 0))
	if err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(magic)
	if err != nil {
		return nil, err
	}

	return data, nil
}

var counter = 0

func traverse(model interface{}, path []interface{}) error {
	modelType := reflect.TypeOf(model)
	modelKind := modelType.Kind()
	modelValue := reflect.ValueOf(model)

	fmt.Println(path)

	switch modelKind {
	case reflect.Map:
		for _, key := range modelValue.MapKeys() {
			keyInterface := key.Interface()
			keyValue := modelValue.MapIndex(key)
			keyValueInterface := keyValue.Interface()

			switch keyValue.Elem().Kind() {
			case reflect.Map, reflect.Slice:
				err := traverse(keyValueInterface, append(path, keyInterface))
				if err != nil {
					return err
				}

			case reflect.String:
				if reflect.DeepEqual(keyInterface, "Type") &&
					reflect.DeepEqual(keyValueInterface, "AWS::CloudFormation::Stack") {
					fmt.Println("FOUND STACK")
				}

				str := processString(keyValueInterface.(string))
				modelValue.SetMapIndex(key, reflect.ValueOf(str))
			}
		}

	case reflect.Slice:
		length := modelValue.Len()
		for index := 0; index < length; index++ {
			indexValue := modelValue.Index(index)
			indexValueInterface := indexValue.Interface()

			switch indexValue.Elem().Kind() {
			case reflect.Map, reflect.Slice:
				err := traverse(indexValueInterface, append(path, index))
				if err != nil {
					return err
				}

			case reflect.String:
				str := processString(indexValueInterface.(string))
				modelValue.Index(index).Set(reflect.ValueOf(str))
			}
		}

	default:
		return fmt.Errorf("unrecognised kind '%+v'", modelKind)
	}

	return nil
}

func processString(str string) string {
	return fmt.Sprintf("%s::LOL", str)
}

var resourceProcessors []func(interface{}) error
var stringProcessors []func(string) (string, error)
