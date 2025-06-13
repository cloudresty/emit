package main

import (
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

// LogrusBenchmarkSet contains all Logrus-specific benchmarks
type LogrusBenchmarkSet struct{}

// GetBenchmarks returns all Logrus benchmark functions
func (l LogrusBenchmarkSet) GetBenchmarks() []BenchmarkFunc {
	return []BenchmarkFunc{
		// Simple message benchmarks
		{"Logrus_SimpleMessage", l.BenchmarkSimpleMessage},

		// Field benchmarks
		{"Logrus_WithFields", l.BenchmarkWithFields},
		{"Logrus_WithFieldsComplex", l.BenchmarkWithFieldsComplex},

		// Entry-based benchmarks
		{"Logrus_Entry", l.BenchmarkEntry},
		{"Logrus_EntryComplex", l.BenchmarkEntryComplex},

		// Security benchmarks
		{"Logrus_SecurityNone", l.BenchmarkSecurityNone},
		{"Logrus_SecurityManual", l.BenchmarkSecurityManual},
	}
}

// Simple message logging benchmarks
func (l LogrusBenchmarkSet) BenchmarkSimpleMessage(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		logrusLogger.Info("Simple log message")
	}
}

// WithFields benchmarks
func (l LogrusBenchmarkSet) BenchmarkWithFields(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		logrusLogger.WithFields(logrus.Fields{
			"user_id":    "12345",
			"action":     "login",
			"ip_address": "192.168.1.100",
			"success":    true,
		}).Info("User action")
	}
}

func (l LogrusBenchmarkSet) BenchmarkWithFieldsComplex(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		logrusLogger.WithFields(logrus.Fields{
			"service":        "user-service",
			"operation":      "create_user",
			"user_id":        "12345",
			"email":          "user@example.com",
			"ip_address":     "192.168.1.100",
			"status_code":    201,
			"duration_ms":    15.75,
			"success":        true,
			"timestamp":      time.Now(),
			"correlation_id": "corr_abc123",
		}).Info("Complex operation")
	}
}

// Entry-based benchmarks (reusable logger)
func (l LogrusBenchmarkSet) BenchmarkEntry(b *testing.B) {
	entry := logrusLogger.WithFields(logrus.Fields{
		"service": "user-service",
		"version": "v1.2.3",
	})

	b.ResetTimer()
	for b.Loop() {
		entry.WithFields(logrus.Fields{
			"user_id":    "12345",
			"action":     "login",
			"ip_address": "192.168.1.100",
			"success":    true,
		}).Info("User action")
	}
}

func (l LogrusBenchmarkSet) BenchmarkEntryComplex(b *testing.B) {
	entry := logrusLogger.WithFields(logrus.Fields{
		"service": "user-service",
		"version": "v1.2.3",
	})

	b.ResetTimer()
	for b.Loop() {
		entry.WithFields(logrus.Fields{
			"operation":      "create_user",
			"user_id":        "12345",
			"email":          "user@example.com",
			"ip_address":     "192.168.1.100",
			"status_code":    201,
			"duration_ms":    15.75,
			"success":        true,
			"timestamp":      time.Now(),
			"correlation_id": "corr_abc123",
		}).Info("Complex operation")
	}
}

// Security benchmarks
func (l LogrusBenchmarkSet) BenchmarkSecurityNone(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		logrusLogger.WithFields(logrus.Fields{
			"password": "super_secret_123", // EXPOSED!
			"email":    "user@example.com", // EXPOSED!
			"api_key":  "sk_live_abc123",   // EXPOSED!
			"user_id":  "12345",
			"ssn":      "123-45-6789", // EXPOSED!
			"action":   "register",
		}).Info("User registration")
	}
}

func (l LogrusBenchmarkSet) BenchmarkSecurityManual(b *testing.B) {
	maskSensitive := func(s string) string {
		if len(s) <= 4 {
			return "***"
		}
		return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
	}

	b.ResetTimer()
	for b.Loop() {
		logrusLogger.WithFields(logrus.Fields{
			"password": maskSensitive("super_secret_123"), // Manually masked
			"email":    maskSensitive("user@example.com"), // Manually masked
			"api_key":  maskSensitive("sk_live_abc123"),   // Manually masked
			"user_id":  "12345",
			"ssn":      maskSensitive("123-45-6789"), // Manually masked
			"action":   "register",
		}).Info("User registration")
	}
}
