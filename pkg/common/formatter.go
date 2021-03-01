package common

import (
	"bytes"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// LogFormatter gives a more user friendly log
type LogFormatter struct {
	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string
}

// Format formats a log entry during output
func (f *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	// output buffer
	b := &bytes.Buffer{}

	// log level with color
	fmt.Fprintf(b, "\x1b[%dm%-7s\x1b[0m", levelColor, entry.Level.String())
	b.WriteByte('\t')

	var callingFile string
	i := strings.LastIndex(entry.Caller.File, "/")
	if i != -1 {
		callingFile = entry.Caller.File[i+1:]
	} else {
		callingFile = entry.Caller.File
	}

	if entry.HasCaller() {
		fmt.Fprintf(b, "%-30s", fmt.Sprintf("%s:%d", callingFile, entry.Caller.Line))
		b.WriteByte('\t')
	}

	// formatted timestamp
	fmt.Fprintf(b, "[%s]", entry.Time.Format(f.TimestampFormat))
	b.WriteByte('\t')

	// message
	fmt.Fprintf(b, "%-50s", entry.Message)
	b.WriteByte('\t')

	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for f, v := range entry.Data {
			fields = append(fields, fmt.Sprintf("'%s' : '%v'", f, v))
		}
		fmt.Fprintf(b, "extraFields: [ %s ]", strings.Join(fields, ", "))
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level log.Level) int {
	switch level {
	case log.DebugLevel:
		return colorGray
	case log.WarnLevel:
		return colorYellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
