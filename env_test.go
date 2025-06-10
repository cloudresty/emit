package emit

import (
	"os"
	"testing"
)

// TestEnvironmentVariableConfiguration tests that environment variables are properly loaded
func TestEnvironmentVariableConfiguration(t *testing.T) {
	// Save original environment
	originalEnv := map[string]string{
		"EMIT_FORMAT":          os.Getenv("EMIT_FORMAT"),
		"EMIT_LEVEL":           os.Getenv("EMIT_LEVEL"),
		"EMIT_SHOW_CALLER":     os.Getenv("EMIT_SHOW_CALLER"),
		"EMIT_MASK_SENSITIVE":  os.Getenv("EMIT_MASK_SENSITIVE"),
		"EMIT_MASK_PII":        os.Getenv("EMIT_MASK_PII"),
		"EMIT_MASK_STRING":     os.Getenv("EMIT_MASK_STRING"),
		"EMIT_PII_MASK_STRING": os.Getenv("EMIT_PII_MASK_STRING"),
	}

	// Restore environment after test
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// Test format configuration
	tests := []struct {
		name     string
		envVar   string
		envValue string
		check    func(*Logger) bool
	}{
		{
			name:     "Plain format",
			envVar:   "EMIT_FORMAT",
			envValue: "plain",
			check:    func(l *Logger) bool { return l.format == PLAIN_FORMAT },
		},
		{
			name:     "JSON format",
			envVar:   "EMIT_FORMAT",
			envValue: "json",
			check:    func(l *Logger) bool { return l.format == JSON_FORMAT },
		},
		{
			name:     "Development format",
			envVar:   "EMIT_FORMAT",
			envValue: "development",
			check:    func(l *Logger) bool { return l.format == PLAIN_FORMAT },
		},
		{
			name:     "Debug level",
			envVar:   "EMIT_LEVEL",
			envValue: "debug",
			check:    func(l *Logger) bool { return l.level == DEBUG },
		},
		{
			name:     "Error level",
			envVar:   "EMIT_LEVEL",
			envValue: "error",
			check:    func(l *Logger) bool { return l.level == ERROR },
		},
		{
			name:     "Show caller true",
			envVar:   "EMIT_SHOW_CALLER",
			envValue: "true",
			check:    func(l *Logger) bool { return l.showCaller == true },
		},
		{
			name:     "Show caller 1",
			envVar:   "EMIT_SHOW_CALLER",
			envValue: "1",
			check:    func(l *Logger) bool { return l.showCaller == true },
		},
		{
			name:     "Show caller false",
			envVar:   "EMIT_SHOW_CALLER",
			envValue: "false",
			check:    func(l *Logger) bool { return l.showCaller == false },
		},
		{
			name:     "Mask sensitive false",
			envVar:   "EMIT_MASK_SENSITIVE",
			envValue: "false",
			check:    func(l *Logger) bool { return l.sensitiveMode == SHOW_SENSITIVE },
		},
		{
			name:     "Mask PII false",
			envVar:   "EMIT_MASK_PII",
			envValue: "false",
			check:    func(l *Logger) bool { return l.piiMode == SHOW_PII },
		},
		{
			name:     "Custom mask string",
			envVar:   "EMIT_MASK_STRING",
			envValue: "[REDACTED]",
			check:    func(l *Logger) bool { return l.maskString == "[REDACTED]" },
		},
		{
			name:     "Custom PII mask string",
			envVar:   "EMIT_PII_MASK_STRING",
			envValue: "[PII_HIDDEN]",
			check:    func(l *Logger) bool { return l.piiMaskString == "[PII_HIDDEN]" },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Clear all environment variables
			for key := range originalEnv {
				os.Unsetenv(key)
			}

			// Set the specific environment variable for this test
			os.Setenv(test.envVar, test.envValue)

			// Create a new logger that will read from environment
			testLogger := &Logger{
				level:           INFO,
				writer:          os.Stdout,
				showCaller:      false,
				format:          JSON_FORMAT,
				sensitiveMode:   MASK_SENSITIVE,
				piiMode:         MASK_PII,
				sensitiveFields: defaultSensitiveFields,
				piiFields:       defaultPIIFields,
				maskString:      "***MASKED***",
				piiMaskString:   "***PII***",
			}

			// Simulate environment initialization
			originalLogger := defaultLogger
			defaultLogger = testLogger
			initFromEnvironment()
			defaultLogger = originalLogger

			// Check if the configuration was applied correctly
			if !test.check(testLogger) {
				t.Errorf("Environment variable %s=%s was not applied correctly", test.envVar, test.envValue)
			}
		})
	}
}

