package emit

import (
	"time"
)

// Zero-allocation logging implementation
// This is the real ZeroAlloc implementation that provides ultra-fast logging
// with minimal memory allocations for hot path scenarios.

// Pre-formatted JSON templates for zero-allocation logging
var (
	zeroAllocJSONPrefix = []byte(`{"timestamp":"`)
	zeroAllocLevelDebug = []byte(`","level":"debug","msg":"`)
	zeroAllocLevelInfo  = []byte(`","level":"info","msg":"`)
	zeroAllocLevelWarn  = []byte(`","level":"warn","msg":"`)
	zeroAllocLevelError = []byte(`","level":"error","msg":"`)
	zeroAllocJSONSuffix = []byte(`"}` + "\n")
)

// logZeroBlazing - Absolute fastest logging possible
func (l *Logger) logZeroBlazing(level LogLevel, message string, fields ...ZField) {
	if level < l.level {
		return // Critical: early exit
	}

	// Use even smaller stack buffer
	var stackBuf [280]byte
	var pos int

	if l.format == JSON_FORMAT {
		pos = l.buildJSONZeroAlloc(stackBuf[:], level, message, fields...)
	} else {
		pos = l.buildPlainZeroAlloc(stackBuf[:], level, message, fields...)
	}

	// Single write operation
	l.writer.Write(stackBuf[:pos])
}

// buildJSONZeroAlloc - Fastest possible JSON building using templates
func (l *Logger) buildJSONZeroAlloc(buf []byte, level LogLevel, message string, fields ...ZField) int {
	pos := 0

	// Template-based approach for maximum speed
	copy(buf[pos:], zeroAllocJSONPrefix)
	pos += len(zeroAllocJSONPrefix)

	// Use ultra-fast timestamp
	timestamp := GetUltraFastTimestamp()
	copy(buf[pos:], timestamp)
	pos += len(timestamp)

	// Pre-formatted level template
	var levelTemplate []byte
	switch level {
	case DEBUG:
		levelTemplate = zeroAllocLevelDebug
	case INFO:
		levelTemplate = zeroAllocLevelInfo
	case WARN:
		levelTemplate = zeroAllocLevelWarn
	case ERROR:
		levelTemplate = zeroAllocLevelError
	default:
		levelTemplate = zeroAllocLevelInfo
	}

	copy(buf[pos:], levelTemplate)
	pos += len(levelTemplate)

	// Message - direct copy (hot path assumes no escaping needed)
	copy(buf[pos:], message)
	pos += len(message)

	// Fields - ultra-minimal processing
	if len(fields) > 0 {
		for _, field := range fields {
			copy(buf[pos:], `","`)
			pos += 3
			pos = l.writeFieldZeroAlloc(buf, pos, field)
		}
	}

	// Close JSON using template
	copy(buf[pos:], zeroAllocJSONSuffix)
	pos += len(zeroAllocJSONSuffix)

	return pos
}

