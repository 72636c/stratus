package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"
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
	if str, ok := model.(string); ok {
		fmt.Println(str)
		return
	}

	str, err := encodeYAML(model)
	if err != nil {
		fmt.Printf("%+v\n", model)
		return
	}

	err = quick.Highlight(os.Stdout, str, "yaml", "terminal256", "pygments")
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
