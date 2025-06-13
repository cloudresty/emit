package emit

import (
	"sync"
	"time"
)

// Memory pools for reducing allocations

var (
	// Pool for map[string]any to reduce field map allocations
	fieldMapPool = sync.Pool{
		New: func() any {
			return make(map[string]any, 8) // Pre-allocate for common field count
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

// Int64 adds an int64 field
func (pf *PooledFields) Int64(key string, value int64) *PooledFields {
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

// Time adds a time field (formats as RFC3339)
func (pf *PooledFields) Time(key string, value time.Time) *PooledFields {
	pf.fields[key] = value.Format(time.RFC3339)
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

// WithPooledFields executes a function with pooled fields and automatically releases them
func WithPooledFields(fn func(*PooledFields)) {
	pf := NewPooledFields()
	defer pf.Release()
	fn(pf)
}
