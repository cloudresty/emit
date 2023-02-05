package log

import (
	"fmt"
	"strings"
	"time"
)

// Plain is a function that prints a plain-text log message
// on the console using a very simple format.

func Plain(severity, applicationName, applicationVersion, message string) {

	// Check if the severity level is valid,
	// it has to be 'info', 'warning', 'error' or 'debug'.
	// If it is not valid, then the default value is 'info'.

	if severity != "info" && severity != "warning" && severity != "error" && severity != "debug" {
		severity = "info"
	}

	// Severity level is case insensitive, so we convert it to lowercase.
	severity = strings.ToLower(severity)

	// Format the severity level to match the same length.
	switch severity {
	case "info":
		severity = "INFO    |"
	case "warning":
		severity = "WARNING |"
	case "error":
		severity = "ERROR   |"
	case "debug":
		severity = "DEBUG   |"
	default:
		severity = "INFO    |"
	}

	// Console output format:
	// {UTC TIME} | {LOGGING LEVEL} | {APPLICATION NAME} {APPLICATION VERSION}: {MESSAGE}

	fmt.Printf("%s | %s %s %s: %s\n", time.Now().UTC().Format("2006-01-02 15:04:05"), severity, applicationName, applicationVersion, message)

}
