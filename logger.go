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

	// Route to appropriate formatter based on format setting
	if l.format == PLAIN_FORMAT {
		l.logPlain(level, message, fields)
	} else {
		l.logJSON(level, message, fields)
	}
}