// TestInvalidEnvironmentValues tests that invalid environment values fall back to defaults
func TestInvalidEnvironmentValues(t *testing.T) {
	// Save original environment
	originalEnv := map[string]string{
		"EMIT_FORMAT":         os.Getenv("EMIT_FORMAT"),
		"EMIT_LEVEL":          os.Getenv("EMIT_LEVEL"),
		"EMIT_MASK_SENSITIVE": os.Getenv("EMIT_MASK_SENSITIVE"),
		"EMIT_MASK_PII":       os.Getenv("EMIT_MASK_PII"),
	}

	// Restore environment after test
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// Test invalid values
	invalidTests := []struct {
		envVar       string
		invalidValue string
		checkDefault func(*Logger) bool
	}{
		{
			envVar:       "EMIT_FORMAT",
			invalidValue: "invalid_format",
			checkDefault: func(l *Logger) bool { return l.format == JSON_FORMAT }, // Should default to JSON
		},
		{
			envVar:       "EMIT_LEVEL",
			invalidValue: "invalid_level",
			checkDefault: func(l *Logger) bool { return l.level == INFO }, // Should default to INFO
		},
		{
			envVar:       "EMIT_MASK_SENSITIVE",
			invalidValue: "maybe",
			checkDefault: func(l *Logger) bool { return l.sensitiveMode == MASK_SENSITIVE }, // Should default to masking
		},
		{
			envVar:       "EMIT_MASK_PII",
			invalidValue: "perhaps",
			checkDefault: func(l *Logger) bool { return l.piiMode == MASK_PII }, // Should default to masking
		},
	}

	for _, test := range invalidTests {
		t.Run("Invalid "+test.envVar, func(t *testing.T) {
			// Clear all environment variables
			for key := range originalEnv {
				os.Unsetenv(key)
			}

			// Set invalid value
			os.Setenv(test.envVar, test.invalidValue)

			// Create a new logger
			testLogger := &Logger{
				level:           INFO,
				writer:          os.Stdout,
				showCaller:      false,
				format:          JSON_FORMAT,
				sensitiveMode:   MASK_SENSITIVE,
				piiMode:         MASK_PII,
				sensitiveFields: defaultSensitiveFields,
				piiFields:       defaultPIIFields,
				maskString:      "***MASKED***",
				piiMaskString:   "***PII***",
			}

			// Simulate environment initialization
			originalLogger := defaultLogger
			defaultLogger = testLogger
			initFromEnvironment()
			defaultLogger = originalLogger

			// Check if it falls back to default
			if !test.checkDefault(testLogger) {
				t.Errorf("Invalid environment variable %s=%s should fall back to default", test.envVar, test.invalidValue)
			}
		})
	}
}

// TestCustomFieldPatterns tests adding custom sensitive and PII field patterns
func TestCustomFieldPatterns(t *testing.T) {
	// Save original logger
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	// Reset to a clean state
	defaultLogger = &Logger{
		level:           INFO,
		writer:          os.Stdout,
		format:          JSON_FORMAT,
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
	}

	// Test adding custom sensitive field
	AddSensitiveField("custom_secret")

	found := false
	for _, field := range defaultLogger.sensitiveFields {
		if field == "custom_secret" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Custom sensitive field 'custom_secret' was not added")
	}

	// Test adding custom PII field
	AddPIIField("employee_id")

	found = false
	for _, field := range defaultLogger.piiFields {
		if field == "employee_id" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Custom PII field 'employee_id' was not added")
	}

	// Test setting custom field arrays
	customSensitive := []string{"secret1", "secret2"}
	SetSensitiveFields(customSensitive)

	if len(defaultLogger.sensitiveFields) != 2 {
		t.Errorf("Expected 2 sensitive fields, got %d", len(defaultLogger.sensitiveFields))
	}

	customPII := []string{"pii1", "pii2", "pii3"}
	SetPIIFields(customPII)

	if len(defaultLogger.piiFields) != 3 {
		t.Errorf("Expected 3 PII fields, got %d", len(defaultLogger.piiFields))
	}
}
