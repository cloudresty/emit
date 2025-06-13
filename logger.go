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

	// This bypasses all security processing overhead for simple logging
	if len(fields) == 0 {
		l.logSimple(level, message)
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

// logSimple writes a simple log entry without fields (optimized fast path)
func (l *Logger) logSimple(level LogLevel, message string) {
	l.logSimpleExtremelyFast(level, message)
}

// logSimpleExtremelyFast provides ultra-fast simple message logging
func (l *Logger) logSimpleExtremelyFast(level LogLevel, message string) {
	l.logZeroBlazing(level, message)
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
