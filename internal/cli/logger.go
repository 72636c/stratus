package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/pretty"
)

var (
	loggers = map[string]Logger{
		"color":  ColourLogger,
		"colour": ColourLogger,
		"plain":  PlainLogger,
	}

	loggerNames = func() string {
		names := make([]string, 0)

		for name := range loggers {
			names = append(names, name)
		}

		return strings.Join(names, "|")
	}()
)

type Logger func(message interface{}, ok bool)

func ColourLogger(message interface{}, ok bool) {
	if !ok {
		fmt.Printf("\n")
		return
	}

	str, ok := message.(string)
	if ok {
		fmt.Printf("\n\n%s\n", str)
		return
	}

	data, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		fmt.Printf("\n%+v", message)
		return
	}

	fmt.Printf("\n%s", pretty.Color(data, nil))
}

func PlainLogger(message interface{}, ok bool) {
	if !ok {
		fmt.Printf("\n")
		return
	}

	str, ok := message.(string)
	if ok {
		fmt.Printf("\n\n%s\n", str)
		return
	}

	fmt.Printf("\n%+v", message)
}
