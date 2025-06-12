package emit

import "time"

// Fields provides a fluent API for building log fields
type Fields map[string]any

// NewFields creates a new Fields instance
func NewFields() Fields {
	return make(Fields)
}

// Set adds or updates a field
func (f Fields) Set(key string, value any) Fields {
	f[key] = value
	return f
}

// Add is an alias for Set for more natural chaining
func (f Fields) Add(key string, value any) Fields {
	return f.Set(key, value)
}

// With is another alias for Set for more natural chaining
func (f Fields) With(key string, value any) Fields {
	return f.Set(key, value)
}

// String adds a string field
func (f Fields) String(key, value string) Fields {
	f[key] = value
	return f
}

// Int adds an integer field
func (f Fields) Int(key string, value int) Fields {
	f[key] = value
	return f
}

// Int64 adds an int64 field
func (f Fields) Int64(key string, value int64) Fields {
	f[key] = value
	return f
}

// Float64 adds a float64 field
func (f Fields) Float64(key string, value float64) Fields {
	f[key] = value
	return f
}

// Bool adds a boolean field
func (f Fields) Bool(key string, value bool) Fields {
	f[key] = value
	return f
}

// Time adds a time field (formats as RFC3339)
func (f Fields) Time(key string, value time.Time) Fields {
	f[key] = value.Format(time.RFC3339)
	return f
}

// Error adds an error field (converts to string)
func (f Fields) Error(key string, err error) Fields {
	if err != nil {
		f[key] = err.Error()
	} else {
		f[key] = nil
	}
	return f
}

// Any adds a field of any type
func (f Fields) Any(key string, value any) Fields {
	f[key] = value
	return f
}

// Merge combines multiple Fields objects
func (f Fields) Merge(other Fields) Fields {
	for k, v := range other {
		f[k] = v
	}
	return f
}

// Clone creates a copy of the Fields
func (f Fields) Clone() Fields {
	clone := make(Fields, len(f))
	for k, v := range f {
		clone[k] = v
	}
	return clone
}

// ToMap converts Fields to map[string]any for internal use
func (f Fields) ToMap() map[string]any {
	return map[string]any(f)
}

// Global helper functions for quick field creation

// F is a shorthand for creating Fields - the shortest possible API
// Deprecated: Use emit.Field() for clearer intent
// func F() Fields {
// 	return NewFields()
// }

// Field creates a single-field Fields object
func Field(key string, value any) Fields {
	return NewFields().Set(key, value)
}

// StringField creates a Fields object with a string field
func StringField(key, value string) Fields {
	return NewFields().String(key, value)
}

// IntField creates a Fields object with an integer field
func IntField(key string, value int) Fields {
	return NewFields().Int(key, value)
}

// ErrorField creates a Fields object with an error field
func ErrorField(key string, err error) Fields {
	return NewFields().Error(key, err)
}

// TimeField creates a Fields object with a time field
func TimeField(key string, value time.Time) Fields {
	return NewFields().Time(key, value)
}
