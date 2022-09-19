package logger

import "fmt"

// Fields Type to pass when we want to call WithFields for structured logging.
type Fields map[string]interface{}

const (
	// DebugLevel has verbose message.
	DebugLevel = "debug"
	// InfoLevel is default log level.
	InfoLevel = "info"
	// WarnLevel is for logging messages about possible issues.
	WarnLevel = "warn"
	// ErrorLevel is for logging errors.
	ErrorLevel = "error"
	// FatalLevel is for logging fatal messages. The system shutdown after logging the message.
	FatalLevel = "fatal"
)

// Configuration stores the config for the logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default.
type Configuration struct {
	EnableConsole     bool   `yaml:"enable console"`
	ConsoleJSONFormat bool   `yaml:"console json format"`
	ConsoleLevel      string `yaml:"console level"`
	EnableFile        bool   `yaml:"enable file"`
	FileJSONFormat    bool   `yaml:"file json format"`
	FileLevel         string `yaml:"file level"`
	FileLocation      string `yaml:"file location"`
}

// Default sets default values in config variables.
func (c *Configuration) Default() {
	c.EnableConsole = true
	c.ConsoleJSONFormat = true
	c.ConsoleLevel = DebugLevel
	c.EnableFile = true
	c.FileJSONFormat = true
	c.FileLevel = InfoLevel
	c.FileLocation = "./tmp/logs.log"
}

// InitLogger returns an instance of logger.
func InitLogger(config Configuration) error {
	logger, err := newZapLogger(config)
	if err != nil {
		return fmt.Errorf("failed to create new zap logger: %w", err)
	}

	global = logger

	return nil
}
