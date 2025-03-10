package golog

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// Plain is a function that prints a plain-text log message
// on the console using a very simple format with color-coded severity levels.
func Plain(severity, message string, optionalParams ...string) {

	if severity != "info" && severity != "warning" && severity != "error" && severity != "debug" {
		severity = "info"
	}

	// Severity level is case insensitive, so we convert it to lowercase.
	severity = strings.ToLower(severity)

	var colorCode string
	switch severity {
	case "info":
		colorCode = "\033[32m" // Green
	case "warning":
		colorCode = "\033[33m" // Yellow
	case "error":
		colorCode = "\033[31m" // Red
	case "debug":
		colorCode = "\033[34m" // Blue
	default:
		colorCode = ""
	}

	resetCode := "\033[0m" // Reset color

	if runtime.GOOS == "windows" {
		// Windows doesn't directly support ANSI escape codes;  you'd need a Windows-specific library
		colorCode = ""
		resetCode = ""
	}

	appName := ""
	appVersion := ""

	// Handle optional parameters based on their count.
	switch len(optionalParams) {
	case 1:
		appName = optionalParams[0]
	case 2:
		appName = optionalParams[0]
		appVersion = optionalParams[1]
	}

	// Console output format:
	// {UTC TIME} | {LOGGING LEVEL} | {APPLICATION NAME} {APPLICATION VERSION}: {MESSAGE}

	fmt.Printf("%s | %s%s-7s%s | %s %s: %s\n", time.Now().UTC().Format("2006-01-02 15:04:05"), colorCode, severity, resetCode, appName, appVersion, message)

}

// JSON is a function that prints a JSON log message.
func JSON(severity, message string, optionalParams ...string) {

	logData := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"severity":  severity,
		"message":   message,
	}

	appName := ""
	appVersion := ""
	switch len(optionalParams) {
	case 1:
		appName = optionalParams[0]
	case 2:
		appName = optionalParams[0]
		appVersion = optionalParams[1]
	}

	if appName != "" {
		logData["application_name"] = appName
	}
	if appVersion != "" {
		logData["application_version"] = appVersion
	}

	jsonData, err := json.Marshal(logData)
	if err != nil {
		fmt.Printf("Error marshaling JSON log: %v\n", err)
		return
	}

	fmt.Println(string(jsonData))

}
