package emit

import "os"

// Global logger instance
var defaultLogger *Logger

// init initializes a default logger
func init() {
	defaultLogger = &Logger{
		level:           INFO,
		writer:          os.Stdout,
		showCaller:      false,
		format:          JSON_FORMAT,    // JSON is default
		sensitiveMode:   MASK_SENSITIVE, // Mask sensitive data by default
		piiMode:         MASK_PII,       // Mask PII data by default
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
	}

	// Initialize from environment variables
	initFromEnvironment()
}

// log writes a log entry at the specified level
func (l *Logger) log(level LogLevel, message string, fields map[string]any) {
	if level < l.level {
		return
	}

	// Ultra-fast path for simple messages (no fields) - OPTIMIZED FOR SPEED
	if len(fields) == 0 {
		l.logSimpleUltraFast(level, message)
		return
	}

	// Route to appropriate formatter based on format setting and field complexity
	if l.format == PLAIN_FORMAT {
		l.logPlain(level, message, fields)
	} else {
		// JSON format
		l.logJSON(level, message, fields)
	}
}

// logSimpleUltraFast - Specialized simple message logger with dynamic buffer
func (l *Logger) logSimpleUltraFast(level LogLevel, message string) {
	// Start with small optimal stack buffer for most common cases
	var stackBuf [128]byte
	var pos int
	var buf []byte = stackBuf[:]

	// First attempt with stack buffer
	if l.format == JSON_FORMAT {
		pos = l.buildSimpleJSONUltraFast(buf, level, message)
	} else {
		pos = l.buildSimplePlainUltraFast(buf, level, message)
	}
	// If buffer overflow detected, use dynamic allocation
	if pos >= len(buf) {
		// Estimate needed size based on format
		var estimatedSize int
		if l.format == JSON_FORMAT {
			estimatedSize = l.estimateJSONSize(level, message)
		} else {
			estimatedSize = l.estimatePlainSize(level, message)
		}
		dynamicBuf := make([]byte, estimatedSize)

		if l.format == JSON_FORMAT {
			pos = l.buildSimpleJSONUltraFast(dynamicBuf, level, message)
		} else {
			pos = l.buildSimplePlainUltraFast(dynamicBuf, level, message)
		}

		// Final safety check - if still overflows, fallback to safe method
		if pos >= len(dynamicBuf) {
			if l.format == JSON_FORMAT {
				l.logJSON(level, message, nil)
			} else {
				l.logPlain(level, message, nil)
			}
			return
		}

		buf = dynamicBuf
	}

	// Single write operation - most critical optimization
	l.writer.Write(buf[:pos])
}

// InfoStructured logs at INFO level with structured fields optimization
func InfoStructured(message string, fields ...ZField) {
	defaultLogger.InfoStructured(message, fields...)
}

func (l *Logger) InfoStructured(message string, fields ...ZField) {
	l.logStructuredFieldsRoute(INFO, message, fields...)
}

// DebugStructured logs at DEBUG level with structured fields optimization
func DebugStructured(message string, fields ...ZField) {
	defaultLogger.DebugStructured(message, fields...)
}

func (l *Logger) DebugStructured(message string, fields ...ZField) {
	l.logStructuredFieldsRoute(DEBUG, message, fields...)
}

// WarnStructured logs at WARN level with structured fields optimization
func WarnStructured(message string, fields ...ZField) {
	defaultLogger.WarnStructured(message, fields...)
}

func (l *Logger) WarnStructured(message string, fields ...ZField) {
	l.logStructuredFieldsRoute(WARN, message, fields...)
}

// ErrorStructured logs at ERROR level with structured fields optimization
func ErrorStructured(message string, fields ...ZField) {
	defaultLogger.ErrorStructured(message, fields...)
}

func (l *Logger) ErrorStructured(message string, fields ...ZField) {
	l.logStructuredFieldsRoute(ERROR, message, fields...)
}

// estimateJSONSize calculates the approximate size needed for JSON output
func (l *Logger) estimateJSONSize(level LogLevel, message string) int {
	// Base JSON structure: {"timestamp":"","level":"","message":""}
	baseSize := 50

	// Timestamp: ISO 8601 format is ~24 characters
	timestampSize := 30

	// Level: debug/info/warn/error (max ~5 chars)
	levelSize := 10

	// Message length
	messageSize := len(message)

	// Component field if present: ,"component":"value"
	componentSize := 0
	if l.component != "" {
		componentSize = 15 + len(l.component) // ,"component":"" + value
	}

	// Version field if present: ,"version":"value"
	versionSize := 0
	if l.version != "" {
		versionSize = 13 + len(l.version) // ,"version":"" + value
	}

	// Calculate total with 25% safety buffer
	totalSize := baseSize + timestampSize + levelSize + messageSize + componentSize + versionSize
	return totalSize + (totalSize / 4) // Add 25% buffer
}

// estimatePlainSize calculates the approximate size needed for plain text output
func (l *Logger) estimatePlainSize(level LogLevel, message string) int {
	// Timestamp: 19 characters (YYYY-MM-DD HH:MM:SS)
	timestampSize := 25

	// Separators: " | " + " | " + ": " + "\n" = ~10 chars
	separatorSize := 15

	// Level: debug/info/warn/error (max ~7 chars with padding)
	levelSize := 10

	// Message length
	messageSize := len(message)

	// Component length
	componentSize := len(l.component)

	// Version length
	versionSize := len(l.version)

	// Calculate total with 25% safety buffer
	totalSize := timestampSize + separatorSize + levelSize + messageSize + componentSize + versionSize
	return totalSize + (totalSize / 4) // Add 25% buffer
}
