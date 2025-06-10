package emit

import (
	"time"
)

// Zero-allocation field types - inspired by Zap but with built-in security

// ZField represents a zero-allocation logging field
type ZField interface {
	WriteToEncoder(enc *ZeroAllocEncoder)
	IsSensitive() bool
	IsPII() bool
}

// StringZField represents a string field with zero allocations
type StringZField struct {
	Key   string
	Value string
}

func (f StringZField) WriteToEncoder(enc *ZeroAllocEncoder) {
	if f.IsSensitive() {
		enc.writeStringField(f.Key, "***MASKED***")
	} else if f.IsPII() {
		enc.writeStringField(f.Key, "***PII***")
	} else {
		enc.writeStringField(f.Key, f.Value)
	}
}

func (f StringZField) IsSensitive() bool {
	// Ultra-fast compile-time optimized security checks
	switch f.Key {
	case "password", "secret", "token", "api_key", "private_key":
		return true
	case "auth", "credential", "session", "jwt", "bearer":
		return true
	default:
		return false // Skip expensive pattern matching in hot path
	}
}

func (f StringZField) IsPII() bool {
	// Ultra-fast compile-time optimized PII checks
	switch f.Key {
	case "email", "phone", "name", "address", "ssn":
		return true
	case "user_email", "full_name", "credit_card", "passport":
		return true
	default:
		return false // Skip expensive pattern matching in hot path
	}
}

// IntZField represents an integer field with zero allocations
type IntZField struct {
	Key   string
	Value int
}

func (f IntZField) WriteToEncoder(enc *ZeroAllocEncoder) {
	enc.writeIntField(f.Key, f.Value)
}

func (f IntZField) IsSensitive() bool { return false }
func (f IntZField) IsPII() bool       { return false }

// Int64ZField represents an int64 field with zero allocations
type Int64ZField struct {
	Key   string
	Value int64
}

func (f Int64ZField) WriteToEncoder(enc *ZeroAllocEncoder) {
	enc.writeInt64Field(f.Key, f.Value)
}

func (f Int64ZField) IsSensitive() bool { return false }
func (f Int64ZField) IsPII() bool       { return false }

// Float64ZField represents a float64 field with zero allocations
type Float64ZField struct {
	Key   string
	Value float64
}

func (f Float64ZField) WriteToEncoder(enc *ZeroAllocEncoder) {
	enc.writeFloat64Field(f.Key, f.Value)
}

func (f Float64ZField) IsSensitive() bool { return false }
func (f Float64ZField) IsPII() bool       { return false }

// BoolZField represents a boolean field with zero allocations
type BoolZField struct {
	Key   string
	Value bool
}

func (f BoolZField) WriteToEncoder(enc *ZeroAllocEncoder) {
	enc.writeBoolField(f.Key, f.Value)
}

func (f BoolZField) IsSensitive() bool { return false }
func (f BoolZField) IsPII() bool       { return false }

// TimeZField represents a time field with zero allocations
type TimeZField struct {
	Key   string
	Value time.Time
}

func (f TimeZField) WriteToEncoder(enc *ZeroAllocEncoder) {
	enc.writeTimeField(f.Key, f.Value)
}

func (f TimeZField) IsSensitive() bool { return false }
func (f TimeZField) IsPII() bool       { return false }

// DurationZField represents a duration field with zero allocations
type DurationZField struct {
	Key   string
	Value time.Duration
}

func (f DurationZField) WriteToEncoder(enc *ZeroAllocEncoder) {
	enc.writeDurationField(f.Key, f.Value)
}

func (f DurationZField) IsSensitive() bool { return false }
func (f DurationZField) IsPII() bool       { return false }

// Zero-allocation field constructors

// ZString creates a zero-allocation string field
func ZString(key, value string) StringZField {
	return StringZField{Key: key, Value: value}
}

// ZInt creates a zero-allocation int field
func ZInt(key string, value int) IntZField {
	return IntZField{Key: key, Value: value}
}

// ZInt64 creates a zero-allocation int64 field
func ZInt64(key string, value int64) Int64ZField {
	return Int64ZField{Key: key, Value: value}
}

// ZFloat64 creates a zero-allocation float64 field
func ZFloat64(key string, value float64) Float64ZField {
	return Float64ZField{Key: key, Value: value}
}

// ZBool creates a zero-allocation bool field
func ZBool(key string, value bool) BoolZField {
	return BoolZField{Key: key, Value: value}
}

// ZTime creates a zero-allocation time field
func ZTime(key string, value time.Time) TimeZField {
	return TimeZField{Key: key, Value: value}
}

// ZDuration creates a zero-allocation duration field
func ZDuration(key string, value time.Duration) DurationZField {
	return DurationZField{Key: key, Value: value}
}

// Zero-allocation logging functions

// InfoZ logs an info message with zero-allocation fields
func InfoZ(message string, fields ...ZField) {
	if defaultLogger != nil && defaultLogger.level <= INFO {
		defaultLogger.logZero(INFO, message, fields...)
	}
}

// ErrorZ logs an error message with zero-allocation fields
func ErrorZ(message string, fields ...ZField) {
	if defaultLogger != nil && defaultLogger.level <= ERROR {
		defaultLogger.logZero(ERROR, message, fields...)
	}
}

// WarnZ logs a warning message with zero-allocation fields
func WarnZ(message string, fields ...ZField) {
	if defaultLogger != nil && defaultLogger.level <= WARN {
		defaultLogger.logZero(WARN, message, fields...)
	}
}

// DebugZ logs a debug message with zero-allocation fields
func DebugZ(message string, fields ...ZField) {
	if defaultLogger != nil && defaultLogger.level <= DEBUG {
		defaultLogger.logZero(DEBUG, message, fields...)
	}
}

// logZero performs zero-allocation logging
func (l *Logger) logZero(level LogLevel, message string, fields ...ZField) {
	if level < l.level {
		return
	}

	// Get encoder from pool
	enc := getZeroAllocEncoder()
	defer putZeroAllocEncoder(enc)

	// Reset encoder
	enc.reset()

	// Write log entry using zero-allocation encoder
	if l.format == PLAIN_FORMAT {
		l.logZeroPlain(enc, level, message, fields...)
	} else {
		l.logZeroJSON(enc, level, message, fields...)
	}

	// Write to output
	l.writer.Write(enc.bytes())
}
