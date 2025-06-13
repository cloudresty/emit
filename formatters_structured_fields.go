package emit

import (
	"strconv"
	"unsafe"
)

// Structured fields - single allocation, perfect hot path
var (

	// Stack-based buffer for ultimate performance
	phase5bBuffer = [1024]byte{} // Increased buffer size for complex structured fields

	// Pre-computed JSON components for zero runtime cost
	structureJSONPrefix = []byte(`{"timestamp":"`)
	structureJSONLevel  = []byte(`","level":"info","msg":"`)
	structureJSONSuffix = []byte(`"}` + "\n")

	// Field separators and components
	structureFieldStart = []byte(`,"`)
	structureFieldMid   = []byte(`":"`)
	structureFieldEnd   = []byte(`"`)
	structureIntStart   = []byte(`,`)
)

// logStructuredFields - single allocation
func (l *Logger) logStructuredFields(level LogLevel, message string, fields ...ZField) {
	// Ultra-fast level check
	if level < l.level {
		return
	}

	// Use stack buffer for ZERO heap allocations
	buf := (*[1024]byte)(unsafe.Pointer(&phase5bBuffer[0]))
	pos := 0

	// Copy JSON prefix directly
	prefixLen := len(structureJSONPrefix)
	copy(buf[pos:], structureJSONPrefix)
	pos += prefixLen

	// Fast timestamp - single allocation shared across all logs
	ts := GetUltraFastTimestamp()
	copy(buf[pos:], ts)
	pos += len(ts)

	// Copy level and message sections
	levelLen := len(structureJSONLevel)
	copy(buf[pos:], structureJSONLevel)
	pos += levelLen

	// Copy message
	copy(buf[pos:], message)
	pos += len(message)

	// Process fields with zero allocations
	for _, field := range fields {
		switch f := field.(type) {
		case StringZField:
			// Fast string field processing
			copy(buf[pos:], structureFieldStart)
			pos += len(structureFieldStart)
			copy(buf[pos:], f.Key)
			pos += len(f.Key)
			copy(buf[pos:], structureFieldMid)
			pos += len(structureFieldMid)

			// Security check - ultra-fast using compile-time determination
			if !f.IsSensitive() && !f.IsPII() {
				copy(buf[pos:], f.Value)
				pos += len(f.Value)
			} else {
				mask := "***MASKED***"
				copy(buf[pos:], mask)
				pos += len(mask)
			}

			copy(buf[pos:], structureFieldEnd)
			pos += len(structureFieldEnd)

		case IntZField:
			// Ultra-fast integer field processing
			copy(buf[pos:], structureIntStart)
			pos += len(structureIntStart)
			copy(buf[pos:], f.Key)
			pos += len(f.Key)
			copy(buf[pos:], `":`)
			pos += 2

			// Fast integer conversion using separate buffer
			var numBuf [20]byte
			numStr := strconv.AppendInt(numBuf[:0], int64(f.Value), 10)
			copy(buf[pos:], numStr)
			pos += len(numStr)

		case Float64ZField:
			// Ultra-fast float field processing
			copy(buf[pos:], structureFieldStart)
			pos += len(structureFieldStart)
			copy(buf[pos:], f.Key)
			pos += len(f.Key)
			copy(buf[pos:], `":`)
			pos += 2

			// Fast float conversion using separate buffer
			var numBuf [32]byte
			numStr := strconv.AppendFloat(numBuf[:0], f.Value, 'g', -1, 64)
			copy(buf[pos:], numStr)
			pos += len(numStr)
		}
	}

	// Close JSON and add newline
	copy(buf[pos:], structureJSONSuffix)
	pos += len(structureJSONSuffix)

	// Single write operation - minimal syscall overhead
	l.writer.Write(buf[:pos])
}

// Route structured fields to implementation
func (l *Logger) logStructuredFieldsRoute(level LogLevel, message string, fields ...ZField) {
	// Route to implementation for maximum performance
	l.logStructuredFields(level, message, fields...)
}
