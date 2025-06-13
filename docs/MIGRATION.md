# Migration Guide

Complete guide for migrating from other logging libraries to emit's elegant API.

## Why Migrate to Emit?

### The Problem with Traditional Loggers

```go
// Traditional loggers have serious issues:

// 1. Easy to accidentally expose sensitive data
log.Printf("User %s logged in with password %s", username, password) // DANGEROUS!

// 2. Verbose and error-prone syntax
logrus.WithFields(logrus.Fields{
    "user_id": 123,
    "email": maskEmail(email),  // Manual masking required!
    "action": "login",
}).Info("User action")

// 3. Complex setup for security
logger.Info("Payment",
    zap.String("email", maskPII(email)),        // Manual masking
    zap.String("card", maskSensitive(card)),    // Manual masking
    zap.Int("amount", amount))
```

### The Emit Solution

```go
// Emit solves all these problems:

// 1. Automatic security - impossible to accidentally expose data
emit.Info.KeyValue("User logged in",
    "username", username,      // Auto-protected if PII
    "password", password)      // Auto-masked

// 2. Clean, simple API
emit.Info.Field("User action",
    emit.NewFields().
        Int("user_id", 123).
        String("email", email).     // Auto-masked
        String("action", "login"))

// 3. Zero-config security
emit.Info.KeyValue("Payment processed",
    "email", email,            // Auto-masked
    "card", card,              // Auto-masked
    "amount", amount)          // Safe
```

## Migration from Standard Log Package

### Before (Unsafe)

```go
package main

import (
    "log"
    "fmt"
)

func main() {
    // DANGEROUS: Exposes sensitive data
    username := "john_doe"
    password := "secret123"

    log.Printf("User %s logged in with password %s", username, password)
    log.Printf("Processing payment for %s with card %s", email, cardNumber)

    // No structured data
    log.Printf("Request took %v milliseconds", duration.Milliseconds())
}
```

### After (Secure & Clean)

```go
package main

import (
    "github.com/cloudresty/emit"
)

func main() {
    // SECURE: Automatic data protection
    username := "john_doe"
    password := "secret123"

    // Option 1: Key-Value (Simple)
    emit.Info.KeyValue("User logged in",
        "username", username,      // Auto-protected
        "password", password)      // Auto-masked

    // Option 2: Structured Fields (Elegant)
    emit.Info.Field("Payment processing",
        emit.NewFields().
            String("email", email).         // Auto-masked
            String("card", cardNumber).     // Auto-masked
            Float64("amount", 99.99))

    // Option 3: Zero-Allocation (Performance)
    emit.Info.ZeroAlloc("Request completed",
        emit.ZDuration("duration", duration),
        emit.ZBool("success", true))
}
```

### Standard Log Migration Examples

| **Standard Log** | **Emit Equivalent** |
|------------------|-------------------|
| `log.Printf("User %s", user)` | `emit.Info.KeyValue("User action", "user", user)` |
| `log.Printf("Error: %v", err)` | `emit.Error.KeyValue("Error occurred", "error", err)` |
| `log.Printf("Count: %d", count)` | `emit.Info.KeyValue("Count", "count", count)` |
| `log.Fatal(err)` | `emit.Error.Msg(err.Error()); os.Exit(1)` |

## Migration from Logrus

### Before (Manual Security)

```go
package main

import (
    "github.com/sirupsen/logrus"
)

func main() {
    // Manual security required
    logrus.WithFields(logrus.Fields{
        "user_id":    123,
        "email":      maskEmail(userEmail),     // Manual masking!
        "password":   "[REDACTED]",             // Manual redaction!
        "action":     "login",
        "ip_address": maskIP(clientIP),         // Manual masking!
    }).Info("User login attempt")

    // Complex error handling
    logrus.WithFields(logrus.Fields{
        "error":      err.Error(),
        "service":    "payment",
        "amount":     amount,
        "card":       maskCard(cardNumber),     // Manual masking!
    }).Error("Payment failed")

    // Different levels
    logrus.WithFields(logrus.Fields{
        "query":     query,
        "rows":      rowCount,
        "duration":  duration.String(),
    }).Debug("Database query")
}
```

### After (Automatic Security)

```go
package main

import (
    "github.com/cloudresty/emit"
)

func main() {
    // Automatic security - no manual masking needed!

    // Option 1: Structured Fields (Recommended)
    emit.Info.Field("User login attempt",
        emit.NewFields().
            Int("user_id", 123).
            String("email", userEmail).         // Auto-masked
            String("password", password).       // Auto-masked
            String("action", "login").
            String("ip_address", clientIP))     // Auto-masked

    // Option 2: Key-Value Pairs (Simple)
    emit.Error.KeyValue("Payment failed",
        "error", err.Error(),
        "service", "payment",
        "amount", amount,
        "card", cardNumber)                     // Auto-masked

    // Option 3: Zero-Allocation (Performance)
    emit.Debug.ZeroAlloc("Database query",
        emit.ZString("query", query),
        emit.ZInt("rows", rowCount),
        emit.ZDuration("duration", duration))
}
```

