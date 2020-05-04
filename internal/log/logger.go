package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/efekarakus/termcolor"
	"github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"
)

var (
	ColourLogger   = newColourLogger()
	StandardLogger = new(standardLogger)
)

func detectFormatter() string {
	switch level := termcolor.SupportLevel(os.Stdout); level {
	case termcolor.Level16M:
		return "terminal16m"
	case termcolor.Level256:
		return "terminal256"
	case termcolor.LevelBasic:
		return "terminal"
	default:
		return ""
	}
}

type Logger interface {
	Data(model interface{})
	Title(format string, arguments ...interface{})
}

type colourLogger struct {
	formatter string
}

func newColourLogger() *colourLogger {
	return &colourLogger{
		formatter: detectFormatter(),
	}
}

func (logger *colourLogger) Data(model interface{}) {
	if str, ok := model.(string); ok {
		fmt.Println(str)
		return
	}

	str, err := encodeYAML(model)
	if err != nil {
		fmt.Printf("%+v\n", model)
		return
	}

	err = quick.Highlight(os.Stdout, str, "yaml", logger.formatter, "arduino")
	if err != nil {
		fmt.Printf("%s", str)
	}
}

func (logger *colourLogger) Title(format string, arguments ...interface{}) {
	title := formatString(format, arguments...)

	fmt.Printf("\n%s\n%s\n", aurora.Bold(title), generateLine(title))
}

type standardLogger struct{}

func (logger *standardLogger) Data(model interface{}) {
	if str, ok := model.(string); ok {
		fmt.Println(str)
		return
	}

	str, err := encodeYAML(model)
	if err != nil {
		fmt.Printf("%+v\n", model)
		return
	}

	fmt.Printf("%s", str)
}

func (logger *standardLogger) Title(format string, arguments ...interface{}) {
	title := formatString(format, arguments...)

	fmt.Printf("\n%s\n%s\n", title, generateLine(title))
}

func formatString(format string, arguments ...interface{}) string {
	if len(arguments) == 0 {
		return format
	}

	return fmt.Sprintf(format, arguments...)
}

func generateLine(str string) string {
	return strings.Repeat("─", len(str))
}

// encodeYAML performs a JSON codec round trip to prevent yaml.Marshal from
// lowercasing struct field names.
func encodeYAML(model interface{}) (string, error) {
	jsonData, err := json.Marshal(model)
	if err != nil {
		return "", err
	}

	var jsonModel interface{}
	err = json.Unmarshal(jsonData, &jsonModel)
	if err != nil {
		return "", err
	}

	data, err := yaml.Marshal(jsonModel)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
