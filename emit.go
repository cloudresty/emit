package emit

import "fmt"

// Public API - Main logging functions

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

// Enhanced API with Fields support

// InfoF logs an info message with Fields
func InfoF(message string, fields Fields) {
	if defaultLogger != nil {
		defaultLogger.log(INFO, message, fields.ToMap())
	}
}

// WarnF logs a warning message with Fields
func WarnF(message string, fields Fields) {
	if defaultLogger != nil {
		defaultLogger.log(WARN, message, fields.ToMap())
	}
}

// ErrorF logs an error message with Fields
func ErrorF(message string, fields Fields) {
	if defaultLogger != nil {
		defaultLogger.log(ERROR, message, fields.ToMap())
	}
}

// DebugF logs a debug message with Fields
func DebugF(message string, fields Fields) {
	if defaultLogger != nil {
		defaultLogger.log(DEBUG, message, fields.ToMap())
	}
}

// Variadic key-value pair API (requires even number of arguments)

// InfoKV logs an info message with key-value pairs
func InfoKV(message string, keyValuePairs ...any) {
	if defaultLogger != nil {
		fields := parseKeyValuePairs(keyValuePairs...)
		defaultLogger.log(INFO, message, fields)
	}
}

// WarnKV logs a warning message with key-value pairs
func WarnKV(message string, keyValuePairs ...any) {
	if defaultLogger != nil {
		fields := parseKeyValuePairs(keyValuePairs...)
		defaultLogger.log(WARN, message, fields)
	}
}

// ErrorKV logs an error message with key-value pairs
func ErrorKV(message string, keyValuePairs ...any) {
	if defaultLogger != nil {
		fields := parseKeyValuePairs(keyValuePairs...)
		defaultLogger.log(ERROR, message, fields)
	}
}

// DebugKV logs a debug message with key-value pairs
func DebugKV(message string, keyValuePairs ...any) {
	if defaultLogger != nil {
		fields := parseKeyValuePairs(keyValuePairs...)
		defaultLogger.log(DEBUG, message, fields)
	}
}

// parseKeyValuePairs converts variadic args to map[string]any
func parseKeyValuePairs(keyValuePairs ...any) map[string]any {
	if len(keyValuePairs) == 0 {
		return nil
	}

	// Ensure even number of arguments
	if len(keyValuePairs)%2 != 0 {
		// Add a placeholder for the missing value
		keyValuePairs = append(keyValuePairs, "<missing_value>")
	}

	fields := make(map[string]any, len(keyValuePairs)/2)
	for i := 0; i < len(keyValuePairs); i += 2 {
		key, ok := keyValuePairs[i].(string)
		if !ok {
			// Convert non-string keys to strings
			key = fmt.Sprintf("%v", keyValuePairs[i])
		}
		fields[key] = keyValuePairs[i+1]
	}
	return fields
}
