package emit

import (
	"sync"
)

// Memory pools for reducing allocations

var (
	// Pool for map[string]any to reduce field map allocations
	fieldMapPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]any, 8) // Pre-allocate for common field count
		},
	}

	// Pool for Fields objects
	fieldsPool = sync.Pool{
		New: func() interface{} {
			return make(Fields, 8)
		},
	}

	// Pool for string slices used in plain text formatting
	stringSlicePool = sync.Pool{
		New: func() interface{} {
			return make([]string, 0, 10)
		},
	}
)

// getFieldMap gets a map from the pool
func getFieldMap() map[string]any {
	m := fieldMapPool.Get().(map[string]any)
	// Clear the map
	for k := range m {
		delete(m, k)
	}
	return m
}

// putFieldMap returns a map to the pool
func putFieldMap(m map[string]any) {
	if len(m) <= 32 { // Don't pool very large maps
		fieldMapPool.Put(m)
	}
}

// getFields gets a Fields object from the pool
func getFields() Fields {
	f := fieldsPool.Get().(Fields)
	// Clear the fields
	for k := range f {
		delete(f, k)
	}
	return f
}

// putFields returns a Fields object to the pool
func putFields(f Fields) {
	if len(f) <= 32 { // Don't pool very large fields
		fieldsPool.Put(f)
	}
}

// getStringSlice gets a string slice from the pool
func getStringSlice() []string {
	s := stringSlicePool.Get().([]string)
	return s[:0] // Reset length but keep capacity
}

// putStringSlice returns a string slice to the pool
func putStringSlice(s []string) {
	if cap(s) <= 50 { // Don't pool very large slices
		stringSlicePool.Put(s)
	}
}

// Optimized Fields implementation with pooling
type PooledFields struct {
	fields map[string]any
}

// NewPooledFields creates a new PooledFields using memory pool
func NewPooledFields() *PooledFields {
	return &PooledFields{
		fields: getFieldMap(),
	}
}

// Release returns the underlying map to the pool
func (pf *PooledFields) Release() {
	if pf.fields != nil {
		putFieldMap(pf.fields)
		pf.fields = nil
	}
}

// String adds a string field
func (pf *PooledFields) String(key, value string) *PooledFields {
	pf.fields[key] = value
	return pf
}

// Int adds an integer field
func (pf *PooledFields) Int(key string, value int) *PooledFields {
	pf.fields[key] = value
	return pf
}

// Bool adds a boolean field
func (pf *PooledFields) Bool(key string, value bool) *PooledFields {
	pf.fields[key] = value
	return pf
}

// Float64 adds a float64 field
func (pf *PooledFields) Float64(key string, value float64) *PooledFields {
	pf.fields[key] = value
	return pf
}

// Error adds an error field
func (pf *PooledFields) Error(key string, err error) *PooledFields {
	if err != nil {
		pf.fields[key] = err.Error()
	} else {
		pf.fields[key] = nil
	}
	return pf
}

// ToMap returns the underlying map
func (pf *PooledFields) ToMap() map[string]any {
	return pf.fields
}

// Pool-based field creation functions for high-performance scenarios

// PF creates a new PooledFields - ultra high performance API
func PF() *PooledFields {
	return NewPooledFields()
}

// WithPooledFields executes a function with pooled fields and automatically releases them
func WithPooledFields(fn func(*PooledFields)) {
	pf := NewPooledFields()
	defer pf.Release()
	fn(pf)
}

// High-performance logging functions that use pooled fields
func InfoFP(message string, fn func(*PooledFields)) {
	if defaultLogger != nil && defaultLogger.level <= INFO {
		pf := NewPooledFields()
		fn(pf)
		defaultLogger.log(INFO, message, pf.ToMap())
		pf.Release()
	}
}

func ErrorFP(message string, fn func(*PooledFields)) {
	if defaultLogger != nil && defaultLogger.level <= ERROR {
		pf := NewPooledFields()
		fn(pf)
		defaultLogger.log(ERROR, message, pf.ToMap())
		pf.Release()
	}
}

func WarnFP(message string, fn func(*PooledFields)) {
	if defaultLogger != nil && defaultLogger.level <= WARN {
		pf := NewPooledFields()
		fn(pf)
		defaultLogger.log(WARN, message, pf.ToMap())
		pf.Release()
	}
}

func DebugFP(message string, fn func(*PooledFields)) {
	if defaultLogger != nil && defaultLogger.level <= DEBUG {
		pf := NewPooledFields()
		fn(pf)
		defaultLogger.log(DEBUG, message, pf.ToMap())
		pf.Release()
	}
}
