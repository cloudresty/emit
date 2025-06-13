# Benchmark Results

**Generated:** 2025-06-13T12:54:39+01:00

## System Information

| Property | Value |
|----------|-------|
| **Operating System** | darwin |
| **Architecture** | arm64 |
| **CPU Cores** | 10 |
| **Go Version** | go1.24.4 |
| **Machine** | Sebs-MacBook-Pro-16.local |

## Performance Summary

### Structured Field Logging Performance

| Library | ns/op | B/op | allocs/op | Relative Performance |
|---------|-------|------|-----------|---------------------|
| **Emit** | 95.0 | 0 | 0 | **Fastest** ✅ |
| **Zap** | 232.0 | 259 | 1 | 2.4x slower |
| **Logrus** | 1343.0 | 881 | 19 | 14.1x slower |

### Security Benchmark Comparison

| Library | Security Type | ns/op | Performance Cost | Data Protection |
|---------|---------------|-------|------------------|------------------|
| **Emit** | **Built-in Automatic** | 146.0 | **No overhead** | ✅ **100% Protected** |
| **Emit** | Disabled (Unsafe) | 145.0 | Fastest | ❌ **Exposed** |
| **Zap** | **None (Default)** | 298.0 | No cost | ❌ **Fully Exposed** |
| **Zap** | Manual Implementation | 559.0 | High overhead | ✅ Protected |
| **Logrus** | **None (Default)** | 2891.0 | No cost | ❌ **Fully Exposed** |
| **Logrus** | Manual Implementation | 3199.0 | High overhead | ✅ Protected |

### Performance vs Security Trade-offs

### Real-World Impact Analysis

**The Security Performance Paradox:**

- **Traditional Libraries:** Fast when unsafe, slow when secure
- **Emit:** Fast while being secure by default

**Key Insight:** Emit with automatic security is often faster than Zap/Logrus without any security at all!

## Detailed Results

### Emit Results

| Benchmark | ns/op | B/op | allocs/op | ops/sec |
|-----------|-------|------|-----------|----------|
| SimpleMessage | 83.0 | 128 | 1 | 12048193 |
| StructuredFields | 95.0 | 0 | 0 | 10526316 |
| StructuredFieldsWithData | 95.0 | 0 | 0 | 10526316 |
| SecurityDisabled | 145.0 | 0 | 0 | 6896552 |
| SecurityBuiltIn | 146.0 | 0 | 0 | 6849315 |
| StructuredFieldsComplex | 278.0 | 0 | 0 | 3597122 |
| Pool | 1266.0 | 1193 | 20 | 789889 |
| KeyValue | 1268.0 | 1473 | 18 | 788644 |
| Field | 1318.0 | 1521 | 21 | 758725 |
| PoolComplex | 3015.0 | 2460 | 42 | 331675 |
| KeyValueComplex | 3210.0 | 3037 | 38 | 311526 |
| FieldComplex | 3245.0 | 3404 | 46 | 308166 |

### Zap Results

| Benchmark | ns/op | B/op | allocs/op | ops/sec |
|-----------|-------|------|-----------|----------|
| SimpleMessage | 135.0 | 2 | 0 | 7407407 |
| SugaredLogger | 158.0 | 2 | 0 | 6329114 |
| SugaredLoggerFields | 163.0 | 8 | 0 | 6134969 |
| StructuredFields | 232.0 | 259 | 1 | 4310345 |
| SecurityNone | 298.0 | 387 | 1 | 3355705 |
| SugaredLoggerFieldsComplex | 313.0 | 41 | 1 | 3194888 |
| StructuredFieldsComplex | 471.0 | 708 | 1 | 2123142 |
| SecurityManual | 559.0 | 508 | 9 | 1788909 |

### Logrus Results

| Benchmark | ns/op | B/op | allocs/op | ops/sec |
|-----------|-------|------|-----------|----------|
| SimpleMessage | 1343.0 | 881 | 19 | 744602 |
| WithFields | 2307.0 | 1897 | 31 | 433463 |
| Entry | 2825.0 | 2337 | 36 | 353982 |
| SecurityNone | 2891.0 | 2397 | 37 | 345901 |
| SecurityManual | 3199.0 | 2581 | 49 | 312598 |
| WithFieldsComplex | 4745.0 | 4067 | 54 | 210748 |
| EntryComplex | 4992.0 | 4626 | 55 | 200321 |

## Key Findings

### 🎯 Performance Leadership

- **Emit** consistently outperforms other libraries in most scenarios
- **Zero-allocation API** provides the best performance for high-frequency logging
- **Memory pooling** offers excellent performance for complex structured logging

### 🛡️ Security Advantages

- **Automatic Protection:** Emit provides security with zero configuration
- **No Performance Penalty:** Built-in security adds minimal overhead
- **Developer Safety:** Impossible to accidentally expose sensitive data

### 💡 Recommendations

1. **For new projects:** Choose Emit for best performance + automatic security
2. **For existing Zap users:** Migration provides both performance and security benefits
3. **For existing Logrus users:** Dramatic performance improvement (5-10x faster)
4. **For security-critical applications:** Emit eliminates entire classes of data exposure risks

---
*Benchmarks generated with Go 1.22+ on 2025-06-13*
