# Security & Compliance Guide

Comprehensive guide to emit's built-in security features and compliance capabilities.

## Automatic Data Protection

Emit automatically protects sensitive information without any configuration required.

### PII (Personally Identifiable Information) Protection

Emit automatically detects and masks PII data:

```go
emit.Info.Field("User registration",
    emit.NewFields().
        String("email", "user@example.com").           // → "***PII***"
        String("phone", "+1-555-123-4567").            // → "***PII***"
        String("full_name", "John Doe").               // → "***PII***"
        String("address", "123 Main St, Anytown").     // → "***PII***"
        String("ssn", "123-45-6789").                  // → "***PII***"
        String("credit_card", "4111-1111-1111-1111").  // → "***PII***"
        String("ip_address", "192.168.1.100").         // → "***PII***"
        Int("user_id", 12345))                         // → 12345 (safe)
```

### Sensitive Data Protection

Automatically masks authentication and security-related data:

```go
emit.Error.Field("Authentication failed",
    emit.NewFields().
        String("api_key", "sk-1234567890abcdef").       // → "***MASKED***"
        String("password", "user_password").            // → "***MASKED***"
        String("access_token", "bearer_token_xyz").     // → "***MASKED***"
        String("private_key", "-----BEGIN PRIVATE").    // → "***MASKED***"
        String("session_id", "sess_abc123").            // → "***MASKED***"
        String("jwt_token", "eyJ0eXAiOiJKV1Q").         // → "***MASKED***"
        String("username", "john_doe").                 // → "john_doe" (safe)
        Int("attempt_count", 3))                        // → 3 (safe)
```

## Protected Field Categories

### Automatically Protected PII Fields

| **Field Pattern** | **Example** | **Masked As** |
|------------------|-------------|---------------|
| Email addresses | `user@example.com` | `***PII***` |
| Phone numbers | `+1-555-123-4567` | `***PII***` |
| Names | `John Doe` | `***PII***` |
| Addresses | `123 Main St` | `***PII***` |
| IP addresses | `192.168.1.100` | `***PII***` |
| Credit cards | `4111-1111-1111-1111` | `***PII***` |
| SSN/Tax IDs | `123-45-6789` | `***PII***` |
| Driver licenses | `DL123456789` | `***PII***` |
| Passport numbers | `P1234567` | `***PII***` |
| Date of birth | `1990-01-01` | `***PII***` |

### Automatically Protected Sensitive Fields

| **Field Pattern** | **Example** | **Masked As** |
|------------------|-------------|---------------|
| Passwords | `my_password` | `***MASKED***` |
| API keys | `sk-1234567890abcdef` | `***MASKED***` |
| Access tokens | `bearer_token_xyz` | `***MASKED***` |
| Private keys | `-----BEGIN PRIVATE` | `***MASKED***` |
| Session IDs | `sess_abc123def456` | `***MASKED***` |
| JWT tokens | `eyJ0eXAiOiJKV1Q` | `***MASKED***` |
| OAuth tokens | `oauth_token_xyz` | `***MASKED***` |
| Certificates | `-----BEGIN CERT` | `***MASKED***` |
| Database URLs | `postgres://user:pass@` | `***MASKED***` |
| Secrets | `secret_value_123` | `***MASKED***` |

## Compliance Features

### GDPR Compliance

Emit helps ensure GDPR compliance by automatically protecting EU personal data:

```go
// GDPR Article 4 - Personal Data Protection
emit.Info.Field("GDPR compliant user event",
    emit.NewFields().
        String("user_email", "eu_user@company.com").    // Auto-masked
        String("user_name", "Hans Mueller").            // Auto-masked
        String("user_address", "Berlin, Germany").      // Auto-masked
        String("user_phone", "+49-30-12345678").        // Auto-masked
        String("user_ip", "85.214.132.117").            // Auto-masked
        String("event_type", "profile_update").         // Safe
        Time("timestamp", time.Now()).                  // Safe
        Bool("consent_given", true))                    // Safe

// Output: All PII automatically masked as "***PII***"
```

