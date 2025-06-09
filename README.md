# GoLog - Kubernetes-Ready JSON Logging for Go

A lightweight, structured logging library for Go applications optimized for Kubernetes environments. GoLog provides JSON logging by default with support for structured fields, making it perfect for cloud-native applications.

&nbsp;

## Features

- **JSON-first logging** - Kubernetes-optimized structured logs by default
- **Environment-based configuration** - Control format via environment variables
- **Structured logging** - Add contextual fields to your logs
- **Multiple log levels** - DEBUG, INFO, WARN, ERROR
- **Caller information** - Optional file, line, and function tracking
- **Global configuration** - Set component, version, and log level globally
- **Plain text fallback** - Colored console output for development
- **Zero dependencies** - Uses only Go standard library
- **Backward compatibility** - Works with existing code

&nbsp;

## Installation

```bash
go get github.com/cloudresty/golog
```

&nbsp;

## Quick Start

```go
package main

import "github.com/cloudresty/golog"

func main() {
    // Simple logging (JSON format by default)
    golog.Info("Application starting")
    golog.Error("Something went wrong")

    // With component and version
    golog.Info("User authenticated", "auth-service", "v1.2.3")
}
```

Output (JSON format):

```json
{"timestamp":"2025-06-09T10:30:45.123456789Z","level":"info","msg":"Application starting"}
{"timestamp":"2025-06-09T10:30:45.124567890Z","level":"error","msg":"Something went wrong"}
{"timestamp":"2025-06-09T10:30:45.125678901Z","level":"info","msg":"User authenticated","component":"auth-service","version":"v1.2.3"}
```

&nbsp;

## Environment-Based Configuration

GoLog automatically configures itself based on environment variables:

&nbsp;

### Environment Variables

- `GOLOG_FORMAT`: Controls output format
  - `json`, `production`, `prod` - JSON format (default)
  - `plain`, `text`, `console`, `development`, `dev` - Plain text format
- `GOLOG_LEVEL`: Sets log level (`debug`, `info`, `warn`, `error`)
- `GOLOG_SHOW_CALLER`: Enable caller info (`true`, `1` to enable)

&nbsp;

### Local Development Setup

```bash
# Set environment for development
export GOLOG_FORMAT=plain
export GOLOG_LEVEL=debug
export GOLOG_SHOW_CALLER=true

# Run your application
go run main.go
```

Output (Plain format):

```plaintext
2025-06-09 10:30:45 | info    | auth-service v1.2.3: Application starting
2025-06-09 10:30:45 | error   | auth-service v1.2.3: Something went wrong
2025-06-09 10:30:45 | info    | auth-service v1.2.3: User authenticated
```

&nbsp;

### Production/Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  template:
    spec:
      containers:
      - name: my-app
        image: my-app:latest
        env:
        - name: GOLOG_LEVEL
          value: "info"
        # GOLOG_FORMAT not set = defaults to JSON
```

&nbsp;

## Configuration

&nbsp;

### Global Configuration

```go
package main

import "github.com/cloudresty/golog"

func main() {
    // Configure the global logger
    golog.SetComponent("user-service")
    golog.SetVersion("v2.1.0")
    golog.SetLevel("debug")
    golog.SetShowCaller(true)

    // All subsequent logs will include these settings
    golog.Info("Service configured")
    golog.Debug("Debug information")
}
```

&nbsp;

### Runtime Format Switching

```go
// Switch to plain format for development
golog.SetPlainFormat()

// Switch back to JSON for production
golog.SetJSONFormat()

// Or use generic SetFormat
golog.SetFormat("plain")  // or "json"
```

&nbsp;

## Logging Levels

```go
// Basic logging functions
golog.Debug("Detailed debugging information")
golog.Info("General information")
golog.Warning("Warning message")  // or golog.Warn()
golog.Error("Error occurred")

// Using the generic Log function
golog.Log("info", "Generic log message")
golog.Log("error", "Generic error message", "my-app", "v1.0.0")
```

&nbsp;

## Structured Logging with Fields

Add contextual information to your logs using the `WithFields` functions:

```go
// Info with fields
golog.InfoWithFields("User login successful", map[string]any{
    "user_id":    12345,
    "username":   "john.doe",
    "ip_address": "192.168.1.100",
    "duration":   "250ms",
})

// Error with fields
golog.ErrorWithFields("Database connection failed", map[string]any{
    "database": "users_db",
    "host":     "db.example.com",
    "port":     5432,
    "error":    "connection timeout",
    "retry":    3,
})