// writeFieldZeroAlloc - Absolute fastest field writing
func (l *Logger) writeFieldZeroAlloc(buf []byte, pos int, field ZField) int {
	switch f := field.(type) {
	case StringZField:
		// Ultra-fast string field: "key":"value"
		copy(buf[pos:], f.Key)
		pos += len(f.Key)
		copy(buf[pos:], `":"`)
		pos += 3
		copy(buf[pos:], f.Value)
		pos += len(f.Value)
		buf[pos] = '"'
		pos++

	case IntZField:
		// Ultra-fast int field: "key":123
		copy(buf[pos:], f.Key)
		pos += len(f.Key)
		copy(buf[pos:], `":`)
		pos += 2
		pos += writeIntZeroAlloc(buf[pos:], f.Value)

	case Float64ZField:
		// Ultra-fast float field: "key":25.4
		copy(buf[pos:], f.Key)
		pos += len(f.Key)
		copy(buf[pos:], `":`)
		pos += 2
		pos += writeFloat64ZeroAlloc(buf[pos:], f.Value)

	case BoolZField:
		// Ultra-fast bool field: "key":true
		copy(buf[pos:], f.Key)
		pos += len(f.Key)
		copy(buf[pos:], `":`)
		pos += 2
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

// writeIntZeroAlloc - Fastest integer conversion
func writeIntZeroAlloc(buf []byte, value int) int {
	// Optimized for common small values
	if value >= 0 && value <= 999 {
		return writeSmallIntZeroAlloc(buf, value)
	}

	// Handle negative and larger numbers
	if value < 0 {
		buf[0] = '-'
		return 1 + writeSmallIntZeroAlloc(buf[1:], -value)
	}

	// For larger positive numbers, use simple digit extraction
	return writeGenericIntZeroAlloc(buf, value)
}

func writeSmallIntZeroAlloc(buf []byte, value int) int {
	if value == 0 {
		buf[0] = '0'
		return 1
	}
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

func writeGenericIntZeroAlloc(buf []byte, value int) int {
	// Simple approach for larger numbers
	digits := 0
	temp := value
	for temp > 0 {
		temp /= 10
		digits++
	}

	for i := digits - 1; i >= 0; i-- {
		buf[i] = byte('0' + value%10)
		value /= 10
	}

	return digits
}

// writeFloat64ZeroAlloc - Fastest float conversion for common values
func writeFloat64ZeroAlloc(buf []byte, value float64) int {
	// Fast path for common float patterns
	if value == 0.0 {
		buf[0] = '0'
		return 1
	}

	// Handle integers that happen to be floats
	if value == float64(int(value)) {
		return writeIntZeroAlloc(buf, int(value))
	}

	// Simple 1-decimal approach for hot path
	intPart := int(value)
	fracPart := int((value-float64(intPart))*10 + 0.5) // Round to 1 decimal

	pos := writeIntZeroAlloc(buf, intPart)
	if fracPart > 0 && fracPart < 10 {
		buf[pos] = '.'
		buf[pos+1] = byte('0' + fracPart)
		return pos + 2
	}

	return pos
}

// buildPlainZeroAlloc - Fastest plain text building
func (l *Logger) buildPlainZeroAlloc(buf []byte, level LogLevel, message string, fields ...ZField) int {
	pos := 0

	// Simplified timestamp for plain text (HH:MM:SS)
	now := time.Now()
	hour, min, sec := now.Clock()

	pos += write2DigitsZeroAlloc(buf[pos:], hour)
	buf[pos] = ':'
	pos++
	pos += write2DigitsZeroAlloc(buf[pos:], min)
	buf[pos] = ':'
	pos++
	pos += write2DigitsZeroAlloc(buf[pos:], sec)

	copy(buf[pos:], " | ")
	pos += 3

	// Level with padding
	pos += writeLevelPaddedZeroAlloc(buf[pos:], level)

	copy(buf[pos:], " | ")
	pos += 3

	// Message
	copy(buf[pos:], message)
	pos += len(message)

	// Fields
	if len(fields) > 0 {
		copy(buf[pos:], " [")
		pos += 2

		for i, field := range fields {
			if i > 0 {
				buf[pos] = ' '
				pos++
			}
			pos = l.writeFieldPlainZeroAlloc(buf, pos, field)
		}

		buf[pos] = ']'
		pos++
	}

	buf[pos] = '\n'
	pos++

	return pos
}

func write2DigitsZeroAlloc(buf []byte, value int) int {
	buf[0] = byte('0' + value/10)
	buf[1] = byte('0' + value%10)
	return 2
}

func writeLevelPaddedZeroAlloc(buf []byte, level LogLevel) int {
	switch level {
	case DEBUG:
		copy(buf, "debug  ")
		return 7
	case INFO:
		copy(buf, "info   ")
		return 7
	case WARN:
		copy(buf, "warn   ")
		return 7
	case ERROR:
		copy(buf, "error  ")
		return 7
	default:
		copy(buf, "info   ")
		return 7
	}
}

func (l *Logger) writeFieldPlainZeroAlloc(buf []byte, pos int, field ZField) int {
	switch f := field.(type) {
	case StringZField:
		copy(buf[pos:], f.Key)
		pos += len(f.Key)
		buf[pos] = '='
		pos++
		copy(buf[pos:], f.Value)
		pos += len(f.Value)

	case IntZField:
		copy(buf[pos:], f.Key)
		pos += len(f.Key)
		buf[pos] = '='
		pos++
		pos += writeIntZeroAlloc(buf[pos:], f.Value)

	case Float64ZField:
		copy(buf[pos:], f.Key)
		pos += len(f.Key)
		buf[pos] = '='
		pos++
		pos += writeFloat64ZeroAlloc(buf[pos:], f.Value)

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
