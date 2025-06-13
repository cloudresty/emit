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

// logSimpleUltraFast - Specialized simple message logger
func (l *Logger) logSimpleUltraFast(level LogLevel, message string) {
	// Use smallest possible stack buffer for simple messages
	var stackBuf [128]byte
	var pos int

	if l.format == JSON_FORMAT {
		pos = l.buildSimpleJSONUltraFast(stackBuf[:], level, message)
	} else {
		pos = l.buildSimplePlainUltraFast(stackBuf[:], level, message)
	}

	// Single write operation - most critical optimization
	l.writer.Write(stackBuf[:pos])
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
