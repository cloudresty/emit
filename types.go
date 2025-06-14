package emit

import "io"

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// OutputFormat represents the output format type
type OutputFormat int

const (
	JSON_FORMAT OutputFormat = iota
	PLAIN_FORMAT
)

// SensitiveDataMode represents how to handle sensitive data
type SensitiveDataMode int

const (
	MASK_SENSITIVE SensitiveDataMode = iota // Default: mask sensitive data
	SHOW_SENSITIVE                          // Show sensitive data (not recommended for production)
)

// PIIDataMode represents how to handle PII data
type PIIDataMode int

const (
	MASK_PII PIIDataMode = iota // Default: mask PII data
	SHOW_PII                    // Show PII data (not recommended for production)
)

// LogEntry represents a structured log entry for Kubernetes
type LogEntry struct {
	Timestamp string         `json:"timestamp"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Component string         `json:"component,omitempty"`
	Version   string         `json:"version,omitempty"`
	File      string         `json:"file,omitempty"`
	Line      int            `json:"line,omitempty"`
	Function  string         `json:"function,omitempty"`
	Fields    map[string]any `json:"fields,omitempty"`
}

// Logger represents the JSON logger
type Logger struct {
	level           LogLevel
	component       string
	version         string
	writer          io.Writer
	showCaller      bool
	format          OutputFormat
	sensitiveMode   SensitiveDataMode
	piiMode         PIIDataMode
	sensitiveFields []string
	piiFields       []string
	maskString      string
	piiMaskString   string
}
