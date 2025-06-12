package emit

import (
	"runtime"
	"strconv"
	"sync"
	"time"
)

// ZeroAllocEncoder is a high-performance, zero-allocation JSON encoder
type ZeroAllocEncoder struct {
	buf        []byte
	scratch    [64]byte // Scratch space for number conversions
	fieldCount int
}

// Pool for zero-allocation encoders
var zeroAllocEncoderPool = sync.Pool{
	New: func() interface{} {
		return &ZeroAllocEncoder{
			buf: make([]byte, 0, 1024), // Start with 1KB capacity
		}
	},
}

// getZeroAllocEncoder gets an encoder from the pool
func getZeroAllocEncoder() *ZeroAllocEncoder {
	return zeroAllocEncoderPool.Get().(*ZeroAllocEncoder)
}

// putZeroAllocEncoder returns an encoder to the pool
func putZeroAllocEncoder(enc *ZeroAllocEncoder) {
	if cap(enc.buf) <= 4096 { // Don't pool very large buffers
		zeroAllocEncoderPool.Put(enc)
	}
}

// reset clears the encoder for reuse
func (e *ZeroAllocEncoder) reset() {
	e.buf = e.buf[:0]
	e.fieldCount = 0
}

// bytes returns the encoded bytes
func (e *ZeroAllocEncoder) bytes() []byte {
	return e.buf
}

// writeString appends a JSON-escaped string to the buffer
func (e *ZeroAllocEncoder) writeString(s string) {
	e.buf = append(e.buf, '"')

	// Fast path for strings without special characters
	needsEscape := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 32 || c == '"' || c == '\\' {
			needsEscape = true
			break
		}
	}

	if !needsEscape {
		e.buf = append(e.buf, s...)
	} else {
		// Slow path with proper escaping
		for _, r := range s {
			switch r {
			case '"':
				e.buf = append(e.buf, '\\', '"')
			case '\\':
				e.buf = append(e.buf, '\\', '\\')
			case '\n':
				e.buf = append(e.buf, '\\', 'n')
			case '\r':
				e.buf = append(e.buf, '\\', 'r')
			case '\t':
				e.buf = append(e.buf, '\\', 't')
			default:
				if r < 32 {
					e.buf = append(e.buf, '\\', 'u', '0', '0')
					e.buf = append(e.buf, "0123456789abcdef"[r>>4])
					e.buf = append(e.buf, "0123456789abcdef"[r&0xF])
				} else {
					// Use unsafe to convert rune to bytes efficiently
					if r <= 0x7F {
						e.buf = append(e.buf, byte(r))
					} else {
						// For non-ASCII, fall back to string conversion
						e.buf = append(e.buf, string(r)...)
					}
				}
			}
		}
	}

	e.buf = append(e.buf, '"')
}

// writeStringField writes a string field to JSON
func (e *ZeroAllocEncoder) writeStringField(key, value string) {
	if e.fieldCount > 0 {
		e.buf = append(e.buf, ',')
	}

	e.writeString(key)
	e.buf = append(e.buf, ':')
	e.writeString(value)

	e.fieldCount++
}

// writeIntField writes an integer field to JSON
func (e *ZeroAllocEncoder) writeIntField(key string, value int) {
	if e.fieldCount > 0 {
		e.buf = append(e.buf, ',')
	}

	e.writeString(key)
	e.buf = append(e.buf, ':')

	// Use scratch space for number conversion to avoid allocations
	num := strconv.AppendInt(e.scratch[:0], int64(value), 10)
	e.buf = append(e.buf, num...)

	e.fieldCount++
}

// writeInt64Field writes an int64 field to JSON
func (e *ZeroAllocEncoder) writeInt64Field(key string, value int64) {
	if e.fieldCount > 0 {
		e.buf = append(e.buf, ',')
	}

	e.writeString(key)
	e.buf = append(e.buf, ':')

	num := strconv.AppendInt(e.scratch[:0], value, 10)
	e.buf = append(e.buf, num...)

	e.fieldCount++
}

// writeFloat64Field writes a float64 field to JSON
func (e *ZeroAllocEncoder) writeFloat64Field(key string, value float64) {
	if e.fieldCount > 0 {
		e.buf = append(e.buf, ',')
	}

	e.writeString(key)
	e.buf = append(e.buf, ':')

	num := strconv.AppendFloat(e.scratch[:0], value, 'f', -1, 64)
	e.buf = append(e.buf, num...)

	e.fieldCount++
}

// writeBoolField writes a boolean field to JSON
func (e *ZeroAllocEncoder) writeBoolField(key string, value bool) {
	if e.fieldCount > 0 {
		e.buf = append(e.buf, ',')
	}

	e.writeString(key)
	e.buf = append(e.buf, ':')

	if value {
		e.buf = append(e.buf, "true"...)
	} else {
		e.buf = append(e.buf, "false"...)
	}

	e.fieldCount++
}

