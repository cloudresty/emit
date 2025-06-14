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
		Message:   message,
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
		fmt.Fprintf(l.writer, `{"timestamp":"%s","level":"error","message":"Failed to marshal log entry: %v","component":"%s"}`+"\n",
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

// buildSimpleJSONUltraFast - Ultra-fast JSON builder for simple messages
func (l *Logger) buildSimpleJSONUltraFast(buf []byte, level LogLevel, message string) int {
	timestamp := GetUltraFastTimestamp()
	levelStr := level.StringFast()

	pos := 0

	// Check buffer space as we build to prevent overflow
	if pos+len(`{"timestamp":"`) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], `{"timestamp":"`)

	if pos+len(timestamp) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], timestamp)

	if pos+len(`","level":"`) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], `","level":"`)

	if pos+len(levelStr) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], levelStr)

	if pos+len(`","message":"`) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], `","message":"`)

	if pos+len(message) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], message)

	if pos+len(`"`) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], `"`)

	if l.component != "" {
		if pos+len(`,"component":"`) >= len(buf) {
			return len(buf)
		}
		pos += copy(buf[pos:], `,"component":"`)

		if pos+len(l.component) >= len(buf) {
			return len(buf)
		}
		pos += copy(buf[pos:], l.component)

		if pos+len(`"`) >= len(buf) {
			return len(buf)
		}
		pos += copy(buf[pos:], `"`)
	}

	if l.version != "" {
		if pos+len(`,"version":"`) >= len(buf) {
			return len(buf)
		}
		pos += copy(buf[pos:], `,"version":"`)

		if pos+len(l.version) >= len(buf) {
			return len(buf)
		}
		pos += copy(buf[pos:], l.version)

		if pos+len(`"`) >= len(buf) {
			return len(buf)
		}
		pos += copy(buf[pos:], `"`)
	}

	if pos+len("}\n") >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], "}\n")
	return pos
}

// buildSimplePlainUltraFast - Ultra-fast plain text builder for simple messages
func (l *Logger) buildSimplePlainUltraFast(buf []byte, level LogLevel, message string) int {
	timestamp := GetUltraFastTimestamp()
	levelStr := level.StringFast()

	pos := 0

	// Check buffer space as we build to prevent overflow
	if pos+19 >= len(buf) {
		return len(buf)
	} // timestamp[:19]
	pos += copy(buf[pos:], timestamp[:19])

	if pos+3 >= len(buf) {
		return len(buf)
	} // " | "
	pos += copy(buf[pos:], " | ")

	if pos+len(levelStr) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], levelStr)

	// Pad level to 7 characters for alignment
	for len(levelStr) < 7 {
		if pos+1 >= len(buf) {
			return len(buf)
		}
		pos += copy(buf[pos:], " ")
		levelStr += " "
	}

	if pos+3 >= len(buf) {
		return len(buf)
	} // " | "
	pos += copy(buf[pos:], " | ")

	if pos+len(l.component) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], l.component)

	if pos+1 >= len(buf) {
		return len(buf)
	} // " "
	pos += copy(buf[pos:], " ")

	if pos+len(l.version) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], l.version)

	if pos+2 >= len(buf) {
		return len(buf)
	} // ": "
	pos += copy(buf[pos:], ": ")

	if pos+len(message) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], message)

	if pos+1 >= len(buf) {
		return len(buf)
	} // "\n"
	pos += copy(buf[pos:], "\n")

	return pos
}
