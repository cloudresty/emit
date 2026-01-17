package emit

import (
	"strings"
	"testing"
)

// TestGetUltraFastTimestampPrecisions tests different timestamp precisions
func TestGetUltraFastTimestampPrecisions(t *testing.T) {
	// Test nanosecond precision
	SetTimestampPrecision(NanosecondPrecision)
	ts := GetUltraFastTimestamp()
	if len(ts) == 0 {
		t.Error("Nanosecond timestamp should not be empty")
	}

	// Test microsecond precision
	SetTimestampPrecision(MicrosecondPrecision)
	ts = GetUltraFastTimestamp()
	if len(ts) == 0 {
		t.Error("Microsecond timestamp should not be empty")
	}

	// Test millisecond precision
	SetTimestampPrecision(MillisecondPrecision)
	ts = GetUltraFastTimestamp()
	if len(ts) == 0 {
		t.Error("Millisecond timestamp should not be empty")
	}

	// Test second precision
	SetTimestampPrecision(SecondPrecision)
	ts = GetUltraFastTimestamp()
	if len(ts) == 0 {
		t.Error("Second timestamp should not be empty")
	}

	// Reset to default
	SetTimestampPrecision(NanosecondPrecision)
}

// TestTimestampFormat tests that timestamps are in expected format
func TestTimestampFormat(t *testing.T) {
	SetTimestampPrecision(NanosecondPrecision)
	ts := GetUltraFastTimestamp()

	// Should contain date separator
	if !strings.Contains(ts, "-") {
		t.Error("Timestamp should contain date separators")
	}

	// Should contain time separator
	if !strings.Contains(ts, ":") {
		t.Error("Timestamp should contain time separators")
	}

	// Should contain T separator
	if !strings.Contains(ts, "T") {
		t.Error("Timestamp should contain T separator")
	}

	// Should end with Z for UTC
	if !strings.HasSuffix(ts, "Z") {
		t.Error("Timestamp should end with Z for UTC")
	}
}

// TestLogLevelString tests LogLevel.String method
func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, "debug"},
		{INFO, "info"},
		{WARN, "warn"},
		{ERROR, "error"},
		{LogLevel(99), "info"}, // Unknown level defaults to info
	}

	for _, tt := range tests {
		result := tt.level.String()
		if result != tt.expected {
			t.Errorf("LogLevel(%d).String() = %s, want %s", tt.level, result, tt.expected)
		}
	}
}

// TestLogLevelStringFast tests LogLevel.StringFast method
func TestLogLevelStringFast(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, "debug"},
		{INFO, "info"},
		{WARN, "warn"},
		{ERROR, "error"},
		{LogLevel(99), "info"}, // Unknown level defaults to info
	}

	for _, tt := range tests {
		result := tt.level.StringFast()
		if result != tt.expected {
			t.Errorf("LogLevel(%d).StringFast() = %s, want %s", tt.level, result, tt.expected)
		}
	}
}

