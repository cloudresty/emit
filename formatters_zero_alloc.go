package emit

import (
	"strconv"
)

// Zero-allocation encoder
type HighFrequencyEncoder struct {
	stackBuf [512]byte // Stack-allocated buffer for hot path
	pos      int       // Current position in buffer
}

// Field type detection and processing
const (
	FIELD_STRING = iota
	FIELD_INT
	FIELD_FLOAT64
	FIELD_BOOL
)

// Pre-computed hot path field metadata for O(1) lookups
type FieldMeta struct {
	fieldType int
	keyBytes  []byte
	sensitive bool
	pii       bool
}

// Global hot path field cache (populated at startup)
var hotPathFields = map[string]FieldMeta{
	"method":      {FIELD_STRING, []byte("method"), false, false},
	"path":        {FIELD_STRING, []byte("path"), false, false},
	"status":      {FIELD_INT, []byte("status"), false, false},
	"duration_ms": {FIELD_FLOAT64, []byte("duration_ms"), false, false},
	"user_id":     {FIELD_STRING, []byte("user_id"), false, false},
	"request_id":  {FIELD_STRING, []byte("request_id"), false, false},
	"service":     {FIELD_STRING, []byte("service"), false, false},
	"endpoint":    {FIELD_STRING, []byte("endpoint"), false, false},
	"cached":      {FIELD_BOOL, []byte("cached"), false, false},
}

// logZeroHighFrequency - Ultra-optimized zero-allocation logging for hot paths
// This is the fastest possible implementation targeting <45ns/op
func (l *Logger) logZeroHighFrequency(level LogLevel, message string, fields ...ZField) {
	if level < l.level {
		return // Early exit - most critical optimization
	}

	// Stack-allocated encoder - zero heap allocation
	var encoder HighFrequencyEncoder

	// Fast path: JSON format optimization (most common case)
	if l.format == JSON_FORMAT {
		l.buildJSONHighFrequency(&encoder, level, message, fields...)
	} else {
		l.buildPlainHighFrequency(&encoder, level, message, fields...)
	}

	// Single write operation
	l.writer.Write(encoder.stackBuf[:encoder.pos])
}

// buildJSONHighFrequency - JSON building using stack buffer
func (l *Logger) buildJSONHighFrequency(enc *HighFrequencyEncoder, level LogLevel, message string, fields ...ZField) {
	buf := enc.stackBuf[:]
	pos := 0

	// Start JSON object
	buf[pos] = '{'
	pos++

	// Write timestamp - use cached timestamp for maximum speed
	const timestampPrefix = `"timestamp":"`
	copy(buf[pos:], timestampPrefix)
	pos += len(timestampPrefix)

	// Use cached timestamp
	timestamp := GetUltraFastTimestamp()
	copy(buf[pos:], timestamp)
	pos += len(timestamp)

	buf[pos] = '"'
	pos++

	// Write level - use pre-computed level strings
	const levelPrefix = `,"level":"`
	copy(buf[pos:], levelPrefix)
	pos += len(levelPrefix)

	// Direct level string copy - avoid function calls
	levelStr := getLevelStringDirect(level)
	copy(buf[pos:], levelStr)
	pos += len(levelStr)

	buf[pos] = '"'
	pos++

	// Write message - direct copy (assume no escaping needed in hot path)
	const msgPrefix = `,"msg":"`
	copy(buf[pos:], msgPrefix)
	pos += len(msgPrefix)

	copy(buf[pos:], message)
	pos += len(message)

	buf[pos] = '"'
	pos++

	// Write component and version if present (most loggers have these)
	if l.component != "" {
		const componentPrefix = `,"component":"`
		copy(buf[pos:], componentPrefix)
		pos += len(componentPrefix)
		copy(buf[pos:], l.component)
		pos += len(l.component)
		buf[pos] = '"'
		pos++
	}

	if l.version != "" {
		const versionPrefix = `,"version":"`
		copy(buf[pos:], versionPrefix)
		pos += len(versionPrefix)
		copy(buf[pos:], l.version)
		pos += len(l.version)
		buf[pos] = '"'
		pos++
	}

	// Write fields - optimized hot path
	if len(fields) > 0 {
		for _, field := range fields {
			buf[pos] = ','
			pos++
			pos = l.writeFieldHighFrequency(buf, pos, field)
		}
	}

	// Close JSON object
	buf[pos] = '}'
	pos++
	buf[pos] = '\n'
	pos++

	enc.pos = pos
}

