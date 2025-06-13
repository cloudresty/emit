package emit

// InfoLogger provides info-level logging methods with clear, simple names
type InfoLogger struct{}

// Field logs an info message with structured fields
func (InfoLogger) Field(msg string, fields Fields) {
	logWithFields(INFO, msg, fields)
}

// KeyValue logs an info message with key-value pairs
func (InfoLogger) KeyValue(msg string, keysAndValues ...interface{}) {
	logWithKeyValues(INFO, msg, keysAndValues...)
}

// StructuredFields logs an info message with structured fields
func (InfoLogger) StructuredFields(msg string, fields ...ZField) {
	if defaultLogger != nil {
		defaultLogger.logStructuredFields(INFO, msg, fields...)
	}
}

// Pool logs an info message using memory-pooled fields
func (InfoLogger) Pool(msg string, fn func(*PooledFields)) {
	logWithPool(INFO, msg, fn)
}

// Msg logs a simple info message
func (InfoLogger) Msg(msg string) {
	if defaultLogger != nil {
		defaultLogger.log(INFO, msg, nil)
	}
}

// ErrorLogger provides error-level logging methods with clear, simple names
type ErrorLogger struct{}

// Field logs an error message with structured fields
func (ErrorLogger) Field(msg string, fields Fields) {
	logWithFields(ERROR, msg, fields)
}

// KeyValue logs an error message with key-value pairs
func (ErrorLogger) KeyValue(msg string, keysAndValues ...interface{}) {
	logWithKeyValues(ERROR, msg, keysAndValues...)
}

// StructuredFields logs an error message with ultra-fast structured fields (Phase 5C)
func (ErrorLogger) StructuredFields(msg string, fields ...ZField) {
	if defaultLogger != nil {
		defaultLogger.logStructuredFields(ERROR, msg, fields...)
	}
}

// Pool logs an error message using memory-pooled fields
func (ErrorLogger) Pool(msg string, fn func(*PooledFields)) {
	logWithPool(ERROR, msg, fn)
}

// Msg logs a simple error message
func (ErrorLogger) Msg(msg string) {
	if defaultLogger != nil {
		defaultLogger.log(ERROR, msg, nil)
	}
}

// WarnLogger provides warn-level logging methods with clear, simple names
type WarnLogger struct{}

// Field logs a warn message with structured fields
func (WarnLogger) Field(msg string, fields Fields) {
	logWithFields(WARN, msg, fields)
}

// KeyValue logs a warn message with key-value pairs
func (WarnLogger) KeyValue(msg string, keysAndValues ...interface{}) {
	logWithKeyValues(WARN, msg, keysAndValues...)
}

// StructuredFields logs a warn message with ultra-fast structured fields
func (WarnLogger) StructuredFields(msg string, fields ...ZField) {
	if defaultLogger != nil {
		defaultLogger.logStructuredFields(WARN, msg, fields...)
	}
}

// Pool logs a warn message using memory-pooled fields
func (WarnLogger) Pool(msg string, fn func(*PooledFields)) {
	logWithPool(WARN, msg, fn)
}

// Msg logs a simple warn message
func (WarnLogger) Msg(msg string) {
	if defaultLogger != nil {
		defaultLogger.log(WARN, msg, nil)
	}
}

// DebugLogger provides debug-level logging methods with clear, simple names
type DebugLogger struct{}

// Field logs a debug message with structured fields
func (DebugLogger) Field(msg string, fields Fields) {
	logWithFields(DEBUG, msg, fields)
}

// KeyValue logs a debug message with key-value pairs
func (DebugLogger) KeyValue(msg string, keysAndValues ...interface{}) {
	logWithKeyValues(DEBUG, msg, keysAndValues...)
}

// StructuredFields logs a debug message with ultra-fast structured fields (Phase 5C)
func (DebugLogger) StructuredFields(msg string, fields ...ZField) {
	if defaultLogger != nil {
		defaultLogger.logStructuredFields(DEBUG, msg, fields...)
	}
}

// Pool logs a debug message using memory-pooled fields
func (DebugLogger) Pool(msg string, fn func(*PooledFields)) {
	logWithPool(DEBUG, msg, fn)
}

// Msg logs a simple debug message
func (DebugLogger) Msg(msg string) {
	if defaultLogger != nil {
		defaultLogger.log(DEBUG, msg, nil)
	}
}

var (
	// Info provides emit.Info.FieldWithMessage() and other methods
	Info = InfoLogger{}

	// Error provides emit.Error.FieldWithMessage() and other methods
	Error = ErrorLogger{}

	// Warn provides emit.Warn.FieldWithMessage() and other methods
	Warn = WarnLogger{}

	// Debug provides emit.Debug.FieldWithMessage() and other methods
	Debug = DebugLogger{}
)
