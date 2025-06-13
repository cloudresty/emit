# Emit API Reference

Complete reference for all logging methods and field types.

## Elegant API Overview

Every logging level (`Info`, `Error`, `Warn`, `Debug`) provides the same consistent interface:

```go
// All levels support the same methods
emit.Info.Field(msg, fields)           // Structured fields
emit.Info.KeyValue(msg, k, v, ...)     // Key-value pairs
emit.Info.ZeroAlloc(msg, zfields...)   // Ultra-fast zero-allocation
emit.Info.Pool(msg, func)              // Memory-pooled performance
emit.Info.Msg(msg)                     // Simple message
```

## 1. Structured Field Logging

### Basic Usage

```go
// Simple structured logging
emit.Info.Field("User registration",
    emit.NewFields().
        String("username", "john_doe").
        String("email", "user@example.com").
        Int("user_id", 12345).
        Bool("newsletter", true).
        Time("created_at", time.Now()))
```

### Advanced Examples

```go
// Complex business logic logging
emit.Error.Field("Payment processing failed",
    emit.NewFields().
        String("transaction_id", "txn_789").
        String("payment_method", "credit_card").
        String("card_number", "4111-1111-1111-1111").  // Auto-masked
        Float64("amount", 99.99).
        String("currency", "USD").
        String("processor", "stripe").
        String("error_code", "insufficient_funds").
        Int("retry_count", 3).
        Time("failed_at", time.Now()).
        Duration("processing_time", 250*time.Millisecond))

// Microservice communication
emit.Debug.Field("Service call completed",
    emit.NewFields().
        String("from_service", "user-service").
        String("to_service", "auth-service").
        String("method", "ValidateToken").
        String("correlation_id", "corr_123").
        Int("response_code", 200).
        Float64("latency_ms", 45.7).
        Bool("cache_hit", true).
        String("version", "v1.2.3"))

// Database operations
emit.Warn.Field("Query performance warning",
    emit.NewFields().
        String("query", "SELECT * FROM users WHERE active = ?").
        String("database", "postgres").
        Int("rows_returned", 15000).
        Duration("execution_time", 2*time.Second).
        Float64("cpu_usage", 0.85).
        Bool("index_used", false).
        String("optimization_hint", "consider adding index on active column"))
```

### Field Reusability

```go
// Create reusable field builders for microservices
baseServiceFields := emit.NewFields().
    String("service", "auth").
    String("version", "v1.2.3").
    String("environment", "production").
    String("region", "us-west-2")

// Clone and extend for different operations
loginAttempt := baseServiceFields.Clone().
    String("operation", "login").
    String("username", username).
    String("ip_address", clientIP).
    Time("timestamp", time.Now())

passwordReset := baseServiceFields.Clone().
    String("operation", "password_reset").
    String("email", userEmail).      // Auto-masked
    String("reset_token", token).    // Auto-masked
    Time("timestamp", time.Now())

emit.Info.Field("User login attempt", loginAttempt)
emit.Info.Field("Password reset initiated", passwordReset)
```

## 2. Key-Value Pair Logging

### Simple Usage

```go
// Quick key-value logging
emit.Info.KeyValue("User action",
    "user_id", 12345,
    "action", "login",
    "ip", "192.168.1.100",
    "success", true,
    "duration_ms", 250)

// Error context
emit.Error.KeyValue("Database connection failed",
    "host", "db.example.com",
    "port", 5432,
    "database", "users",
    "error", "connection timeout",
    "retry_count", 3)
```

### Complex Examples

```go
// API request logging
emit.Info.KeyValue("HTTP request processed",
    "method", "POST",
    "endpoint", "/api/v1/users",
    "status_code", 201,
    "request_id", "req_abc123",
    "user_agent", "MyApp/1.0",
    "content_length", 1024,
    "response_time_ms", 125,
    "client_ip", "203.0.113.1")

// Business metrics
emit.Warn.KeyValue("High resource usage detected",
    "metric", "memory_usage",
    "current_percent", 87.5,
    "threshold_percent", 80.0,
    "service", "image-processor",
    "instance_id", "i-1234567890abcdef0",
    "auto_scale_triggered", true,
    "alert_sent", true)

// Security events
emit.Error.KeyValue("Suspicious activity detected",
    "event_type", "failed_login_attempts",
    "user_email", "attacker@evil.com",  // Auto-masked
    "ip_address", "198.51.100.1",
    "attempt_count", 15,
    "time_window_minutes", 5,
    "account_locked", true,
    "admin_notified", true)
```

