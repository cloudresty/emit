package emit

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

// TestFieldsBuilder tests the Fields builder methods
func TestFieldsBuilder(t *testing.T) {
	f := NewFields()

	// Test String
	f.String("name", "John")

	// Test Int
	f.Int("age", 30)

	// Test Int64
	f.Int64("id", 1234567890123)

	// Test Float64
	f.Float64("score", 98.5)

	// Test Bool
	f.Bool("active", true)

	// Test Time
	testTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	f.Time("created", testTime)

	// Test Error
	testErr := errors.New("test error")
	f.Error("error", testErr)

	m := f.ToMap()

	if m["name"] != "John" {
		t.Error("String field not set correctly")
	}
	if m["age"] != 30 {
		t.Error("Int field not set correctly")
	}
	if m["id"] != int64(1234567890123) {
		t.Error("Int64 field not set correctly")
	}
	if m["score"] != 98.5 {
		t.Error("Float64 field not set correctly")
	}
	if m["active"] != true {
		t.Error("Bool field not set correctly")
	}
	if m["created"] != testTime.Format(time.RFC3339) {
		t.Error("Time field not set correctly")
	}
	if m["error"] != "test error" {
		t.Error("Error field not set correctly")
	}
}

// TestFieldsAny tests the Any method
func TestFieldsAny(t *testing.T) {
	f := NewFields()

	type CustomStruct struct {
		Name string
		Value int
	}

	f.Any("custom", CustomStruct{Name: "test", Value: 42})
	f.Any("slice", []int{1, 2, 3})
	f.Any("map", map[string]int{"a": 1})

	m := f.ToMap()

	if _, ok := m["custom"].(CustomStruct); !ok {
		t.Error("Any should preserve struct type")
	}
	if _, ok := m["slice"].([]int); !ok {
		t.Error("Any should preserve slice type")
	}
}

// TestFieldsSet tests the Set method (alias for String)
func TestFieldsSet(t *testing.T) {
	f := NewFields()
	f.Set("key", "value")

	m := f.ToMap()
	if m["key"] != "value" {
		t.Error("Set should work like String")
	}
}

// TestFieldsAdd tests the Add method (alias for String)
func TestFieldsAdd(t *testing.T) {
	f := NewFields()
	f.Add("key", "value")

	m := f.ToMap()
	if m["key"] != "value" {
		t.Error("Add should work like String")
	}
}

// TestFieldsWith tests the With method (alias for String)
func TestFieldsWith(t *testing.T) {
	f := NewFields()
	f.With("key", "value")

	m := f.ToMap()
	if m["key"] != "value" {
		t.Error("With should work like String")
	}
}

// TestFieldsMerge tests merging two Fields
func TestFieldsMerge(t *testing.T) {
	f1 := NewFields().String("key1", "value1")
	f2 := NewFields().String("key2", "value2")

	f1.Merge(f2)

	m := f1.ToMap()
	if m["key1"] != "value1" || m["key2"] != "value2" {
		t.Error("Merge should combine fields from both")
	}
}

// TestFieldsClone tests cloning Fields
func TestFieldsClone(t *testing.T) {
	f1 := NewFields().String("key", "original")
	f2 := f1.Clone()

	// Modify original
	f1.String("key", "modified")

	m1 := f1.ToMap()
	m2 := f2.ToMap()

	if m1["key"] == m2["key"] {
		t.Error("Clone should create independent copy")
	}
}

// TestFieldHelpers tests standalone field helper functions
func TestFieldHelpers(t *testing.T) {
	// Test Field (generic)
	f := Field("key", "value")
	if f["key"] != "value" {
		t.Error("Field helper should create map with key-value")
	}

	// Test StringField
	sf := StringField("name", "test")
	if sf["name"] != "test" {
		t.Error("StringField should create string field")
	}

	// Test IntField
	inf := IntField("count", 42)
	if inf["count"] != 42 {
		t.Error("IntField should create int field")
	}

	// Test ErrorField
	err := errors.New("test error")
	ef := ErrorField("error", err)
	if ef["error"] != "test error" {
		t.Error("ErrorField should create error field")
	}

	// Test TimeField
	testTime2 := time.Now()
	tf := TimeField("created", testTime2)
	if tf["created"] != testTime2.Format(time.RFC3339) {
		t.Error("TimeField should create time field as RFC3339 string")
	}
}

// TestPooledFields tests the memory-pooled fields implementation
func TestPooledFields(t *testing.T) {
	pf := NewPooledFields()
	defer pf.Release()

	// Test all pooled field methods
	pf.String("name", "test")
	pf.Int("count", 42)
	pf.Int64("big", 9223372036854775807)
	pf.Bool("active", true)
	pf.Float64("rate", 3.14)
	pf.Time("created", time.Now())
	pf.Error("err", errors.New("test error"))

	m := pf.ToMap()

	if m["name"] != "test" {
		t.Error("PooledFields.String failed")
	}
	if m["count"] != 42 {
		t.Error("PooledFields.Int failed")
	}
	if m["big"] != int64(9223372036854775807) {
		t.Error("PooledFields.Int64 failed")
	}
	if m["active"] != true {
		t.Error("PooledFields.Bool failed")
	}
	if m["rate"] != 3.14 {
		t.Error("PooledFields.Float64 failed")
	}
}

// TestWithPooledFields tests the WithPooledFields helper
func TestWithPooledFields(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	WithPooledFields(func(pf *PooledFields) {
		pf.String("key", "value")
		Info.Field("Pooled test", Fields(pf.ToMap()))
	})

	output := buf.String()
	if !strings.Contains(output, "Pooled test") {
		t.Error("WithPooledFields should enable logging")
	}
}

// TestFieldsErrorWithNil tests Error field with nil error
func TestFieldsErrorWithNil(t *testing.T) {
	f := NewFields()
	f.Error("error", nil)

	m := f.ToMap()
	if m["error"] != nil {
		t.Error("Error with nil should set nil")
	}
}

// TestPooledFieldsError tests the PooledFields.Error method
func TestPooledFieldsError(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	// Test with non-nil error
	testErr := errors.New("test error")
	WithPooledFields(func(pf *PooledFields) {
		pf.Error("err", testErr)
		Info.Field("Error test", Fields(pf.ToMap()))
	})

	output := buf.String()
	if !strings.Contains(output, "test error") {
		t.Error("PooledFields.Error should include error message")
	}

	buf.Reset()

	// Test with nil error
	WithPooledFields(func(pf *PooledFields) {
		pf.Error("err", nil)
		Info.Field("Nil error test", Fields(pf.ToMap()))
	})
}
