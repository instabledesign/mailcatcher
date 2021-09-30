package log

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gol4ng/logger"
)

// DefaultFormatter is the default Entry formatter
type DefaultFormatter struct {
	colored       bool
	verboseLevel  logger.Level
	prefixMessage string
}

// Format will return Entry as string
func (n *DefaultFormatter) Format(entry logger.Entry) string {
	builder := &strings.Builder{}

	builder.WriteString(time.Now().Format("2006/01/02 15:04:05.000 "))
	if n.colored {
		switch entry.Level {
		case logger.DebugLevel:
			builder.WriteString("\x1b[1;36m")
		case logger.InfoLevel:
			builder.WriteString("\x1b[1;32m")
		case logger.NoticeLevel:
			builder.WriteString("\x1b[1;34m")
		case logger.WarningLevel:
			builder.WriteString("\x1b[1;33m")
		case logger.ErrorLevel:
			builder.WriteString("\x1b[1;31m")
		case logger.CriticalLevel:
			builder.WriteString("\x1b[1;30;47m")
		case logger.AlertLevel:
			builder.WriteString("\x1b[1;30;43m")
		case logger.EmergencyLevel:
			builder.WriteString("\x1b[1;37;41m")
		}
	}
	builder.WriteString("[ ")
	builder.WriteString(strings.ToUpper(entry.Level.String()))
	builder.WriteString(" ]\x1b[m ")
	builder.WriteString(n.prefixMessage)
	builder.WriteString(entry.Message)
	if entry.Level <= n.verboseLevel && entry.Context != nil {
		builder.WriteString(" ")
		ContextToJSON(entry.Context, builder, n.colored)
	}

	return builder.String()
}

// NewDefaultFormatter will create a new DefaultFormatter
func NewDefaultFormatter(colored bool, verboseLevel logger.Level) *DefaultFormatter {
	return &DefaultFormatter{
		colored:      colored,
		verboseLevel: verboseLevel,
	}
}

func NewPrefixFormatter(prefix string, colored bool, verboseLevel logger.Level) *DefaultFormatter {
	return &DefaultFormatter{
		colored:       colored,
		verboseLevel:  verboseLevel,
		prefixMessage: prefix,
	}
}

func ContextToJSON(context *logger.Context, builder *strings.Builder, colored bool) {
	if context == nil || len(*context) == 0 {
		if colored {
			builder.WriteString("\x1b[3;40mnull\x1b[m")
			return
		}
		builder.WriteString("null")
		return
	}

	builder.WriteString("{")
	i := 0
	for name, field := range *context {
		if i != 0 {
			builder.WriteRune(',')
		}
		builder.WriteRune('"')
		if colored {
			builder.WriteString("\x1b[1;34m")
		}
		builder.WriteString(name)
		if colored {
			builder.WriteString("\x1b[m")
		}
		builder.WriteString("\":")
		d, _ := json.Marshal(field.Value)
		builder.WriteString(string(d))
		i++
	}
	builder.WriteString("}")
}