### CCPA Compliance

California Consumer Privacy Act compliance for US companies:

```go
// CCPA - California Personal Information Protection
emit.Info.Field("CCPA compliant data processing",
    emit.NewFields().
        String("california_resident_email", "ca_user@email.com"). // Auto-masked
        String("personal_identifier", "CA_ID_123456").            // Auto-masked
        String("biometric_data", "fingerprint_hash").             // Auto-masked
        String("location_data", "San Francisco, CA").             // Auto-masked
        String("purchase_history", "electronics, books").         // Business data
        Float64("transaction_amount", 299.99).                    // Business data
        Bool("opt_out_requested", false))                         // Safe

// Personal information automatically protected
```

### HIPAA Compliance (Healthcare)

For healthcare applications, configure additional protected health information:

```go
// Configure HIPAA-specific protected fields
emit.AddPIIField("patient_id")
emit.AddPIIField("medical_record_number")
emit.AddPIIField("diagnosis")
emit.AddPIIField("treatment_plan")
emit.AddPIIField("prescription")

// Healthcare logging with automatic PHI protection
emit.Info.Field("Patient treatment logged",
    emit.NewFields().
        String("patient_id", "P123456").                    // Auto-masked
        String("medical_record_number", "MR789012").        // Auto-masked
        String("diagnosis", "Type 2 Diabetes").             // Auto-masked
        String("treatment_plan", "Insulin therapy").        // Auto-masked
        String("prescription", "Metformin 500mg").          // Auto-masked
        String("healthcare_provider", "Dr. Smith").         // Auto-masked
        String("facility", "General Hospital").             // Safe
        Time("treatment_date", time.Now()).                 // Safe
        Float64("dosage_mg", 500.0))                        // Safe

// All PHI automatically protected
```

### PCI DSS Compliance (Payment Cards)

For payment processing, emit automatically protects cardholder data:

```go
// PCI DSS - Payment Card Industry Data Security Standard
emit.Info.Field("Payment processed securely",
    emit.NewFields().
        String("card_number", "4111-1111-1111-1111").       // Auto-masked
        String("cardholder_name", "John Doe").              // Auto-masked
        String("expiry_date", "12/25").                     // Auto-masked
        String("cvv", "123").                               // Auto-masked
        String("billing_address", "123 Main St").           // Auto-masked
        String("transaction_id", "TXN_ABC123").             // Safe
        Float64("amount", 99.99).                           // Safe
        String("currency", "USD").                          // Safe
        String("merchant_id", "MERCHANT_789").              // Safe
        Bool("approved", true))                             // Safe

// All cardholder data automatically protected
```

## Custom Security Configuration

### Environment-Based Security

```bash
# Production (maximum security)
export EMIT_MASK_SENSITIVE=true
export EMIT_MASK_PII=true
export EMIT_MASK_STRING="[REDACTED]"
export EMIT_PII_MASK_STRING="[PERSONAL_DATA]"

# Development (show data for debugging)
export EMIT_MASK_SENSITIVE=false
export EMIT_MASK_PII=false

# Custom compliance (healthcare example)
export EMIT_MASK_STRING="[PHI_PROTECTED]"
export EMIT_PII_MASK_STRING="[PATIENT_DATA]"
```

### Programmatic Security Configuration

```go
// Custom sensitive field patterns
emit.AddSensitiveField("internal_token")
emit.AddSensitiveField("company_secret")
emit.AddSensitiveField("encryption_key")

// Custom PII field patterns
emit.AddPIIField("employee_id")
emit.AddPIIField("patient_record")
emit.AddPIIField("customer_account")

// Industry-specific field sets
healthcareFields := []string{
    "patient_id", "medical_record", "diagnosis",
    "treatment", "prescription", "doctor_name",
}
emit.SetPIIFields(healthcareFields)

financialSensitive := []string{
    "account_number", "routing_number", "swift_code",
    "transaction_key", "settlement_id", "wire_reference",
}
emit.SetSensitiveFields(financialSensitive)

// Custom mask strings for different data types
emit.SetMaskString("[CLASSIFIED]")          // For sensitive data
emit.SetPIIMaskString("[PERSONAL_INFO]")    // For PII data
```