// Warning with fields
golog.WarnWithFields("Rate limit approaching", map[string]any{
    "current_requests": 850,
    "limit":           1000,
    "window":          "1m",
})
```

JSON Output:

```json
{"timestamp":"2025-06-09T10:30:45.123456789Z","level":"info","msg":"User login successful","fields":{"duration":"250ms","ip_address":"192.168.1.100","user_id":12345,"username":"john.doe"}}
{"timestamp":"2025-06-09T10:30:45.124567890Z","level":"error","msg":"Database connection failed","fields":{"database":"users_db","error":"connection timeout","host":"db.example.com","port":5432,"retry":3}}
```

Plain Text Output:

```plaintext
2025-06-09 10:30:45 | info    | my-app v1.0.0: User login successful [user_id=12345 username=john.doe ip_address=192.168.1.100 duration=250ms]
2025-06-09 10:30:45 | error   | my-app v1.0.0: Database connection failed [database=users_db host=db.example.com port=5432 error=connection timeout retry=3]
```

&nbsp;

## Development Workflows

&nbsp;

### Docker Compose for Development

```yaml
version: '3'
services:
  app:
    build: .
    environment:
      - GOLOG_FORMAT=plain
      - GOLOG_LEVEL=debug
      - GOLOG_SHOW_CALLER=true
```

&nbsp;

### Makefile for Easy Switching

```makefile
.PHONY: dev prod

dev:
    GOLOG_FORMAT=plain GOLOG_LEVEL=debug GOLOG_SHOW_CALLER=true go run main.go

prod:
    GOLOG_FORMAT=json GOLOG_LEVEL=info go run main.go

test:
    GOLOG_FORMAT=plain GOLOG_LEVEL=debug go test ./...
```

&nbsp;

### IDE/Editor Setup

VS Code `.vscode/settings.json`:

```json
{
    "go.testEnvVars": {
        "GOLOG_FORMAT": "plain",
        "GOLOG_LEVEL": "debug"
    }
}
```

&nbsp;

## Kubernetes Integration

&nbsp;

### Example Application

```go
package main

import (
    "os"
    "github.com/cloudresty/golog"
)

func main() {
    // Configure from environment variables
    golog.SetComponent("my-app")
    golog.SetVersion(os.Getenv("APP_VERSION"))

    golog.Info("Application started")

    // Your application logic here
    handleRequest()
}

func handleRequest() {
    golog.InfoWithFields("Processing request", map[string]any{
        "request_id": "req-123",
        "method":     "GET",
        "path":       "/api/users",
    })
}
```

&nbsp;

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  template:
    spec:
      containers:
      - name: my-app
        image: my-app:latest
        env:
        - name: APP_VERSION
          value: "v1.2.3"
        - name: GOLOG_LEVEL
          value: "info"
        # Production uses JSON format by default
```

&nbsp;

## Advanced Usage

&nbsp;

### Custom Log Levels

```go
// Set different log levels
golog.SetLevel("debug")  // Shows all logs
golog.SetLevel("info")   // Shows info, warn, error
golog.SetLevel("warn")   // Shows warn, error
golog.SetLevel("error")  // Shows only error
```

&nbsp;

### Caller Information

```go
// Enable caller information for debugging
golog.SetShowCaller(true)

golog.Info("This will include file and line info")
```

&nbsp;

### Force Specific Format

```go
// Force JSON output regardless of environment
golog.JSON("info", "Always JSON", "app", "v1.0.0")

// Force plain output regardless of environment
golog.Plain("info", "Always plain text", "app", "v1.0.0")
```

&nbsp;

## Migration from Other Loggers

&nbsp;

### From Standard Log Package

```go
// Old
log.Printf("User %s logged in", username)

// New
golog.InfoWithFields("User logged in", map[string]any{
    "username": username,
})
```

&nbsp;

### From Logrus

```go
// Old
logrus.WithFields(logrus.Fields{
    "user_id": 123,
}).Info("User action")

// New
golog.InfoWithFields("User action", map[string]any{
    "user_id": 123,
})
```

&nbsp;

## Best Practices

1. **Use environment variables** for configuration instead of hardcoding
2. **Set global configuration early** in your application startup
3. **Use structured logging** with fields for better observability
4. **Include correlation IDs** in your log fields for tracing
5. **Use appropriate log levels** - avoid debug logs in production
6. **Include version and component** information for better debugging
7. **Use plain format for development**, JSON for production

&nbsp;

## Performance

GoLog is designed to be lightweight and fast:

- Minimal allocations for log level filtering
- Efficient JSON marshaling
- Early return for filtered log levels
- No external dependencies
- Zero-allocation environment variable reads (cached in init)

&nbsp;

---

Made with ♥️ by [Cloudresty](https://cloudresty.com).
