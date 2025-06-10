package emit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// Optimized formatters with buffer pooling and faster serialization

var (
	// Buffer pool for JSON formatting
	bufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 512)) // 512 bytes initial capacity
		},
	}

	// String builder pool for plain text formatting
	stringBuilderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
)

// getBuffer gets a buffer from the pool
func getBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// putBuffer returns a buffer to the pool
func putBuffer(buf *bytes.Buffer) {
	if buf.Cap() <= 2048 { // Don't pool very large buffers
		bufferPool.Put(buf)
	}
}

// getStringBuilder gets a string builder from the pool
func getStringBuilder() *strings.Builder {
	sb := stringBuilderPool.Get().(*strings.Builder)
	sb.Reset()
	return sb
}

// putStringBuilder returns a string builder to the pool
func putStringBuilder(sb *strings.Builder) {
	if sb.Cap() <= 2048 { // Don't pool very large builders
		stringBuilderPool.Put(sb)
	}
}

// Fast JSON serialization with manual field writing for common cases
func (l *Logger) logJSONFast(level LogLevel, message string, fields map[string]any) {
	buf := getBuffer()
	defer putBuffer(buf)

	// Start JSON object
	buf.WriteByte('{')

	// Write timestamp
	buf.WriteString(`"timestamp":"`)
	buf.WriteString(time.Now().UTC().Format(time.RFC3339Nano))
	buf.WriteByte('"')

	// Write level
	buf.WriteString(`,"level":"`)
	buf.WriteString(level.String())
	buf.WriteByte('"')

	// Write message
	buf.WriteString(`,"msg":"`)
	l.writeEscapedString(buf, message)
	buf.WriteByte('"')

	// Write component if present
	if l.component != "" {
		buf.WriteString(`,"component":"`)
		l.writeEscapedString(buf, l.component)
		buf.WriteByte('"')
	}

	// Write version if present
	if l.version != "" {
		buf.WriteString(`,"version":"`)
		l.writeEscapedString(buf, l.version)
		buf.WriteByte('"')
	}

	// Write fields if present
	if len(fields) > 0 {
		buf.WriteString(`,"fields":`)

		// Use optimized masking
		maskedFields := l.maskSensitiveFieldsFast(fields)

		// Fast path for small field sets (most common case)
		if len(maskedFields) <= 5 {
			l.writeFieldsFast(buf, maskedFields)
		} else {
			// Fallback to standard JSON marshaling for complex cases
			if fieldsJSON, err := json.Marshal(maskedFields); err == nil {
				buf.Write(fieldsJSON)
			} else {
				buf.WriteString(`{}`)
			}
		}
	}

	// Write caller info if enabled
	if l.showCaller {
		if pc, file, line, ok := runtime.Caller(4); ok {
			buf.WriteString(`,"file":"`)
			l.writeEscapedString(buf, file)
			buf.WriteString(`","line":`)
			buf.WriteString(strconv.Itoa(line))
			if fn := runtime.FuncForPC(pc); fn != nil {
				buf.WriteString(`,"function":"`)
				l.writeEscapedString(buf, fn.Name())
				buf.WriteByte('"')
			}
		}
	}

	// Close JSON object
	buf.WriteByte('}')
	buf.WriteByte('\n')

	// Write to output
	l.writer.Write(buf.Bytes())
}

// writeFieldsFast manually writes common field types for better performance
func (l *Logger) writeFieldsFast(buf *bytes.Buffer, fields map[string]any) {
	buf.WriteByte('{')

	first := true
	for key, value := range fields {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		// Write key
		buf.WriteByte('"')
		l.writeEscapedString(buf, key)
		buf.WriteString(`":`)

		// Write value based on type (fast path for common types)
		switch v := value.(type) {
		case string:
			buf.WriteByte('"')
			l.writeEscapedString(buf, v)
			buf.WriteByte('"')
		case int:
			buf.WriteString(strconv.Itoa(v))
		case int64:
			buf.WriteString(strconv.FormatInt(v, 10))
		case float64:
			buf.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
		case bool:
			if v {
				buf.WriteString("true")
			} else {
				buf.WriteString("false")
			}
		case nil:
			buf.WriteString("null")
		default:
			// Fallback to JSON marshaling for complex types
			if valueJSON, err := json.Marshal(v); err == nil {
				buf.Write(valueJSON)
			} else {
				buf.WriteString(`"<marshal_error>"`)
			}
		}
	}

	buf.WriteByte('}')
}

// writeEscapedString writes a JSON-escaped string to the buffer
func (l *Logger) writeEscapedString(buf *bytes.Buffer, s string) {
	for _, r := range s {
		switch r {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		case '\b':
			buf.WriteString(`\b`)
		case '\f':
			buf.WriteString(`\f`)
		default:
			if r < 32 {
				buf.WriteString(fmt.Sprintf(`\u%04x`, r))
			} else {
				buf.WriteRune(r)
			}
		}
	}
}

// Fast plain text formatting with string builder
func (l *Logger) logPlainFast(level LogLevel, message string, fields map[string]any) {
	sb := getStringBuilder()
	defer putStringBuilder(sb)

	severity := level.String()

	// Color codes
	var colorCode string
	switch severity {
	case "info":
		colorCode = "\033[32m" // Green
	case "warn":
		colorCode = "\033[33m" // Yellow
	case "error":
		colorCode = "\033[31m" // Red
	case "debug":
		colorCode = "\033[34m" // Blue
	default:
		colorCode = ""
	}

	resetCode := "\033[0m"
	if runtime.GOOS == "windows" {
		colorCode = ""
		resetCode = ""
	}

	// Build timestamp
	sb.WriteString(time.Now().UTC().Format("2006-01-02 15:04:05"))
	sb.WriteString(" | ")

	// Build level with color
	sb.WriteString(colorCode)
	sb.WriteString(severity)
	// Pad to 7 characters for alignment
	for i := len(severity); i < 7; i++ {
		sb.WriteByte(' ')
	}
	sb.WriteString(resetCode)
	sb.WriteString(" | ")

	// Build component and version
	sb.WriteString(l.component)
	sb.WriteByte(' ')
	sb.WriteString(l.version)
	sb.WriteString(": ")

	// Build message
	sb.WriteString(message)

	// Build fields if present
	if len(fields) > 0 {
		maskedFields := l.maskSensitiveFieldsFast(fields)
		sb.WriteString(" [")

		first := true
		for k, v := range maskedFields {
			if !first {
				sb.WriteByte(' ')
			}
			first = false

			sb.WriteString(k)
			sb.WriteByte('=')
			sb.WriteString(fmt.Sprintf("%v", v))
		}

		sb.WriteByte(']')
	}

	sb.WriteByte('\n')

	// Write to output using unsafe string conversion for better performance
	l.writer.Write(stringToBytes(sb.String()))
}

// stringToBytes converts string to []byte without allocation using unsafe
func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		Cap int
	}{s, len(s)}))
}