### Logrus Migration Mapping

| **Logrus** | **Emit Equivalent** |
|------------|-------------------|
| `logrus.Info("message")` | `emit.Info.Msg("message")` |
| `logrus.WithField("key", value)` | `emit.Info.KeyValue("message", "key", value)` |
| `logrus.WithFields(fields)` | `emit.Info.Field("message", emit.NewFields()...)` |
| `logrus.Error(err)` | `emit.Error.KeyValue("Error", "error", err)` |
| `logrus.SetLevel(logrus.DebugLevel)` | `emit.SetLevel("debug")` |

## Migration from Zap

### Before (Complex Setup)

```go
package main

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func main() {
    // Complex setup required
    config := zap.NewProductionConfig()
    config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
    logger, _ := config.Build()
    defer logger.Sync()

    // Manual security masking required
    logger.Info("User action",
        zap.String("email", maskPII(email)),           // Manual masking!
        zap.String("password", "[REDACTED]"),          // Manual redaction!
        zap.Int("user_id", 123),
        zap.String("action", "login"))

    // Verbose structured logging
    logger.Error("Payment failed",
        zap.String("transaction_id", txnID),
        zap.String("card", maskSensitive(cardNumber)), // Manual masking!
        zap.Float64("amount", 99.99),
        zap.String("currency", "USD"),
        zap.Error(err))

    // Performance logging
    logger.Debug("Database operation",
        zap.String("query", query),
        zap.Int("rows", rows),
        zap.Duration("duration", duration))
}
```

### After (Simple & Secure)

```go
package main

import (
    "github.com/cloudresty/emit"
)

func main() {
    // Zero setup required - works out of the box!

    // Automatic security - no manual masking needed
    emit.Info.Field("User action",
        emit.NewFields().
            String("email", email).             // Auto-masked
            String("password", password).       // Auto-masked
            Int("user_id", 123).
            String("action", "login"))

    // Clean structured logging
    emit.Error.Field("Payment failed",
        emit.NewFields().
            String("transaction_id", txnID).
            String("card", cardNumber).         // Auto-masked
            Float64("amount", 99.99).
            String("currency", "USD").
            Error("error", err))

    // High-performance zero-allocation logging
    emit.Debug.ZeroAlloc("Database operation",   // Faster than Zap!
        emit.ZString("query", query),
        emit.ZInt("rows", rows),
        emit.ZDuration("duration", duration))
}
```

### Zap Migration Mapping

| **Zap** | **Emit Equivalent** |
|---------|-------------------|
| `logger.Info("msg")` | `emit.Info.Msg("msg")` |
| `logger.Info("msg", fields...)` | `emit.Info.Field("msg", emit.NewFields()...)` |
| `zap.String("key", value)` | `.String("key", value)` (in Field builder) |
| `zap.Int("key", value)` | `.Int("key", value)` (in Field builder) |
| `zap.Error(err)` | `.Error("error", err)` (in Field builder) |
| `logger.With(fields...).Info()` | Use `emit.NewFields().Clone()` pattern |

## Migration from Zerolog

### Before

```go
package main

import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func main() {
    // Manual security setup required
    log.Info().
        Str("email", maskEmail(email)).          // Manual masking!
        Str("password", "[REDACTED]").           // Manual redaction!
        Int("user_id", 123).
        Str("action", "login").
        Msg("User login")

    // Error logging
    log.Error().
        Err(err).
        Str("service", "payment").
        Float64("amount", 99.99).
        Str("card", maskCard(cardNumber)).       // Manual masking!
        Msg("Payment failed")
}
```

### After

```go
package main

import (
    "github.com/cloudresty/emit"
)

func main() {
    // Automatic security - no setup required
    emit.Info.Field("User login",
        emit.NewFields().
            String("email", email).              // Auto-masked
            String("password", password).        // Auto-masked
            Int("user_id", 123).
            String("action", "login"))

    // Error logging with automatic security
    emit.Error.Field("Payment failed",
        emit.NewFields().
            Error("error", err).
            String("service", "payment").
            Float64("amount", 99.99).
            String("card", cardNumber))          // Auto-masked
}
```

## Advanced Migration Examples

### Microservices Logging Migration

#### Before (Zap + Manual Security)

```go
// Complex setup for each service
func setupLogger(serviceName, version string) *zap.Logger {
    config := zap.NewProductionConfig()
    config.InitialFields = map[string]interface{}{
        "service": serviceName,
        "version": version,
    }
    logger, _ := config.Build()
    return logger
}

func handleUserRequest(logger *zap.Logger, userID int, email string) {
    // Manual security required
    logger.Info("User request",
        zap.Int("user_id", userID),
        zap.String("email", maskPII(email)),     // Manual masking!
        zap.String("endpoint", "/api/users"),
        zap.Time("timestamp", time.Now()))
}
```

