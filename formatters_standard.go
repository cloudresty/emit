package emit

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// logJSON writes a JSON formatted log entry
func (l *Logger) logJSON(level LogLevel, message string, fields map[string]any) {
	entry := LogEntry{
		Timestamp: GetUltraFastTimestamp(),
		Level:     level.StringFast(),
		Msg:       message,
	}

	if l.component != "" {
		entry.Component = l.component
	}

	if l.version != "" {
		entry.Version = l.version
	}

	if len(fields) > 0 {
		entry.Fields = l.maskSensitiveFieldsFast(fields)
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
			GetUltraFastTimestamp(), err, l.component)
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

	// Build the message with fields if present (with masking)
	finalMessage := message
	if len(fields) > 0 {
		maskedFields := l.maskSensitiveFieldsFast(fields)
		var fieldParts []string
		for k, v := range maskedFields {
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, v))
		}
		finalMessage = fmt.Sprintf("%s [%s]", message, strings.Join(fieldParts, " "))
	}

	// Console output format:
	// {UTC TIME} | {LOGGING LEVEL} | {COMPONENT} {VERSION}: {MESSAGE}
	fmt.Fprintf(l.writer, "%s | %s%-7s%s | %s %s: %s\n",
		GetUltraFastTimestamp()[:19],
		colorCode, severity, resetCode, l.component, l.version, finalMessage)
}

// These bypass all field processing and security checks for simple message-only logging

// logJSONSimple writes a simple JSON log entry without fields (optimized)
func (l *Logger) logJSONSimple(level LogLevel, message string) {
	entry := LogEntry{
		Timestamp: GetUltraFastTimestamp(),
		Level:     level.String(),
		Msg:       message,
	}

	if l.component != "" {
		entry.Component = l.component
	}

	if l.version != "" {
		entry.Version = l.version
	}

	// No fields processing - this is the key optimization

	if l.showCaller {
		if pc, file, line, ok := runtime.Caller(5); ok { // One level deeper due to logSimple
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
			GetUltraFastTimestamp(), err, l.component)
		return
	}

	fmt.Fprintln(l.writer, string(data))
}

// logPlainSimple writes a simple plain text log entry without fields (optimized)
func (l *Logger) logPlainSimple(level LogLevel, message string) {
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

	// Simple message without any field processing
	// Console output format:
	// {UTC TIME} | {LOGGING LEVEL} | {COMPONENT} {VERSION}: {MESSAGE}
	fmt.Fprintf(l.writer, "%s | %s%-7s%s | %s %s: %s\n",
		GetUltraFastTimestamp()[:19], // PHASE 3: Use optimized timestamp, truncate to YYYY-MM-DD HH:MM:SS
		colorCode, severity, resetCode, l.component, l.version, message)
}