## 3. Zero-Allocation Logging

### Performance-Critical Usage

```go
// Ultra-fast basic logging (174 ns/op)
emit.Info.ZeroAlloc("Request processed")
emit.Error.ZeroAlloc("Connection failed")

// High-frequency structured logging (345 ns/op)
emit.Debug.ZeroAlloc("Database query executed",
    emit.ZString("query_type", "SELECT"),
    emit.ZString("table", "users"),
    emit.ZInt("rows", 1234),
    emit.ZFloat64("duration_ms", 15.7),
    emit.ZBool("cache_hit", false))
```

### Complex Zero-Allocation Examples

```go
// High-throughput API logging
func handleAPIRequest(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    // Process request...

    emit.Info.ZeroAlloc("API request completed",
        emit.ZString("method", r.Method),
        emit.ZString("path", r.URL.Path),
        emit.ZString("user_agent", r.UserAgent()),
        emit.ZInt("status", 200),
        emit.ZInt64("duration_ns", time.Since(start).Nanoseconds()),
        emit.ZBool("authenticated", true),
        emit.ZTime("timestamp", time.Now()))
}

// Real-time metrics collection
func collectMetrics() {
    cpuUsage := getCurrentCPUUsage()
    memUsage := getCurrentMemoryUsage()

    emit.Debug.ZeroAlloc("System metrics",
        emit.ZFloat64("cpu_percent", cpuUsage),
        emit.ZFloat64("memory_percent", memUsage),
        emit.ZInt64("goroutines", int64(runtime.NumGoroutine())),
        emit.ZInt64("heap_bytes", getHeapSize()),
        emit.ZBool("gc_running", isGCRunning()),
        emit.ZTime("collected_at", time.Now()))
}

// Financial transaction logging (ultra-fast for compliance)
func logTransaction(txn Transaction) {
    emit.Info.ZeroAlloc("Transaction processed",
        emit.ZString("transaction_id", txn.ID),
        emit.ZString("account_from", txn.FromAccount),  // Auto-masked
        emit.ZString("account_to", txn.ToAccount),      // Auto-masked
        emit.ZFloat64("amount", txn.Amount),
        emit.ZString("currency", txn.Currency),
        emit.ZString("type", txn.Type),
        emit.ZBool("approved", txn.Approved),
        emit.ZTime("processed_at", txn.ProcessedAt))
}
```

## 4. Memory-Pooled Logging

### High-Throughput Applications

```go
// Memory-efficient logging for high-volume scenarios
emit.Info.Pool("Bulk operation completed", func(pf *emit.PooledFields) {
    pf.String("operation", "bulk_user_import").
       Int("total_records", 50000).
       Int("successful", 49850).
       Int("failed", 150).
       Float64("success_rate", 99.7).
       Duration("total_time", 5*time.Minute).
       Float64("records_per_second", 166.8).
       Bool("validation_enabled", true).
       Time("completed_at", time.Now())
})
```

### Complex Pooled Examples

```go
// Distributed system coordination
emit.Debug.Pool("Distributed lock acquired", func(pf *emit.PooledFields) {
    pf.String("lock_key", "user_update_12345").
       String("node_id", "node_us_west_1").
       String("process_id", os.Getenv("HOSTNAME")).
       Duration("wait_time", 50*time.Millisecond).
       Duration("lease_duration", 30*time.Second).
       Int("retry_count", 2).
       Bool("auto_renewal", true).
       Time("acquired_at", time.Now()).
       String("correlation_id", "dist_op_789")
})

// Machine learning inference logging
emit.Info.Pool("ML model inference completed", func(pf *emit.PooledFields) {
    pf.String("model_name", "user_recommendation_v2").
       String("model_version", "2.1.3").
       String("input_features", "user_history,preferences,context").
       Int("feature_count", 847).
       Float64("confidence_score", 0.923).
       Duration("inference_time", 12*time.Millisecond).
       Int("recommendations_generated", 10).
       Bool("cache_used", false).
       String("gpu_used", "tesla_v100").
       Time("inference_at", time.Now())
})

// Container orchestration events
emit.Warn.Pool("Container resource limit exceeded", func(pf *emit.PooledFields) {
    pf.String("container_id", "cont_abc123def456").
       String("image", "myapp:v1.2.3").
       String("namespace", "production").
       String("pod_name", "myapp-deployment-xyz").
       String("node", "worker-node-05").
       Float64("cpu_limit", 2.0).
       Float64("cpu_usage", 2.3).
       Int64("memory_limit_bytes", 4*1024*1024*1024).
       Int64("memory_usage_bytes", 4.2*1024*1024*1024).
       Bool("oom_killed", false).
       Bool("auto_scaled", true).
       Time("detected_at", time.Now())
})
```

