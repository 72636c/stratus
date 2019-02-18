package cli

import (
	"strings"

	"github.com/72636c/stratus/internal/log"
)

var (
	nameToLogger = map[string]log.Logger{
		"color":  log.ColourLogger,
		"colour": log.ColourLogger,
		"plain":  log.StandardLogger,
	}

	loggerNames = func() string {
		names := make([]string, 0)

		for name := range nameToLogger {
			names = append(names, name)
		}

		return strings.Join(names, "|")
	}()
)
