package emit

// This file is reserved for the formatters_zero_alloc.go implementation.
// Currently, the ZeroAlloc API routes through formatters_template.go (logZeroBlazing).
// This file exists to maintain the clean API-to-file mapping but contains no active code.
// The previous implementation was an unused optimization experiment.

// Zero-allocation logging is implemented in:
// emit.Info.ZeroAlloc() → logWithZeroAlloc() → defaultLogger.logZero() → logZeroBlazing() (formatters_template.go)
