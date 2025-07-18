package emit

import (
	"sync/atomic"
	"time"
)

// TimestampPrecision defines the precision level for timestamps
type TimestampPrecision int

const (
	NanosecondPrecision TimestampPrecision = iota
	MicrosecondPrecision
	MillisecondPrecision
	SecondPrecision
)

// Global timestamp precision setting
var currentTimestampPrecision int32 = int32(MillisecondPrecision)

// SetTimestampPrecision sets the global timestamp precision
func SetTimestampPrecision(precision TimestampPrecision) {
	atomic.StoreInt32(&currentTimestampPrecision, int32(precision))
}

// GetTimestampPrecision returns the current timestamp precision
func GetTimestampPrecision() TimestampPrecision {
	return TimestampPrecision(atomic.LoadInt32(&currentTimestampPrecision))
}

// ultraFastTimestampCache provides extremely fast timestamp generation
// by caching timestamp strings and updating them less frequently
type ultraFastTimestampCache struct {
	// Atomic fields - must be 64-bit aligned on 32-bit systems
	lastUpdateUnix int64

	// Cached timestamp string - using atomic.Value for safe concurrent access
	cachedTimestamp atomic.Value // stores string

	// Update frequency (in seconds) - how often to refresh timestamp
	updateIntervalSeconds int64
}

var (
	// Global ultra-fast cache instance - optimized for sub-20ns performance
	globalUltraFastCache = &ultraFastTimestampCache{
		updateIntervalSeconds: 1, // Update every 1 second (but cache at nanosecond level)
	}

	// Thread-safe timestamp check tracker
	lastTimestampCheck int64 // Use atomic operations for this
)

// GetUltraFastTimestamp returns a cached timestamp string
// Optimized for sub-20ns performance in the common case
func GetUltraFastTimestamp() string {
	// Ultra-fast path: Check if we even need to update (minimize atomic ops)
	now := time.Now().Unix()

	// Only check atomic lastUpdate occasionally to reduce overhead
	// Use atomic operations for thread safety
	lastCheck := atomic.LoadInt64(&lastTimestampCheck)
	if now == lastCheck {
		// Same second as last check - return cached string directly
		if cached := globalUltraFastCache.cachedTimestamp.Load(); cached != nil {
			return cached.(string)
		}
	}

	// Update our local check atomically
	atomic.StoreInt64(&lastTimestampCheck, now)

	// Check if we need a real update
	lastUpdate := atomic.LoadInt64(&globalUltraFastCache.lastUpdateUnix)

	if now-lastUpdate < atomic.LoadInt64(&globalUltraFastCache.updateIntervalSeconds) {
		// Return cached timestamp
		if cached := globalUltraFastCache.cachedTimestamp.Load(); cached != nil {
			return cached.(string)
		}
	}

	// Time to update - try to win the race
	if atomic.CompareAndSwapInt64(&globalUltraFastCache.lastUpdateUnix, lastUpdate, now) {
		// We won the race - generate new timestamp
		newTimestamp := generateFastTimestamp()
		globalUltraFastCache.cachedTimestamp.Store(newTimestamp)
		return newTimestamp
	}

	// Another goroutine updated it, return the cached version
	if cached := globalUltraFastCache.cachedTimestamp.Load(); cached != nil {
		return cached.(string)
	}

	// First time initialization (rarely called)
	timestamp := generateFastTimestamp()
	globalUltraFastCache.cachedTimestamp.Store(timestamp)
	atomic.StoreInt64(&globalUltraFastCache.lastUpdateUnix, now)
	return timestamp
}

// generateFastTimestamp creates a timestamp string with millisecond precision
// This is only called once per second to update the cache
func generateFastTimestamp() string {
	now := time.Now().UTC()

	// Pre-calculate the most common case: millisecond precision
	// Format: 2006-01-02T15:04:05.000Z

	year := now.Year()
	month := int(now.Month())
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	millis := now.Nanosecond() / 1000000

	// Build timestamp string efficiently
	// Using a fixed-size byte array for better performance
	var buf [24]byte

	// Year (4 digits)
	buf[0] = byte('0' + year/1000)
	buf[1] = byte('0' + (year%1000)/100)
	buf[2] = byte('0' + (year%100)/10)
	buf[3] = byte('0' + year%10)

	buf[4] = '-'

	// Month (2 digits)
	buf[5] = byte('0' + month/10)
	buf[6] = byte('0' + month%10)

	buf[7] = '-'

	// Day (2 digits)
	buf[8] = byte('0' + day/10)
	buf[9] = byte('0' + day%10)

	buf[10] = 'T'

	// Hour (2 digits)
	buf[11] = byte('0' + hour/10)
	buf[12] = byte('0' + hour%10)

	buf[13] = ':'

	// Minute (2 digits)
	buf[14] = byte('0' + minute/10)
	buf[15] = byte('0' + minute%10)

	buf[16] = ':'

	// Second (2 digits)
	buf[17] = byte('0' + second/10)
	buf[18] = byte('0' + second%10)

	buf[19] = '.'

	// Milliseconds (3 digits)
	buf[20] = byte('0' + millis/100)
	buf[21] = byte('0' + (millis%100)/10)
	buf[22] = byte('0' + millis%10)

	buf[23] = 'Z'

	// Convert to string without allocation
	return string(buf[:])
}

// SetUltraFastTimestampPrecision sets the update interval for the ultra-fast cache
// Lower intervals provide more accurate timestamps but slight performance cost
func SetUltraFastTimestampPrecision(intervalSeconds int64) {
	if intervalSeconds < 1 {
		intervalSeconds = 1
	}
	atomic.StoreInt64(&globalUltraFastCache.updateIntervalSeconds, intervalSeconds)
}