## Industry-Specific Examples

### Financial Services

```go
// Banking transaction logging
emit.Info.Field("Wire transfer initiated",
    emit.NewFields().
        String("from_account", "ACC_123456789").        // Auto-masked PII
        String("to_account", "ACC_987654321").          // Auto-masked PII
        String("routing_number", "021000021").          // Auto-masked PII
        String("swift_code", "CHASUS33").               // Auto-masked sensitive
        Float64("amount", 50000.00).                    // Business data
        String("currency", "USD").                      // Business data
        String("reference", "WIRE_REF_ABC123").         // Business data
        String("initiator_ip", "203.0.113.45").         // Auto-masked PII
        Time("initiated_at", time.Now()).               // Safe
        Bool("compliance_checked", true))               // Safe

// Anti-money laundering (AML) logging
emit.Warn.Field("Suspicious transaction flagged",
    emit.NewFields().
        String("customer_id", "CUST_789012").           // Auto-masked PII
        String("customer_name", "Jane Smith").          // Auto-masked PII
        String("customer_ssn", "987-65-4321").          // Auto-masked PII
        Float64("transaction_amount", 15000.00).        // Business data
        String("transaction_type", "cash_deposit").     // Business data
        String("flagged_reason", "unusual_pattern").    // Business data
        Int("risk_score", 85).                          // Business data
        Bool("manual_review_required", true).           // Safe
        Time("flagged_at", time.Now()))                 // Safe
```

### Healthcare Systems

```go
// Patient visit logging
emit.Info.Field("Patient appointment completed",
    emit.NewFields().
        String("patient_id", "PAT_456789").            // Auto-masked PII
        String("patient_name", "Bob Johnson").         // Auto-masked PII
        String("patient_dob", "1975-03-15").           // Auto-masked PII
        String("patient_ssn", "456-78-9012").          // Auto-masked PII
        String("diagnosis_code", "E11.9").             // Auto-masked PII
        String("diagnosis", "Type 2 diabetes").        // Auto-masked PII
        String("treatment", "Lifestyle counseling").   // Auto-masked PII
        String("physician", "Dr. Emily Chen").         // Auto-masked PII
        String("facility", "Metro Health Center").     // Business data
        String("appointment_type", "follow_up").       // Business data
        Duration("visit_duration", 30*time.Minute).    // Business data
        Time("visit_date", time.Now()))                // Business data

// Medical device data logging
emit.Debug.Field("Medical device reading",
    emit.NewFields().
        String("device_id", "DEV_GLUC_001").           // Business data
        String("patient_id", "PAT_123456").            // Auto-masked PII
        Float64("glucose_level", 125.5).               // Medical data (safe for analytics)
        String("measurement_unit", "mg/dL").           // Business data
        Time("reading_time", time.Now()).              // Business data
        Bool("within_normal_range", true).             // Business data
        String("device_location", "home").             // Business data
        Bool("alert_triggered", false))                // Business data
```

### E-commerce Platforms

```go
// Customer order processing
emit.Info.Field("Order placed successfully",
    emit.NewFields().
        String("customer_email", "shopper@email.com").  // Auto-masked PII
        String("customer_name", "Alice Wilson").        // Auto-masked PII
        String("customer_phone", "+1-555-987-6543").    // Auto-masked PII
        String("shipping_address", "456 Oak Ave").      // Auto-masked PII
        String("payment_method", "credit_card").        // Business data
        String("card_last_four", "1234").               // Auto-masked PII
        String("order_id", "ORD_789123").               // Business data
        Float64("order_total", 127.99).                 // Business data
        String("currency", "USD").                      // Business data
        Int("item_count", 3).                           // Business data
        Time("order_time", time.Now()).                 // Business data
        Bool("express_shipping", true))                 // Business data

// Customer support interaction
emit.Info.Field("Support ticket resolved",
    emit.NewFields().
        String("customer_email", "help@customer.com"). // Auto-masked PII
        String("ticket_id", "TICK_456789").            // Business data
        String("issue_category", "billing").           // Business data
        String("resolution", "refund_processed").      // Business data
        String("agent_id", "AGT_123").                 // Business data
        Duration("resolution_time", 2*time.Hour).      // Business data
        Int("satisfaction_score", 5).                  // Business data
        Time("resolved_at", time.Now()).               // Business data
        Bool("escalated", false))                      // Business data
```

