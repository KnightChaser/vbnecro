package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// CustomFormatter formats log entries in the desired style.
type CustomFormatter struct{}

// Format builds the final log line.
// The format is: [timestamp] [LEVEL] message
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Format the timestamp with microsecond precision.
	timestamp := entry.Time.Format("2006/01/02 15:04:05.000000")
	// Convert level to uppercase.
	level := strings.ToUpper(entry.Level.String())

	// Define ANSI color codes for various levels.
	var levelColor string
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = "\x1b[36m" // Cyan
	case logrus.InfoLevel:
		levelColor = "\x1b[32m" // Green
	case logrus.WarnLevel:
		levelColor = "\x1b[33m" // Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = "\x1b[31m" // Red
	default:
		levelColor = ""
	}
	reset := "\x1b[0m"

	// Wrap the level in its color if available.
	if levelColor != "" {
		level = levelColor + level + reset
	}

	// Compose the final log line.
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, entry.Message)
	return []byte(logLine), nil
}

func init() {
	// Set logrus to output to STDOUT.
	logrus.SetOutput(os.Stdout)

	// Set our custom formatter.
	logrus.SetFormatter(&CustomFormatter{})

	// Optionally, set the log level.
	logrus.SetLevel(logrus.DebugLevel)
}
