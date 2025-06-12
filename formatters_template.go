package emit

import (
	"sync/atomic"
	"time"
)

var (
	// Pre-formatted JSON templates for hot path
	jsonTemplatePrefix = []byte(`{"timestamp":"`)
	jsonLevelDebug     = []byte(`","level":"debug","msg":"`)
	jsonLevelInfo      = []byte(`","level":"info","msg":"`)
	jsonLevelWarn      = []byte(`","level":"warn","msg":"`)
	jsonLevelError     = []byte(`","level":"error","msg":"`)
	jsonSuffix         = []byte(`"}` + "\n")

	// Cached timestamp (1-second precision)
	cachedTimestamp atomic.Value
)

func init() {
	updateCachedTimestamp()
	go cachedTimestampUpdater()
}

func cachedTimestampUpdater() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		updateCachedTimestamp()
	}
}

func updateCachedTimestamp() {
	now := time.Now().UTC()
	// Pre-formatted timestamp: 2006-01-02T15:04:05Z
	timestamp := []byte(now.Format("2006-01-02T15:04:05Z"))
	cachedTimestamp.Store(timestamp)
}

func getCachedTimestamp() []byte {
	return cachedTimestamp.Load().([]byte)
}

// Template-based logging functionality
// This file now focuses on template management and cached components
// ZeroAlloc logging has been moved to formatters_zero_alloc.go for proper organization
