package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ttacon/chalk"
	"os"
)

func NewLogger() *logrus.Logger { //todo in further include to the prject instead of usual logger
	// Create a new instance of the logger
	logger := logrus.New()

	// Set the output to stdout
	logger.SetOutput(os.Stdout)

	// Set the log level to Debug for demonstration purposes
	logger.SetLevel(logrus.DebugLevel)

	// Add a hook to customize log message colors
	logger.AddHook(NewCustomColorHook())

	// Log some example messages
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
	return logger
}

// CustomColorHook is a logrus hook to customize log message colors
type CustomColorHook struct{}

// NewCustomColorHook creates a new instance of the CustomColorHook
func NewCustomColorHook() *CustomColorHook {
	return &CustomColorHook{}
}

// Levels returns the log levels for the hook to fire
func (hook *CustomColorHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is called when a log event is fired
func (hook *CustomColorHook) Fire(entry *logrus.Entry) error {
	switch entry.Level {
	case logrus.InfoLevel:
		entry.Message = fmt.Sprintf("%s %s", chalk.Blue.Color("[INFO]"), entry.Message)
	case logrus.WarnLevel:
		entry.Message = fmt.Sprintf("%s %s", chalk.Yellow.Color("[WARNING]"), entry.Message)
	case logrus.ErrorLevel:
		entry.Message = fmt.Sprintf("%s %s", chalk.Red.Color("[ERROR]"), entry.Message)
	}

	return nil
}
