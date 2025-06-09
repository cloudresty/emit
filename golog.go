package golog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

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

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	default:
		return "info"
	}
}

// ParseLogLevel parses a string into a LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info", "information":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}

// LogEntry represents a structured log entry for Kubernetes
type LogEntry struct {
	Timestamp string         `json:"timestamp"`
	Level     string         `json:"level"`
	Msg       string         `json:"msg"`
	Component string         `json:"component,omitempty"`
	Version   string         `json:"version,omitempty"`
	File      string         `json:"file,omitempty"`
	Line      int            `json:"line,omitempty"`
	Function  string         `json:"function,omitempty"`
	Fields    map[string]any `json:"fields,omitempty"`
}

// Logger represents the JSON logger
type Logger struct {
	level      LogLevel
	component  string
	version    string
	writer     io.Writer
	showCaller bool
	format     OutputFormat
}

// Global logger instance
var defaultLogger *Logger

// init initializes a default logger
func init() {
	defaultLogger = &Logger{
		level:      INFO,
		writer:     os.Stdout,
		showCaller: false,
		format:     JSON_FORMAT, // JSON is default
	}

	// Check environment variable for format override
	if logFormat := os.Getenv("GOLOG_FORMAT"); logFormat != "" {
		switch strings.ToLower(logFormat) {
		case "plain", "text", "console", "development", "dev":
			defaultLogger.format = PLAIN_FORMAT
		case "json", "production", "prod":
			defaultLogger.format = JSON_FORMAT
		default:
			// Invalid value, stick with JSON default
			defaultLogger.format = JSON_FORMAT
		}
	}

	// Also check for log level from environment
	if logLevel := os.Getenv("GOLOG_LEVEL"); logLevel != "" {
		defaultLogger.level = ParseLogLevel(logLevel)
	}

	// Check for caller information setting
	if showCaller := os.Getenv("GOLOG_SHOW_CALLER"); showCaller != "" {
		defaultLogger.showCaller = strings.ToLower(showCaller) == "true" || showCaller == "1"
	}
}

// SetComponent sets the component name for the default logger
func SetComponent(component string) {
	if defaultLogger != nil {
		defaultLogger.component = component
	}
}

// SetVersion sets the version for the default logger
func SetVersion(version string) {
	if defaultLogger != nil {
		defaultLogger.version = version
	}
}

// SetLevel sets the log level for the default logger
func SetLevel(level string) {
	if defaultLogger != nil {
		defaultLogger.level = ParseLogLevel(level)
	}
}

// SetShowCaller enables or disables caller information
func SetShowCaller(show bool) {
	if defaultLogger != nil {
		defaultLogger.showCaller = show
	}
}

// SetFormat sets the output format (JSON or Plain)
func SetFormat(format string) {
	if defaultLogger != nil {
		switch strings.ToLower(format) {
		case "plain", "text", "console":
			defaultLogger.format = PLAIN_FORMAT
		case "json":
			defaultLogger.format = JSON_FORMAT
		default:
			defaultLogger.format = JSON_FORMAT
		}
	}
}

// SetPlainFormat switches to plain text output for development
func SetPlainFormat() {
	SetFormat("plain")
}

// SetJSONFormat switches to JSON output for production
func SetJSONFormat() {
	SetFormat("json")
}

// log writes a log entry at the specified level
func (l *Logger) log(level LogLevel, message string, fields map[string]any) {
	if level < l.level {
		return
	}

	// Route to appropriate formatter based on format setting
	if l.format == PLAIN_FORMAT {
		l.logPlain(level, message, fields)
	} else {
		l.logJSON(level, message, fields)
	}
}

// logJSON writes a JSON formatted log entry
func (l *Logger) logJSON(level LogLevel, message string, fields map[string]any) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level.String(),
		Msg:       message,
	}

	if l.component != "" {
		entry.Component = l.component
	}

	if l.version != "" {
		entry.Version = l.version
	}

	if len(fields) > 0 {
		entry.Fields = fields
	}

	if l.showCaller {
		if pc, file, line, ok := runtime.Caller(4); ok {
			entry.File = file
			entry.Line = line
			if fn := runtime.FuncForPC(pc); fn != nil {
				entry.Function = fn.Name()
			}
		}
	}

	data, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple format if JSON marshaling fails
		fmt.Fprintf(l.writer, `{"timestamp":"%s","level":"error","msg":"Failed to marshal log entry: %v","component":"%s"}`+"\n",
			time.Now().UTC().Format(time.RFC3339Nano), err, l.component)
		return
	}

	fmt.Fprintln(l.writer, string(data))
}