## Audit Trail & Compliance Reporting

### Compliance Audit Logging

```go
// Track data access for compliance audits
emit.Info.Field("PII data accessed",
    emit.NewFields().
        String("accessed_by", "john.admin@company.com"). // Auto-masked
        String("employee_id", "EMP_12345").              // Auto-masked
        String("data_type", "customer_personal_info").   // Business data
        String("access_reason", "customer_support").     // Business data
        String("customer_affected", "CUST_67890").       // Auto-masked
        String("access_method", "admin_portal").         // Business data
        String("ip_address", "10.0.1.100").              // Auto-masked
        Time("access_time", time.Now()).                 // Business data
        Bool("authorized", true).                        // Business data
        String("session_id", "sess_audit_123"))          // Auto-masked

// Data retention compliance
emit.Info.Field("Data retention policy applied",
    emit.NewFields().
        String("data_type", "customer_communications").  // Business data
        Int("records_processed", 15000).                 // Business data
        Int("records_deleted", 8500).                    // Business data
        Int("records_archived", 6500).                   // Business data
        String("retention_period", "7_years").           // Business data
        String("legal_basis", "regulatory_requirement"). // Business data
        Time("policy_applied_at", time.Now()).           // Business data
        Bool("gdpr_compliant", true))                    // Business data
```

## Security Best Practices

### DO: Use Structured Logging

```go
// ✅ GOOD: Structured logging with automatic protection
emit.Info.Field("User login",
    emit.NewFields().
        String("username", username).      // Auto-protected
        String("ip", clientIP).            // Auto-protected
        Bool("success", true))

// ✅ GOOD: Key-value logging with automatic protection
emit.Info.KeyValue("Payment processed",
    "card_number", cardNumber,  // Auto-protected
    "amount", amount)           // Safe
```

### DON'T: Use String Formatting

```go
// ❌ BAD: String formatting can expose sensitive data
emit.Info.Msg(fmt.Sprintf("User %s with card %s paid %f",
    username, cardNumber, amount))  // Exposes PII!

// ❌ BAD: Manual concatenation
emit.Info.Msg("User: " + username + " Card: " + cardNumber) // Exposes PII!
```

### Security Checklist

- ✅ Use structured logging (`Field`, `KeyValue`, etc.)
- ✅ Enable PII masking in production (`EMIT_MASK_PII=true`)
- ✅ Enable sensitive data masking (`EMIT_MASK_SENSITIVE=true`)
- ✅ Configure industry-specific protected fields
- ✅ Use environment variables for security settings
- ✅ Regularly audit logged data for compliance
- ✅ Train developers on secure logging practices
- ❌ Never use string formatting for logs with user data
- ❌ Never disable masking in production environments
- ❌ Never log raw authentication credentials

## Zero-Configuration Security

The best part about emit's security features is that they work automatically:

```go
// Just use emit normally - security is automatic!
emit.Info.Field("User registered",
    emit.NewFields().
        String("email", "user@example.com").     // → "***PII***"
        String("password", "secret123").         // → "***MASKED***"
        String("name", "John Doe").              // → "***PII***"
        Int("user_id", 12345))                   // → 12345 (safe)

// No additional configuration required!
// Automatic compliance with GDPR, CCPA, HIPAA, PCI DSS
```

This makes emit the **secure-by-default** choice for logging in any Go application.
