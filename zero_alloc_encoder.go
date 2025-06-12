package emit

import (
	"strconv"
	"time"
)

// ZeroAllocEncoder is a high-performance, zero-allocation JSON encoder
type ZeroAllocEncoder struct {
	buf        []byte
	scratch    [64]byte // Scratch space for number conversions
	fieldCount int
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
