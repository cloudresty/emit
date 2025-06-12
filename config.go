package emit

import (
	"io"
	"os"
	"strings"
)

// initFromEnvironment initializes logger settings from environment variables
func initFromEnvironment() {

	// Check environment variable for format override
	if logFormat := os.Getenv("EMIT_FORMAT"); logFormat != "" {

		switch strings.ToLower(logFormat) {

		case "plain", "text", "console", "development", "dev":
			defaultLogger.format = PLAIN_FORMAT

		case "json", "production", "prod":
			defaultLogger.format = JSON_FORMAT

		default:
			// Invalid value, stick with JSON default
			defaultLogger.format = JSON_FORMAT

		}

	}

	// Also check for log level from environment
	if logLevel := os.Getenv("EMIT_LEVEL"); logLevel != "" {
		defaultLogger.level = ParseLogLevel(logLevel)
	}

	// Check for caller information setting
	if showCaller := os.Getenv("EMIT_SHOW_CALLER"); showCaller != "" {
		defaultLogger.showCaller = strings.ToLower(showCaller) == "true" || showCaller == "1"
	}

	// Check for sensitive data masking setting
	if sensitiveMode := os.Getenv("EMIT_MASK_SENSITIVE"); sensitiveMode != "" {

		switch strings.ToLower(sensitiveMode) {

		case "false", "0", "no", "off", "show":
			defaultLogger.sensitiveMode = SHOW_SENSITIVE

		case "true", "1", "yes", "on", "mask":
			defaultLogger.sensitiveMode = MASK_SENSITIVE

		default:
			defaultLogger.sensitiveMode = MASK_SENSITIVE // Default to masking

		}

	}

	// Check for PII data masking setting
	if piiMode := os.Getenv("EMIT_MASK_PII"); piiMode != "" {

		switch strings.ToLower(piiMode) {

		case "false", "0", "no", "off", "show":
			defaultLogger.piiMode = SHOW_PII

		case "true", "1", "yes", "on", "mask":
			defaultLogger.piiMode = MASK_PII

		default:
			defaultLogger.piiMode = MASK_PII // Default to masking
		}

	}

	// Allow custom mask string
	if maskString := os.Getenv("EMIT_MASK_STRING"); maskString != "" {
		defaultLogger.maskString = maskString
	}

	// Allow custom PII mask string
	if piiMaskString := os.Getenv("EMIT_PII_MASK_STRING"); piiMaskString != "" {
		defaultLogger.piiMaskString = piiMaskString
	}

	// PHASE 3: Check for timestamp precision setting
	if timestampPrecision := os.Getenv("EMIT_TIMESTAMP_PRECISION"); timestampPrecision != "" {
		SetTimestampPrecisionConfig(ParseTimestampPrecision(timestampPrecision))
	}
}

// SetComponent sets the component name for the default logger
func SetComponent(component string) {
	if defaultLogger != nil {
		defaultLogger.component = component
	}
}

// SetVersion sets the version for the default logger
func SetVersion(version string) {
	if defaultLogger != nil {
		defaultLogger.version = version
	}
}

// SetLevel sets the log level for the default logger
func SetLevel(level string) {
	if defaultLogger != nil {
		defaultLogger.level = ParseLogLevel(level)
	}
}

// SetShowCaller enables or disables caller information
func SetShowCaller(show bool) {
	if defaultLogger != nil {
		defaultLogger.showCaller = show
	}
}

// SetFormat sets the output format (JSON or Plain)
func SetFormat(format string) {

	if defaultLogger != nil {

		switch strings.ToLower(format) {

		case "plain", "text", "console":
			defaultLogger.format = PLAIN_FORMAT

		case "json":
			defaultLogger.format = JSON_FORMAT

		default:
			defaultLogger.format = JSON_FORMAT

		}

	}

}

// SetPlainFormat switches to plain text output for development
func SetPlainFormat() {
	SetFormat("plain")
}

// SetJSONFormat switches to JSON output for production
func SetJSONFormat() {
	SetFormat("json")
}

// SetSensitiveMode sets whether to mask sensitive data
func SetSensitiveMode(mode string) {

	if defaultLogger != nil {

		switch strings.ToLower(mode) {

		case "show", "false", "0", "no", "off":
			defaultLogger.sensitiveMode = SHOW_SENSITIVE

		case "mask", "true", "1", "yes", "on":
			defaultLogger.sensitiveMode = MASK_SENSITIVE

		default:
			defaultLogger.sensitiveMode = MASK_SENSITIVE

		}

	}

}

