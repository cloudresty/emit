package emit

import (
	"strconv"
)

// Structured fields - single allocation, perfect hot path
var (
	// Pre-computed JSON components for zero runtime cost
	structureJSONPrefix = []byte(`{"timestamp":"`)
)

// logStructuredFields - safe zero allocation with dynamic buffer support
func (l *Logger) logStructuredFields(level LogLevel, message string, fields ...ZField) {
	// Ultra-fast level check
	if level < l.level {
		return
	}

	// Start with stack buffer for ZERO heap allocations in common case
	var stackBuf [1024]byte
	buf := stackBuf[:]
	pos := 0

	// Estimate total size needed
	estimatedSize := l.estimateStructuredFieldsSize(level, message, fields...)
	var dynamicBuf []byte

	// If estimated size exceeds stack buffer, use dynamic allocation
	if estimatedSize > len(stackBuf) {
		dynamicBuf = make([]byte, estimatedSize)
		buf = dynamicBuf
	}

	// Build JSON with overflow protection
	pos = l.buildStructuredFieldsJSON(buf, level, message, fields...)

	// Check for overflow and retry with larger buffer if needed
	if pos >= len(buf) {
		// Buffer overflow - allocate larger buffer and retry
		retrySize := max(estimatedSize*2, 2048)
		dynamicBuf = make([]byte, retrySize)
		pos = l.buildStructuredFieldsJSON(dynamicBuf, level, message, fields...)

		// Final safety check
		if pos >= len(dynamicBuf) {
			// Fallback to safe JSON marshaling
			fieldMap := make(map[string]any)
			for _, field := range fields {
				switch f := field.(type) {
				case StringZField:
					fieldMap[f.Key] = f.Value
				case IntZField:
					fieldMap[f.Key] = f.Value
				case Float64ZField:
					fieldMap[f.Key] = f.Value
				}
			}
			l.logJSON(level, message, fieldMap)
			return
		}
		buf = dynamicBuf
	}

	// Single write operation - minimal syscall overhead
	l.writer.Write(buf[:pos])
}

// Route structured fields to implementation
func (l *Logger) logStructuredFieldsRoute(level LogLevel, message string, fields ...ZField) {
	// Route to implementation for maximum performance
	l.logStructuredFields(level, message, fields...)
}

// estimateStructuredFieldsSize calculates the approximate size needed for structured fields JSON
func (l *Logger) estimateStructuredFieldsSize(level LogLevel, message string, fields ...ZField) int {
	// Base JSON structure
	baseSize := 100

	// Timestamp and level
	timestampLevelSize := 50

	// Message length
	messageSize := len(message)

	// Component and version
	componentSize := 0
	if l.component != "" {
		componentSize = 15 + len(l.component)
	}
	versionSize := 0
	if l.version != "" {
		versionSize = 13 + len(l.version)
	}

	// Estimate field sizes
	fieldsSize := 0
	for _, field := range fields {
		switch f := field.(type) {
		case StringZField:
			fieldsSize += 10 + len(f.Key) + len(f.Value) // JSON overhead + key + value
		case IntZField:
			fieldsSize += 15 + len(f.Key) // JSON overhead + key + estimated number length
		case Float64ZField:
			fieldsSize += 25 + len(f.Key) // JSON overhead + key + estimated float length
		default:
			fieldsSize += 50 // Conservative estimate for unknown field types
		}
	}

	// Calculate total with 30% safety buffer
	totalSize := baseSize + timestampLevelSize + messageSize + componentSize + versionSize + fieldsSize
	return totalSize + (totalSize * 3 / 10) // Add 30% buffer
}

// buildStructuredFieldsJSON builds JSON with bounds checking
func (l *Logger) buildStructuredFieldsJSON(buf []byte, level LogLevel, message string, fields ...ZField) int {
	pos := 0

	// Check space and copy JSON prefix
	if pos+len(structureJSONPrefix) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], structureJSONPrefix)

	// Fast timestamp
	ts := GetUltraFastTimestamp()
	if pos+len(ts) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], ts)

	// Level and message sections
	levelStr := `","level":"` + level.StringFast() + `","msg":"`
	if pos+len(levelStr) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], levelStr)

	// Message
	if pos+len(message) >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], message)

	if pos+1 >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], `"`)

	// Process fields with bounds checking
	for _, field := range fields {
		switch f := field.(type) {
		case StringZField:
			fieldJSON := `,"` + f.Key + `":"`
			if pos+len(fieldJSON) >= len(buf) {
				return len(buf)
			}
			pos += copy(buf[pos:], fieldJSON)

			// Security check
			value := f.Value
			if f.IsSensitive() || f.IsPII() {
				value = "***MASKED***"
			}

			if pos+len(value) >= len(buf) {
				return len(buf)
			}
			pos += copy(buf[pos:], value)

			if pos+1 >= len(buf) {
				return len(buf)
			}
			pos += copy(buf[pos:], `"`)

		case IntZField:
			fieldJSON := `,"` + f.Key + `":`
			if pos+len(fieldJSON) >= len(buf) {
				return len(buf)
			}
			pos += copy(buf[pos:], fieldJSON)

			// Convert int to string
			var numBuf [20]byte
			numStr := strconv.AppendInt(numBuf[:0], int64(f.Value), 10)
			if pos+len(numStr) >= len(buf) {
				return len(buf)
			}
			pos += copy(buf[pos:], numStr)

		case Float64ZField:
			fieldJSON := `,"` + f.Key + `":`
			if pos+len(fieldJSON) >= len(buf) {
				return len(buf)
			}
			pos += copy(buf[pos:], fieldJSON)

			// Convert float to string
			var numBuf [32]byte
			numStr := strconv.AppendFloat(numBuf[:0], f.Value, 'g', -1, 64)
			if pos+len(numStr) >= len(buf) {
				return len(buf)
			}
			pos += copy(buf[pos:], numStr)
		}
	}

	// Add component if present
	if l.component != "" {
		componentJSON := `,"component":"` + l.component + `"`
		if pos+len(componentJSON) >= len(buf) {
			return len(buf)
		}
		pos += copy(buf[pos:], componentJSON)
	}

	// Add version if present
	if l.version != "" {
		versionJSON := `,"version":"` + l.version + `"`
		if pos+len(versionJSON) >= len(buf) {
			return len(buf)
		}
		pos += copy(buf[pos:], versionJSON)
	}

	// Close JSON and add newline
	if pos+2 >= len(buf) {
		return len(buf)
	}
	pos += copy(buf[pos:], "}\n")

	return pos
}