// writeFieldHighFrequency - Optimized field writing using hot path optimization
func (l *Logger) writeFieldHighFrequency(buf []byte, pos int, field ZField) int {
	switch f := field.(type) {
	case StringZField:
		// Hot path field lookup
		if meta, found := hotPathFields[f.Key]; found {
			buf[pos] = '"'
			pos++
			copy(buf[pos:], meta.keyBytes)
			pos += len(meta.keyBytes)
			buf[pos] = '"'
			pos++
			buf[pos] = ':'
			pos++
			buf[pos] = '"'
			pos++

			// Security check only if needed
			if meta.sensitive && l.sensitiveMode == MASK_SENSITIVE {
				copy(buf[pos:], l.maskString)
				pos += len(l.maskString)
			} else if meta.pii && l.piiMode == MASK_PII {
				copy(buf[pos:], l.piiMaskString)
				pos += len(l.piiMaskString)
			} else {
				copy(buf[pos:], f.Value)
				pos += len(f.Value)
			}

			buf[pos] = '"'
			pos++
		} else {
			// Fallback for unknown fields
			pos = l.writeStringFieldSlow(buf, pos, f.Key, f.Value, f.IsSensitive(), f.IsPII())
		}

	case IntZField:
		if meta, found := hotPathFields[f.Key]; found {
			buf[pos] = '"'
			pos++
			copy(buf[pos:], meta.keyBytes)
			pos += len(meta.keyBytes)
			buf[pos] = '"'
			pos++
			buf[pos] = ':'
			pos++

			// Direct integer conversion to buffer
			pos += writeIntDirect(buf[pos:], f.Value)
		} else {
			pos = l.writeIntFieldSlow(buf, pos, f.Key, f.Value)
		}

	case Float64ZField:
		if meta, found := hotPathFields[f.Key]; found {
			buf[pos] = '"'
			pos++
			copy(buf[pos:], meta.keyBytes)
			pos += len(meta.keyBytes)
			buf[pos] = '"'
			pos++
			buf[pos] = ':'
			pos++

			// Direct float conversion to buffer
			pos += writeFloat64Direct(buf[pos:], f.Value)
		} else {
			pos = l.writeFloat64FieldSlow(buf, pos, f.Key, f.Value)
		}

	case BoolZField:
		if meta, found := hotPathFields[f.Key]; found {
			buf[pos] = '"'
			pos++
			copy(buf[pos:], meta.keyBytes)
			pos += len(meta.keyBytes)
			buf[pos] = '"'
			pos++
			buf[pos] = ':'
			pos++

			if f.Value {
				copy(buf[pos:], "true")
				pos += 4
			} else {
				copy(buf[pos:], "false")
				pos += 5
			}
		} else {
			pos = l.writeBoolFieldSlow(buf, pos, f.Key, f.Value)
		}
	}

	return pos
}

// Ultra-fast direct number conversion functions
func writeIntDirect(buf []byte, value int) int {
	if value == 0 {
		buf[0] = '0'
		return 1
	}

	// Fast path for common small integers
	if value > 0 && value < 1000 {
		return writeSmallIntDirect(buf, value)
	}

	// Use strconv for larger numbers
	str := strconv.Itoa(value)
	copy(buf, str)
	return len(str)
}

func writeSmallIntDirect(buf []byte, value int) int {
	if value < 10 {
		buf[0] = byte('0' + value)
		return 1
	}
	if value < 100 {
		buf[0] = byte('0' + value/10)
		buf[1] = byte('0' + value%10)
		return 2
	}
	// value < 1000
	buf[0] = byte('0' + value/100)
	buf[1] = byte('0' + (value/10)%10)
	buf[2] = byte('0' + value%10)
	return 3
}

func writeFloat64Direct(buf []byte, value float64) int {
	// For hot path, use optimized float conversion
	str := strconv.FormatFloat(value, 'f', -1, 64)
	copy(buf, str)
	return len(str)
}

// getLevelStringDirect - Ultra-fast level string lookup
func getLevelStringDirect(level LogLevel) string {
	switch level {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	default:
		return "info"
	}
}

// Fallback methods for unknown fields (slower but safe)
func (l *Logger) writeStringFieldSlow(buf []byte, pos int, key, value string, sensitive, pii bool) int {
	buf[pos] = '"'
	pos++
	copy(buf[pos:], key)
	pos += len(key)
	buf[pos] = '"'
	pos++
	buf[pos] = ':'
	pos++
	buf[pos] = '"'
	pos++

	if sensitive && l.sensitiveMode == MASK_SENSITIVE {
		copy(buf[pos:], l.maskString)
		pos += len(l.maskString)
	} else if pii && l.piiMode == MASK_PII {
		copy(buf[pos:], l.piiMaskString)
		pos += len(l.piiMaskString)
	} else {
		copy(buf[pos:], value)
		pos += len(value)
	}

	buf[pos] = '"'
	pos++
	return pos
}