#### After (Emit - Simple & Secure)

```go
// Zero setup - configure once globally
func init() {
    emit.SetComponent("user-service")
    emit.SetVersion("v1.2.3")
}

func handleUserRequest(userID int, email string) {
    // Automatic security - no manual work needed
    emit.Info.Field("User request",
        emit.NewFields().
            Int("user_id", userID).
            String("email", email).              // Auto-masked
            String("endpoint", "/api/users").
            Time("timestamp", time.Now()))
}
```

### High-Performance Migration

#### Before (Zap Performance)

```go
// Zap performance logging
func highFrequencyOperation() {
    logger.Debug("Operation",
        zap.String("type", "critical"),
        zap.Int("count", counter),
        zap.Duration("duration", duration))
}

func bulkOperation(items []Item) {
    for _, item := range items {
        logger.Debug("Processing item",
            zap.String("id", item.ID),
            zap.String("type", item.Type))
    }
}
```

#### After (Emit - Faster Performance)

```go
// Emit zero-allocation logging (faster than Zap!)
func highFrequencyOperation() {
    emit.Debug.ZeroAlloc("Operation",            // 174 ns/op vs Zap's 300+
        emit.ZString("type", "critical"),
        emit.ZInt("count", counter),
        emit.ZDuration("duration", duration))
}

func bulkOperation(items []Item) {
    // Memory-pooled for bulk operations
    emit.Debug.Pool("Bulk processing", func(pf *emit.PooledFields) {
        pf.Int("item_count", len(items)).
           Time("started_at", time.Now())
    })

    for _, item := range items {
        emit.Debug.ZeroAlloc("Processing item",   // Ultra-fast per-item
            emit.ZString("id", item.ID),
            emit.ZString("type", item.Type))
    }
}
```

## Migration Checklist

### Phase 1: Setup

- [ ] Add emit to your dependencies: `go get github.com/cloudresty/emit`
- [ ] Set global configuration (component, version, level)
- [ ] Choose production environment variables

### Phase 2: Basic Migration

- [ ] Replace simple log statements with `emit.Info.Msg()`
- [ ] Replace formatted strings with `emit.Info.KeyValue()`
- [ ] Test that sensitive data is automatically masked

### Phase 3: Advanced Migration

- [ ] Convert complex structured logs to `emit.Info.Field()`
- [ ] Optimize hot paths with `emit.Info.ZeroAlloc()`
- [ ] Use `emit.Info.Pool()` for bulk operations

### Phase 4: Security Validation

- [ ] Verify all PII is automatically masked
- [ ] Verify all sensitive data is automatically protected
- [ ] Configure industry-specific protected fields if needed
- [ ] Test in production-like environment

### Phase 5: Performance Optimization

- [ ] Profile your application with emit
- [ ] Optimize hot paths with zero-allocation API
- [ ] Set appropriate log levels for production
- [ ] Monitor performance improvements

## Migration Benefits Summary

### Security Improvements

- ✅ **Automatic PII protection** - No more manual masking
- ✅ **Automatic sensitive data masking** - Zero configuration
- ✅ **Compliance by default** - GDPR, CCPA, HIPAA, PCI DSS
- ✅ **Impossible to accidentally expose data** - Built-in safety

### Performance Improvements

- ✅ **Faster than Zap** - Zero-allocation API outperforms industry leaders
- ✅ **174 ns/op basic logging** - 15-75% faster than alternatives
- ✅ **345 ns/op structured logging** - With automatic security included
- ✅ **Superior memory efficiency** - Fewer allocations, less memory

### Developer Experience Improvements

- ✅ **Clean API** - `emit.Info.Field()` vs cryptic alternatives
- ✅ **Zero setup required** - Works out of the box
- ✅ **Self-documenting code** - Clear intent in method names
- ✅ **IDE-friendly** - Perfect autocomplete and discovery

### Operational Improvements

- ✅ **Kubernetes optimized** - Perfect JSON structure
- ✅ **Environment-aware** - Automatic adaptation
- ✅ **Zero dependencies** - Minimal attack surface
- ✅ **Production ready** - Battle-tested performance

## Get Started Today

1. **Install emit**: `go get github.com/cloudresty/emit`
2. **Start simple**: Replace basic logs with `emit.Info.Msg()`
3. **Add structure**: Use `emit.Info.KeyValue()` for key-value data
4. **Go advanced**: Use `emit.Info.Field()` for complex structured logging
5. **Optimize performance**: Use `emit.Info.ZeroAlloc()` for hot paths

**Experience the future of secure, performant, and elegant logging in Go!**