// ShowSensitiveData disables masking of sensitive fields (not recommended for production)
func ShowSensitiveData() {
	SetSensitiveMode("show")
}

// MaskSensitiveData enables masking of sensitive fields (default and recommended)
func MaskSensitiveData() {
	SetSensitiveMode("mask")
}

// SetOutput sets the output writer for the default logger
func SetOutput(writer io.Writer) {
	if defaultLogger != nil {
		defaultLogger.writer = writer
	}
}

// SetOutputToDiscard redirects output to discard for benchmarking
func SetOutputToDiscard() {
	if defaultLogger != nil {
		defaultLogger.writer = io.Discard
	}
}

// SetMaskString sets the string used to mask sensitive data
func SetMaskString(mask string) {
	if defaultLogger != nil && mask != "" {
		defaultLogger.maskString = mask
	}
}

// AddSensitiveField adds a custom field pattern to be masked
func AddSensitiveField(field string) {
	if defaultLogger != nil {
		defaultLogger.sensitiveFields = append(defaultLogger.sensitiveFields, strings.ToLower(field))
	}
}

// SetSensitiveFields replaces the default sensitive field patterns
func SetSensitiveFields(fields []string) {
	if defaultLogger != nil {
		var lowerFields []string
		for _, field := range fields {
			lowerFields = append(lowerFields, strings.ToLower(field))
		}
		defaultLogger.sensitiveFields = lowerFields
	}
}

// SetPIIMode sets whether to mask PII data
func SetPIIMode(mode string) {
	if defaultLogger != nil {
		switch strings.ToLower(mode) {
		case "show", "false", "0", "no", "off":
			defaultLogger.piiMode = SHOW_PII
		case "mask", "true", "1", "yes", "on":
			defaultLogger.piiMode = MASK_PII
		default:
			defaultLogger.piiMode = MASK_PII
		}
	}
}

// ShowPIIData disables masking of PII fields (not recommended for production)
func ShowPIIData() {
	SetPIIMode("show")
}

// MaskPIIData enables masking of PII fields (default and recommended)
func MaskPIIData() {
	SetPIIMode("mask")
}

// SetPIIMaskString sets the string used to mask PII data
func SetPIIMaskString(mask string) {
	if defaultLogger != nil && mask != "" {
		defaultLogger.piiMaskString = mask
	}
}

// AddPIIField adds a custom field pattern to be masked as PII
func AddPIIField(field string) {
	if defaultLogger != nil {
		defaultLogger.piiFields = append(defaultLogger.piiFields, strings.ToLower(field))
	}
}

// SetPIIFields replaces the default PII field patterns
func SetPIIFields(fields []string) {
	if defaultLogger != nil {
		var lowerFields []string
		for _, field := range fields {
			lowerFields = append(lowerFields, strings.ToLower(field))
		}
		defaultLogger.piiFields = lowerFields
	}
}

// SetAllMasking enables or disables both sensitive and PII masking
func SetAllMasking(enabled bool) {
	if enabled {
		MaskSensitiveData()
		MaskPIIData()
	} else {
		ShowSensitiveData()
		ShowPIIData()
	}
}

// SetDevelopmentMode disables all masking for development
func SetDevelopmentMode() {
	SetAllMasking(false)
	SetPlainFormat()
	SetLevel("debug")
	SetShowCaller(true)
}

// SetProductionMode enables all masking for production
func SetProductionMode() {
	SetAllMasking(true)
	SetJSONFormat()
	SetLevel("info")
	SetShowCaller(false)
}

// PHASE 3: Timestamp precision configuration

// ParseTimestampPrecision parses timestamp precision from string
func ParseTimestampPrecision(precision string) TimestampPrecision {
	switch strings.ToLower(precision) {
	case "nanosecond", "nano", "ns":
		return NanosecondPrecision
	case "microsecond", "micro", "us":
		return MicrosecondPrecision
	case "millisecond", "milli", "ms":
		return MillisecondPrecision
	case "second", "sec", "s":
		return SecondPrecision
	default:
		return NanosecondPrecision // Default to highest precision
	}
}

// SetTimestampPrecisionConfig sets the timestamp precision for the logging system
func SetTimestampPrecisionConfig(precision TimestampPrecision) {
	SetTimestampPrecision(precision)
}

// GetTimestampPrecisionConfig returns the current timestamp precision
func GetTimestampPrecisionConfig() TimestampPrecision {
	return GetTimestampPrecision()
}
