package emit

import (
	"regexp"
	"strings"
	"sync"
)

// Optimized security implementation with caching and pre-compilation

// Field pattern cache for faster lookup
type fieldPatternCache struct {
	mu             sync.RWMutex
	piiCache       map[string]bool
	sensitiveCache map[string]bool
}

var (
	// Pre-compiled regex patterns for better performance
	emailRegex      = regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	phoneRegex      = regexp.MustCompile(`\b(?:\+?1[-.\s]?)?\(?[0-9]{3}\)?[-.\s]?[0-9]{3}[-.\s]?[0-9]{4}\b`)
	creditCardRegex = regexp.MustCompile(`\b(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|3[47][0-9]{13}|3[0-9]{13}|6(?:011|5[0-9]{2})[0-9]{12})\b`)

	// Global field cache for faster lookups
	fieldCache = &fieldPatternCache{
		piiCache:       make(map[string]bool, 100),
		sensitiveCache: make(map[string]bool, 100),
	}

	// Pre-built lookup maps for O(1) field checking
	piiFieldsMap       map[string]bool
	sensitiveFieldsMap map[string]bool
	onceInit           sync.Once
)

// initializeFieldMaps builds lookup maps for O(1) field pattern matching
func initializeFieldMaps() {
	onceInit.Do(func() {
		// Build PII fields map
		piiFieldsMap = make(map[string]bool, len(defaultPIIFields)*2)
		for _, pattern := range defaultPIIFields {
			piiFieldsMap[pattern] = true
			piiFieldsMap[strings.ToUpper(pattern)] = true // Add uppercase variant
		}

		// Build sensitive fields map
		sensitiveFieldsMap = make(map[string]bool, len(defaultSensitiveFields)*2)
		for _, pattern := range defaultSensitiveFields {
			sensitiveFieldsMap[pattern] = true
			sensitiveFieldsMap[strings.ToUpper(pattern)] = true // Add uppercase variant
		}
	})
}

// Fast PII field checking with caching
func (l *Logger) isPIIFieldFast(fieldName string) bool {
	if l.piiMode == SHOW_PII {
		return false
	}

	initializeFieldMaps()

	// Check cache first
	fieldCache.mu.RLock()
	if cached, exists := fieldCache.piiCache[fieldName]; exists {
		fieldCache.mu.RUnlock()
		return cached
	}
	fieldCache.mu.RUnlock()

	// Fast lookup in pre-built map
	lowerFieldName := strings.ToLower(fieldName)
	isPII := piiFieldsMap[lowerFieldName]

	if !isPII {
		// Fallback to substring search only if direct lookup fails
		for pattern := range piiFieldsMap {
			if strings.Contains(lowerFieldName, pattern) {
				isPII = true
				break
			}
		}
	}

	// Cache the result
	fieldCache.mu.Lock()
	fieldCache.piiCache[fieldName] = isPII
	fieldCache.mu.Unlock()

	return isPII
}

// Fast sensitive field checking with caching
func (l *Logger) isSensitiveFieldFast(fieldName string) bool {
	if l.sensitiveMode == SHOW_SENSITIVE {
		return false
	}

	initializeFieldMaps()

	// Check cache first
	fieldCache.mu.RLock()
	if cached, exists := fieldCache.sensitiveCache[fieldName]; exists {
		fieldCache.mu.RUnlock()
		return cached
	}
	fieldCache.mu.RUnlock()

	// Fast lookup in pre-built map
	lowerFieldName := strings.ToLower(fieldName)
	isSensitive := sensitiveFieldsMap[lowerFieldName]

	if !isSensitive {
		// Fallback to substring search only if direct lookup fails
		for pattern := range sensitiveFieldsMap {
			if strings.Contains(lowerFieldName, pattern) {
				isSensitive = true
				break
			}
		}
	}

	// Cache the result
	fieldCache.mu.Lock()
	fieldCache.sensitiveCache[fieldName] = isSensitive
	fieldCache.mu.Unlock()

	return isSensitive
}

// Optimized field masking with pre-allocated map and minimal allocations
func (l *Logger) maskSensitiveFieldsFast(fields map[string]any) map[string]any {
	if (l.sensitiveMode == SHOW_SENSITIVE && l.piiMode == SHOW_PII) || len(fields) == 0 {
		return fields
	}

	// Pre-allocate with exact capacity to avoid map growth
	maskedFields := make(map[string]any, len(fields))

	for key, value := range fields {
		// Fast path: check PII first (more specific), then sensitive data
		if l.isPIIFieldFast(key) {
			maskedFields[key] = l.piiMaskString
		} else if l.isSensitiveFieldFast(key) {
			maskedFields[key] = l.maskString
		} else {
			// Handle nested maps recursively
			if nestedMap, ok := value.(map[string]any); ok {
				maskedFields[key] = l.maskSensitiveFieldsFast(nestedMap)
			} else {
				maskedFields[key] = value
			}
		}
	}

	return maskedFields
}

// Content-based masking for string values (optional advanced feature)
func (l *Logger) maskSensitiveContent(value string) string {
	// Mask email patterns
	value = emailRegex.ReplaceAllString(value, "***EMAIL***")

	// Mask phone patterns
	value = phoneRegex.ReplaceAllString(value, "***PHONE***")

	// Mask credit card patterns
	value = creditCardRegex.ReplaceAllString(value, "***CARD***")

	return value
}

// ClearFieldCache clears the field pattern cache (for testing or dynamic field updates)
func ClearFieldCache() {
	fieldCache.mu.Lock()
	defer fieldCache.mu.Unlock()

	fieldCache.piiCache = make(map[string]bool, 100)
	fieldCache.sensitiveCache = make(map[string]bool, 100)
}