// writeTimeField writes a time field to JSON
func (e *ZeroAllocEncoder) writeTimeField(key string, value time.Time) {
	if e.fieldCount > 0 {
		e.buf = append(e.buf, ',')
	}

	e.writeString(key)
	e.buf = append(e.buf, ':')
	e.writeString(value.Format(time.RFC3339Nano))

	e.fieldCount++
}

// writeDurationField writes a duration field to JSON as nanoseconds
func (e *ZeroAllocEncoder) writeDurationField(key string, value time.Duration) {
	if e.fieldCount > 0 {
		e.buf = append(e.buf, ',')
	}

	e.writeString(key)
	e.buf = append(e.buf, ':')

	num := strconv.AppendInt(e.scratch[:0], int64(value), 10)
	e.buf = append(e.buf, num...)

	e.fieldCount++
}

// Zero-allocation JSON logging implementation
func (l *Logger) logZeroJSON(enc *ZeroAllocEncoder, level LogLevel, message string, fields ...ZField) {
	// Start JSON object
	enc.buf = append(enc.buf, '{')

	// Write timestamp
	enc.writeStringField("timestamp", time.Now().UTC().Format(time.RFC3339Nano))

	// Write level
	enc.writeStringField("level", level.String())

	// Write message
	enc.writeStringField("msg", message)

	// Write component if present
	if l.component != "" {
		enc.writeStringField("component", l.component)
	}

	// Write version if present
	if l.version != "" {
		enc.writeStringField("version", l.version)
	}

	// Write fields
	if len(fields) > 0 {
		if enc.fieldCount > 0 {
			enc.buf = append(enc.buf, ',')
		}

		enc.buf = append(enc.buf, `"fields":{`...)
		fieldEnc := &ZeroAllocEncoder{buf: make([]byte, 0, 256)}

		for _, field := range fields {
			field.WriteToEncoder(fieldEnc)
		}

		enc.buf = append(enc.buf, fieldEnc.buf...)
		enc.buf = append(enc.buf, '}')
	}

	// Write caller info if enabled
	if l.showCaller {
		if pc, file, line, ok := runtime.Caller(5); ok {
			enc.writeStringField("file", file)
			enc.writeIntField("line", line)
			if fn := runtime.FuncForPC(pc); fn != nil {
				enc.writeStringField("function", fn.Name())
			}
		}
	}

	// Close JSON object and add newline
	enc.buf = append(enc.buf, '}', '\n')
}

// Zero-allocation plain text logging implementation
func (l *Logger) logZeroPlain(enc *ZeroAllocEncoder, level LogLevel, message string, fields ...ZField) {
	severity := level.String()

	// Color codes
	var colorCode string
	switch severity {
	case "info":
		colorCode = "\033[32m"
	case "warn":
		colorCode = "\033[33m"
	case "error":
		colorCode = "\033[31m"
	case "debug":
		colorCode = "\033[34m"
	}

	resetCode := "\033[0m"
	if runtime.GOOS == "windows" {
		colorCode = ""
		resetCode = ""
	}

	// Build timestamp
	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")
	enc.buf = append(enc.buf, timestamp...)
	enc.buf = append(enc.buf, " | "...)

	// Build level with color
	enc.buf = append(enc.buf, colorCode...)
	enc.buf = append(enc.buf, severity...)

	// Pad to 7 characters for alignment
	for i := len(severity); i < 7; i++ {
		enc.buf = append(enc.buf, ' ')
	}
	enc.buf = append(enc.buf, resetCode...)
	enc.buf = append(enc.buf, " | "...)

	// Build component and version
	enc.buf = append(enc.buf, l.component...)
	enc.buf = append(enc.buf, ' ')
	enc.buf = append(enc.buf, l.version...)
	enc.buf = append(enc.buf, ": "...)

	// Build message
	enc.buf = append(enc.buf, message...)

	// Build fields if present
	if len(fields) > 0 {
		enc.buf = append(enc.buf, " ["...)

		for i, field := range fields {
			if i > 0 {
				enc.buf = append(enc.buf, ' ')
			}

			// For plain text, we need to extract key-value pairs
			// This is a simplified implementation - can be optimized further
			switch f := field.(type) {
			case StringZField:
				enc.buf = append(enc.buf, f.Key...)
				enc.buf = append(enc.buf, '=')
				if f.IsSensitive() {
					enc.buf = append(enc.buf, "***MASKED***"...)
				} else if f.IsPII() {
					enc.buf = append(enc.buf, "***PII***"...)
				} else {
					enc.buf = append(enc.buf, f.Value...)
				}
			case IntZField:
				enc.buf = append(enc.buf, f.Key...)
				enc.buf = append(enc.buf, '=')
				num := strconv.AppendInt(enc.scratch[:0], int64(f.Value), 10)
				enc.buf = append(enc.buf, num...)
			case BoolZField:
				enc.buf = append(enc.buf, f.Key...)
				enc.buf = append(enc.buf, '=')
				if f.Value {
					enc.buf = append(enc.buf, "true"...)
				} else {
					enc.buf = append(enc.buf, "false"...)
				}
			}
		}

		enc.buf = append(enc.buf, ']')
	}

	enc.buf = append(enc.buf, '\n')
}
