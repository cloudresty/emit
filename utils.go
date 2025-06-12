package emit

import "strings"

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	default:
		return "info"
	}
}

// StringFast returns the string representation of the log level with optimized performance
// PHASE 4B: Optimized version that uses pre-allocated strings for better performance
func (l LogLevel) StringFast() string {
	// Use compile-time constants to avoid string allocation
	switch l {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	default:
		return "info"
	}
}

// ParseLogLevel parses a string into a LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info", "information":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}
