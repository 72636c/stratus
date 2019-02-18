package log

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/pretty"
)

var (
	ColourLogger   = new(colourLogger)
	StandardLogger = new(standardLogger)
)

type Logger interface {
	Data(model interface{})
	Title(format string, arguments ...interface{})
}

type colourLogger struct{}

func (logger *colourLogger) Data(model interface{}) {
	data, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		fmt.Printf("%+v\n", model)
		return
	}

	fmt.Printf("%s\n", pretty.Color(data, nil))
}

func (logger *colourLogger) Title(format string, arguments ...interface{}) {
	if len(arguments) == 0 {
		fmt.Println(format)
		return
	}

	fmt.Println(fmt.Sprintf(format, arguments...))
}

type standardLogger struct{}

func (logger *standardLogger) Data(model interface{}) {
	fmt.Println(model)
}

func (logger *standardLogger) Title(format string, arguments ...interface{}) {
	if len(arguments) == 0 {
		fmt.Println(format)
		return
	}

	fmt.Println(fmt.Sprintf(format, arguments...))
}