## Field Types Reference

### All Available Types

```go
fields := emit.NewFields().
    String("name", "value").                          // String values
    Int("count", 42).                                // Integer values
    Int64("big_number", 1234567890).                 // 64-bit integers
    Float64("percentage", 95.7).                     // Floating point
    Bool("enabled", true).                           // Boolean values
    Time("timestamp", time.Now()).                   // Time (RFC3339)
    Duration("elapsed", 50*time.Millisecond).        // Duration
    Error("error", fmt.Errorf("something went wrong")). // Error values
    Any("metadata", complexStruct)                   // Any type (JSON)
```

### Zero-Allocation Types

```go
emit.Info.ZeroAlloc("Event occurred",
    emit.ZString("service", "auth"),
    emit.ZInt("count", 100),
    emit.ZInt64("bytes", 1048576),
    emit.ZFloat64("ratio", 0.75),
    emit.ZBool("success", true),
    emit.ZTime("when", time.Now()),
    emit.ZDuration("took", 25*time.Millisecond))
```

### Pooled Field Types

```go
emit.Info.Pool("Operation", func(pf *emit.PooledFields) {
    pf.String("key", "value").
       Int("number", 123).
       Int64("big", 9876543210).
       Float64("decimal", 3.14159).
       Bool("flag", false).
       Time("timestamp", time.Now()).
       Duration("elapsed", time.Second).
       Error("err", someError)
})
```

## Performance Guidelines

### When to Use Each API

| **API Type** | **Use Case** | **Performance** | **Memory** |
|--------------|--------------|-----------------|------------|
| **Field()** | Complex structured logs | Good | Moderate |
| **KeyValue()** | Simple key-value pairs | Very Good | Low |
| **ZeroAlloc()** | High-frequency logging | Excellent | Very Low |
| **Pool()** | Memory-sensitive bulk ops | Excellent | Minimal |
| **Msg()** | Simple messages | Excellent | Minimal |

### Performance Tips

1. **Use ZeroAlloc()** for hot paths and high-frequency logging
2. **Use Pool()** for bulk operations and memory-sensitive scenarios
3. **Use Field()** for complex business logic logging
4. **Use KeyValue()** for simple, readable logging
5. **Use Msg()** for basic status messages

### Benchmark Results

```plaintext
BenchmarkField-10           1,000,000   1,112 ns/op   1,201 B/op    13 allocs/op
BenchmarkKeyValue-10        2,575,186     469 ns/op     464 B/op     5 allocs/op
BenchmarkZeroAlloc-10       6,895,036     174 ns/op      32 B/op     1 allocs/op
BenchmarkPool-10            3,414,537     345 ns/op     464 B/op     6 allocs/op
BenchmarkMsg-10            10,000,000     150 ns/op      24 B/op     1 allocs/op
```

## Migration Examples

### From Other Loggers

```go
// From standard log
log.Printf("User %s logged in", username)
// ↓ Becomes (secure)
emit.Info.KeyValue("User logged in", "username", username)  // Auto-masked

// From logrus
logrus.WithFields(logrus.Fields{"user": id}).Info("Action")
// ↓ Becomes
emit.Info.Field("Action", emit.NewFields().Int("user", id))

// From zap
logger.Info("Event", zap.String("key", value))
// ↓ Becomes
emit.Info.KeyValue("Event", "key", value)  // Auto-secured
```

This API reference provides comprehensive examples for all logging patterns in emit.
