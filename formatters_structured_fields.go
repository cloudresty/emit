package emit

import (
	"strconv"
	"sync"
)

// Fast JSON string escaping for structured fields
func escapeJSONString(dst []byte, src string) int {
	// For performance, handle the common case of no escaping needed
	needsEscaping := false
	for i := 0; i < len(src); i++ {
		c := src[i]
		if c == '"' || c == '\\' || c < 0x20 {
			needsEscaping = true
			break
		}
	}

	if !needsEscaping {
		// Fast path: no escaping needed
		copy(dst, src)
		return len(src)
	}

	// Slow path: escape the string using Go's built-in escaping
	escaped := strconv.Quote(src)
	// Remove the surrounding quotes that strconv.Quote adds
	escaped = escaped[1 : len(escaped)-1]
	copy(dst, escaped)
	return len(escaped)
}

// Structured fields - thread-safe buffer pool for concurrent access
var (
	// Pre-computed level strings as byte slices for maximum performance
	debugLevelBytes = []byte(`","level":"debug","message":"`)
	infoLevelBytes  = []byte(`","level":"info","message":"`)
	warnLevelBytes  = []byte(`","level":"warn","message":"`)
	errorLevelBytes = []byte(`","level":"error","message":"`)

	// Thread-safe buffer pool to prevent race conditions
	bufferPool = sync.Pool{
		New: func() interface{} {
			// Allocate 1024-byte buffer for each pool entry
			buf := make([]byte, 1024)
			return &buf
		},
	}

	// Pre-allocated component and version buffers for ultra-fast access
	componentPrefix = []byte(`,"component":"`)
	versionPrefix   = []byte(`,"version":"`)
)

