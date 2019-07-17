package log

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"
)

var (
	// not implemented; ideally roll a YAML encoder with colour terminal output
	ColourLogger = StandardLogger

	StandardLogger = new(standardLogger)
)

type Logger interface {
	Data(model interface{})
	Title(format string, arguments ...interface{})
}

type standardLogger struct{}

func (logger *standardLogger) Data(model interface{}) {
	if reflect.TypeOf(model).Kind() == reflect.String {
		fmt.Printf("%+v\n", model)
		return
	}

	// round trip to prevent yaml.v2 from lowercasing struct field names

	jsonData, err := json.Marshal(model)
	if err != nil {
		fmt.Printf("%+v\n", model)
		return
	}

	var jsonModel interface{}
	err = json.Unmarshal(jsonData, &jsonModel)
	if err != nil {
		fmt.Printf("%+v\n", model)
		return
	}

	data, err := yaml.Marshal(jsonModel)
	if err != nil {
		fmt.Printf("%+v\n", model)
		return
	}

	fmt.Printf("%s", data)
}

func (logger *standardLogger) Title(format string, arguments ...interface{}) {
	fmt.Println()

	title := format
	if len(arguments) != 0 {
		title = fmt.Sprintf(format, arguments...)
	}

	fmt.Println(aurora.Bold(title))

	fmt.Println(strings.Repeat("─", len(title)))
}
