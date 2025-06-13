package emit

import "fmt"

// parseKeyValuePairs converts variadic args to map[string]any
// Used internally by the API for emit.Info.KeyValue() etc.
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

// Internal helper functions for the API
// These provide the actual logging implementation for the API namespace

func logWithFields(level LogLevel, message string, fields Fields) {
	if defaultLogger != nil {
		defaultLogger.log(level, message, fields.ToMap())
	}
}

func logWithKeyValues(level LogLevel, message string, keysAndValues ...interface{}) {
	if defaultLogger != nil {
		fields := parseKeyValuePairs(keysAndValues...)
		defaultLogger.log(level, message, fields)
	}
}

func logWithZeroAlloc(level LogLevel, message string, fields ...ZField) {
	if defaultLogger != nil {
		defaultLogger.logZero(level, message, fields...)
	}
}

func logWithPool(level LogLevel, message string, fn func(*PooledFields)) {
	if defaultLogger != nil {
		pf := NewPooledFields()
		fn(pf)
		defaultLogger.log(level, message, pf.ToMap())
		pf.Release()
	}
}

// Simple message logging functions with clear names
func InfoMsg(message string) {
	if defaultLogger != nil {
		defaultLogger.log(INFO, message, nil)
	}
}

func ErrorMsg(message string) {
	if defaultLogger != nil {
		defaultLogger.log(ERROR, message, nil)
	}
}

func WarnMsg(message string) {
	if defaultLogger != nil {
		defaultLogger.log(WARN, message, nil)
	}
}

func DebugMsg(message string) {
	if defaultLogger != nil {
		defaultLogger.log(DEBUG, message, nil)
	}
}

// InfoWithFields logs an info message with a map of fields
func InfoWithFields(message string, fields map[string]any) {
	if defaultLogger != nil {
		defaultLogger.log(INFO, message, fields)
	}
}

// Utility functions for custom integrations and special cases

// Log is a generic logging function that can be used for custom integrations
func Log(level, message string, optionalParams ...string) {
	logLevel := ParseLogLevel(level)

	// Handle optional parameters for component and version
	if len(optionalParams) >= 1 && defaultLogger.component == "" {
		SetComponent(optionalParams[0])
	}
	if len(optionalParams) >= 2 && defaultLogger.version == "" {
		SetVersion(optionalParams[1])
	}

	defaultLogger.log(logLevel, message, nil)
}

// JSON forces JSON output for a single log entry (for special cases)
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

// Plain forces plain output for a single log entry (for special cases)
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