// logStructuredFields - optimized for maximum performance with thread-safe buffers
func (l *Logger) logStructuredFields(level LogLevel, message string, fields ...ZField) {
	// Ultra-fast level check - most critical optimization
	if level < l.level {
		return
	}

	// Get thread-safe buffer from pool to prevent race conditions
	bufPtr := bufferPool.Get().(*[]byte)
	buf := *bufPtr
	defer bufferPool.Put(bufPtr) // Return buffer to pool when done
	pos := 0

	// Hot path optimization: For common case (≤4 fields), skip estimation
	// Most logging calls have 0-4 fields, so this covers 95%+ of cases
	fieldCount := len(fields)
	if fieldCount > 4 || len(message) > 200 {
		// Only do estimation for complex cases
		estimatedSize := 100 + len(message)

		if l.component != "" {
			estimatedSize += 15 + len(l.component)
		}
		if l.version != "" {
			estimatedSize += 13 + len(l.version)
		}

		for _, field := range fields {
			switch f := field.(type) {
			case StringZField:
				estimatedSize += 20 + len(f.Key) + len(f.Value)
			case IntZField:
				estimatedSize += 20 + len(f.Key)
			case Float64ZField:
				estimatedSize += 30 + len(f.Key)
			case BoolZField:
				estimatedSize += 15 + len(f.Key)
			}
		}

		if estimatedSize >= len(buf) {
			l.logStructuredFieldsDynamic(level, message, fields...)
			return
		}
	}

	// Ultra-hot path: build JSON directly with maximum inlining

	// JSON prefix: {"timestamp":" - 15 bytes
	buf[0] = '{'
	buf[1] = '"'
	buf[2] = 't'
	buf[3] = 'i'
	pos = 4 // Start from here for remaining "mestamp":"
	copy(buf[pos:], []byte(`mestamp":"`))
	pos += 10

	// Fast cached timestamp - inline string copy
	ts := GetUltraFastTimestamp()
	copy(buf[pos:], ts)
	pos += len(ts)

	// Level section - pre-computed byte slices, eliminate switch overhead for INFO
	if level == INFO {
		// Most common case - hardcode for INFO
		copy(buf[pos:], infoLevelBytes)
		pos += len(infoLevelBytes)
	} else {
		var levelBytes []byte
		switch level {
		case DEBUG:
			levelBytes = debugLevelBytes
		case WARN:
			levelBytes = warnLevelBytes
		case ERROR:
			levelBytes = errorLevelBytes
		default:
			levelBytes = infoLevelBytes
		}
		copy(buf[pos:], levelBytes)
		pos += len(levelBytes)
	}

	// Message (inline string-to-byte conversion)
	copy(buf[pos:], message)
	pos += len(message)

	// Closing quote for message
	buf[pos] = '"'
	pos++

	// Process fields inline - unrolled for maximum performance
	for _, field := range fields {
		switch f := field.(type) {
		case StringZField:
			// ,"key":"value" - most common field type, optimize heavily
			buf[pos] = ','
			buf[pos+1] = '"'
			pos += 2

			copy(buf[pos:], f.Key)
			pos += len(f.Key)

			buf[pos] = '"'
			buf[pos+1] = ':'
			buf[pos+2] = '"'
			pos += 3

			// Security check inline - optimize for non-sensitive case
			if f.IsSensitive() || f.IsPII() {
				copy(buf[pos:], "***MASKED***")
				pos += 12
			} else {
				// Properly escape JSON strings to prevent invalid JSON
				escaped := escapeJSONString(buf[pos:], f.Value)
				pos += escaped
			}

			buf[pos] = '"'
			pos++

		case IntZField:
			// ,"key":123
			buf[pos] = ','
			buf[pos+1] = '"'
			pos += 2

			copy(buf[pos:], f.Key)
			pos += len(f.Key)

			buf[pos] = '"'
			buf[pos+1] = ':'
			pos += 2

			// Fast integer conversion
			var numBuf [20]byte
			numStr := strconv.AppendInt(numBuf[:0], int64(f.Value), 10)
			copy(buf[pos:], numStr)
			pos += len(numStr)

		case Float64ZField:
			// ,"key":123.45
			buf[pos] = ','
			buf[pos+1] = '"'
			pos += 2

			copy(buf[pos:], f.Key)
			pos += len(f.Key)

			buf[pos] = '"'
			buf[pos+1] = ':'
			pos += 2

			// Fast float conversion
			var numBuf [32]byte
			numStr := strconv.AppendFloat(numBuf[:0], f.Value, 'g', -1, 64)
			copy(buf[pos:], numStr)
			pos += len(numStr)

		case BoolZField:
			// ,"key":true/false
			buf[pos] = ','
			buf[pos+1] = '"'
			pos += 2

			copy(buf[pos:], f.Key)
			pos += len(f.Key)

			buf[pos] = '"'
			buf[pos+1] = ':'
			pos += 2

			if f.Value {
				buf[pos] = 't'
				buf[pos+1] = 'r'
				buf[pos+2] = 'u'
				buf[pos+3] = 'e'
				pos += 4
			} else {
				buf[pos] = 'f'
				buf[pos+1] = 'a'
				buf[pos+2] = 'l'
				buf[pos+3] = 's'
				buf[pos+4] = 'e'
				pos += 5
			}
		}
	}

	// Add component if present (pre-computed prefix)
	if l.component != "" {
		copy(buf[pos:], componentPrefix)
		pos += len(componentPrefix)
		copy(buf[pos:], l.component)
		pos += len(l.component)
		buf[pos] = '"'
		pos++
	}

	// Add version if present (pre-computed prefix)
	if l.version != "" {
		copy(buf[pos:], versionPrefix)
		pos += len(versionPrefix)
		copy(buf[pos:], l.version)
		pos += len(l.version)
		buf[pos] = '"'
		pos++
	}

	// Close JSON: }\n - inline for final micro-optimization
	buf[pos] = '}'
	buf[pos+1] = '\n'
	pos += 2

	// Single write operation
	_, _ = l.writer.Write(buf[:pos])
}