func (l *Logger) writeIntFieldSlow(buf []byte, pos int, key string, value int) int {
	buf[pos] = '"'
	pos++
	copy(buf[pos:], key)
	pos += len(key)
	buf[pos] = '"'
	pos++
	buf[pos] = ':'
	pos++
	pos += writeIntDirect(buf[pos:], value)
	return pos
}

func (l *Logger) writeFloat64FieldSlow(buf []byte, pos int, key string, value float64) int {
	buf[pos] = '"'
	pos++
	copy(buf[pos:], key)
	pos += len(key)
	buf[pos] = '"'
	pos++
	buf[pos] = ':'
	pos++
	pos += writeFloat64Direct(buf[pos:], value)
	return pos
}

func (l *Logger) writeBoolFieldSlow(buf []byte, pos int, key string, value bool) int {
	buf[pos] = '"'
	pos++
	copy(buf[pos:], key)
	pos += len(key)
	buf[pos] = '"'
	pos++
	buf[pos] = ':'
	pos++

	if value {
		copy(buf[pos:], "true")
		pos += 4
	} else {
		copy(buf[pos:], "false")
		pos += 5
	}
	return pos
}

// buildPlainHighFrequency - Ultra-fast plain text building
func (l *Logger) buildPlainHighFrequency(enc *HighFrequencyEncoder, level LogLevel, message string, fields ...ZField) {
	buf := enc.stackBuf[:]
	pos := 0

	// Build timestamp (truncated for speed)
	timestamp := GetUltraFastTimestamp()
	if len(timestamp) >= 19 {
		copy(buf[pos:], timestamp[:19])
		pos += 19
	} else {
		copy(buf[pos:], timestamp)
		pos += len(timestamp)
	}

	copy(buf[pos:], " | ")
	pos += 3

	// Build level (pre-padded)
	levelStr := getLevelStringPadded(level)
	copy(buf[pos:], levelStr)
	pos += 7 // All level strings are padded to 7 chars

	copy(buf[pos:], " | ")
	pos += 3

	// Build component and version
	if l.component != "" {
		copy(buf[pos:], l.component)
		pos += len(l.component)
		buf[pos] = ' '
		pos++
	}

	if l.version != "" {
		copy(buf[pos:], l.version)
		pos += len(l.version)
		copy(buf[pos:], ": ")
		pos += 2
	}

	// Build message
	copy(buf[pos:], message)
	pos += len(message)

	// Build fields
	if len(fields) > 0 {
		copy(buf[pos:], " [")
		pos += 2

		for i, field := range fields {
			if i > 0 {
				buf[pos] = ' '
				pos++
			}
			pos = l.writeFieldPlainHighFrequency(buf, pos, field)
		}

		buf[pos] = ']'
		pos++
	}

	buf[pos] = '\n'
	pos++

	enc.pos = pos
}

// getLevelStringPadded - Pre-padded level strings for alignment
func getLevelStringPadded(level LogLevel) string {
	switch level {
	case DEBUG:
		return "debug  " // Padded to 7 chars
	case INFO:
		return "info   " // Padded to 7 chars
	case WARN:
		return "warn   " // Padded to 7 chars
	case ERROR:
		return "error  " // Padded to 7 chars
	default:
		return "info   " // Padded to 7 chars
	}
}

// writeFieldPlainHighFrequency - Ultra-fast plain text field writing
func (l *Logger) writeFieldPlainHighFrequency(buf []byte, pos int, field ZField) int {
	switch f := field.(type) {
	case StringZField:
		copy(buf[pos:], f.Key)
		pos += len(f.Key)
		buf[pos] = '='
		pos++

		if f.IsSensitive() && l.sensitiveMode == MASK_SENSITIVE {
			copy(buf[pos:], l.maskString)
			pos += len(l.maskString)
		} else if f.IsPII() && l.piiMode == MASK_PII {
			copy(buf[pos:], l.piiMaskString)
			pos += len(l.piiMaskString)
		} else {
			copy(buf[pos:], f.Value)
			pos += len(f.Value)
		}

	case IntZField:
		copy(buf[pos:], f.Key)
		pos += len(f.Key)
		buf[pos] = '='
		pos++
		pos += writeIntDirect(buf[pos:], f.Value)

	case Float64ZField:
		copy(buf[pos:], f.Key)
		pos += len(f.Key)
		buf[pos] = '='
		pos++
		pos += writeFloat64Direct(buf[pos:], f.Value)

	case BoolZField:
		copy(buf[pos:], f.Key)
		pos += len(f.Key)
		buf[pos] = '='
		pos++
		if f.Value {
			copy(buf[pos:], "true")
			pos += 4
		} else {
			copy(buf[pos:], "false")
			pos += 5
		}
	}

	return pos
}