// logPlain writes a plain text formatted log entry
func (l *Logger) logPlain(level LogLevel, message string, fields map[string]any) {
	severity := level.String()

	var colorCode string
	switch severity {
	case "info":
		colorCode = "\033[32m" // Green
	case "warn":
		colorCode = "\033[33m" // Yellow
	case "error":
		colorCode = "\033[31m" // Red
	case "debug":
		colorCode = "\033[34m" // Blue
	default:
		colorCode = ""
	}

	resetCode := "\033[0m" // Reset color

	if runtime.GOOS == "windows" {
		// Windows doesn't directly support ANSI escape codes
		colorCode = ""
		resetCode = ""
	}

	// Build the message with fields if present
	finalMessage := message
	if len(fields) > 0 {
		var fieldParts []string
		for k, v := range fields {
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, v))
		}
		finalMessage = fmt.Sprintf("%s [%s]", message, strings.Join(fieldParts, " "))
	}

	// Console output format:
	// {UTC TIME} | {LOGGING LEVEL} | {COMPONENT} {VERSION}: {MESSAGE}
	fmt.Fprintf(l.writer, "%s | %s%-7s%s | %s %s: %s\n",
		time.Now().UTC().Format("2006-01-02 15:04:05"),
		colorCode, severity, resetCode, l.component, l.version, finalMessage)
}

// Default logging functions - these automatically use the configured format

// Log is the default function that prints a log message in the configured format
func Log(level, message string, optionalParams ...string) {
	logLevel := ParseLogLevel(level)

	// Handle optional parameters
	if len(optionalParams) >= 1 && defaultLogger.component == "" {
		SetComponent(optionalParams[0])
	}
	if len(optionalParams) >= 2 && defaultLogger.version == "" {
		SetVersion(optionalParams[1])
	}

	defaultLogger.log(logLevel, message, nil)
}

// Info logs an info level message
func Info(message string, optionalParams ...string) {
	Log("info", message, optionalParams...)
}

// Warning logs a warning level message
func Warning(message string, optionalParams ...string) {
	Log("warn", message, optionalParams...)
}

// Error logs an error level message
func Error(message string, optionalParams ...string) {
	Log("error", message, optionalParams...)
}

// Debug logs a debug level message
func Debug(message string, optionalParams ...string) {
	Log("debug", message, optionalParams...)
}

// Rich logging functions with fields

// InfoWithFields logs an info message with additional fields
func InfoWithFields(message string, fields map[string]any) {
	if defaultLogger != nil {
		defaultLogger.log(INFO, message, fields)
	}
}

// WarnWithFields logs a warning message with additional fields
func WarnWithFields(message string, fields map[string]any) {
	if defaultLogger != nil {
		defaultLogger.log(WARN, message, fields)
	}
}

// ErrorWithFields logs an error message with additional fields
func ErrorWithFields(message string, fields map[string]any) {
	if defaultLogger != nil {
		defaultLogger.log(ERROR, message, fields)
	}
}

// DebugWithFields logs a debug message with additional fields
func DebugWithFields(message string, fields map[string]any) {
	if defaultLogger != nil {
		defaultLogger.log(DEBUG, message, fields)
	}
}

// Backward compatibility functions

// JSON function for backward compatibility - forces JSON output
func JSON(severity, message string, optionalParams ...string) {
	logLevel := ParseLogLevel(severity)

	// Handle optional parameters
	if len(optionalParams) >= 1 && defaultLogger.component == "" {
		SetComponent(optionalParams[0])
	}
	if len(optionalParams) >= 2 && defaultLogger.version == "" {
		SetVersion(optionalParams[1])
	}

	// Force JSON format for this call
	defaultLogger.logJSON(logLevel, message, nil)
}

// Plain function for backward compatibility - forces plain output
func Plain(severity, message string, optionalParams ...string) {
	logLevel := ParseLogLevel(severity)

	// Handle optional parameters
	if len(optionalParams) >= 1 && defaultLogger.component == "" {
		SetComponent(optionalParams[0])
	}
	if len(optionalParams) >= 2 && defaultLogger.version == "" {
		SetVersion(optionalParams[1])
	}

	// Force plain format for this call
	defaultLogger.logPlain(logLevel, message, nil)
}