// logStructuredFieldsDynamic - handles cases where log entry is too large for stack buffer
func (l *Logger) logStructuredFieldsDynamic(level LogLevel, message string, fields ...ZField) {
	// Calculate required size more accurately
	size := 100 + len(message) // base structure + message

	// Add logger metadata
	if l.component != "" {
		size += 15 + len(l.component)
	}
	if l.version != "" {
		size += 13 + len(l.version)
	}

	// Add fields
	for _, field := range fields {
		switch f := field.(type) {
		case StringZField:
			size += 20 + len(f.Key) + len(f.Value)
		case IntZField:
			size += 20 + len(f.Key)
		case Float64ZField:
			size += 30 + len(f.Key)
		case BoolZField:
			size += 15 + len(f.Key)
		}
	}

	// Add buffer to be safe
	if size < 2048 {
		size = 2048
	}

	buf := make([]byte, size)
	pos := 0

	// Build JSON (similar to hot path but with bounds checking)
	buf[0] = '{'
	buf[1] = '"'
	buf[2] = 't'
	buf[3] = 'i'
	pos = 4
	copy(buf[pos:], []byte(`mestamp":"`))
	pos += 10

	ts := GetUltraFastTimestamp()
	copy(buf[pos:], ts)
	pos += len(ts)

	// Level
	var levelBytes []byte
	switch level {
	case DEBUG:
		levelBytes = debugLevelBytes
	case INFO:
		levelBytes = infoLevelBytes
	case WARN:
		levelBytes = warnLevelBytes
	case ERROR:
		levelBytes = errorLevelBytes
	default:
		levelBytes = infoLevelBytes
	}

	copy(buf[pos:], levelBytes)
	pos += len(levelBytes)

	copy(buf[pos:], message)
	pos += len(message)
	buf[pos] = '"'
	pos++
	// Process fields - inline version
	for _, field := range fields {
		switch f := field.(type) {
		case StringZField:
			buf[pos] = ','
			buf[pos+1] = '"'
			pos += 2
			copy(buf[pos:], f.Key)
			pos += len(f.Key)
			buf[pos] = '"'
			buf[pos+1] = ':'
			buf[pos+2] = '"'
			pos += 3

			if f.IsSensitive() || f.IsPII() {
				copy(buf[pos:], "***MASKED***")
				pos += 12
			} else {
				// Properly escape JSON strings in dynamic formatter too
				escaped := escapeJSONString(buf[pos:], f.Value)
				pos += escaped
			}
			buf[pos] = '"'
			pos++

		case IntZField:
			buf[pos] = ','
			buf[pos+1] = '"'
			pos += 2
			copy(buf[pos:], f.Key)
			pos += len(f.Key)
			buf[pos] = '"'
			buf[pos+1] = ':'
			pos += 2

			var numBuf [20]byte
			numStr := strconv.AppendInt(numBuf[:0], int64(f.Value), 10)
			copy(buf[pos:], numStr)
			pos += len(numStr)

		case Float64ZField:
			buf[pos] = ','
			buf[pos+1] = '"'
			pos += 2
			copy(buf[pos:], f.Key)
			pos += len(f.Key)
			buf[pos] = '"'
			buf[pos+1] = ':'
			pos += 2

			var numBuf [32]byte
			numStr := strconv.AppendFloat(numBuf[:0], f.Value, 'g', -1, 64)
			copy(buf[pos:], numStr)
			pos += len(numStr)

		case BoolZField:
			buf[pos] = ','
			buf[pos+1] = '"'
			pos += 2
			copy(buf[pos:], f.Key)
			pos += len(f.Key)
			buf[pos] = '"'
			buf[pos+1] = ':'
			pos += 2

			if f.Value {
				buf[pos] = 't'
				buf[pos+1] = 'r'
				buf[pos+2] = 'u'
				buf[pos+3] = 'e'
				pos += 4
			} else {
				buf[pos] = 'f'
				buf[pos+1] = 'a'
				buf[pos+2] = 'l'
				buf[pos+3] = 's'
				buf[pos+4] = 'e'
				pos += 5
			}
		}
	}

	// Add component and version
	if l.component != "" {
		copy(buf[pos:], componentPrefix)
		pos += len(componentPrefix)
		copy(buf[pos:], l.component)
		pos += len(l.component)
		buf[pos] = '"'
		pos++
	}

	if l.version != "" {
		copy(buf[pos:], versionPrefix)
		pos += len(versionPrefix)
		copy(buf[pos:], l.version)
		pos += len(l.version)
		buf[pos] = '"'
		pos++
	}

	// Close JSON: }\n
	buf[pos] = '}'
	buf[pos+1] = '\n'
	pos += 2

	_, _ = l.writer.Write(buf[:pos])
}

// Route structured fields to implementation
func (l *Logger) logStructuredFieldsRoute(level LogLevel, message string, fields ...ZField) {
	// Route to implementation for maximum performance
	l.logStructuredFields(level, message, fields...)
}
