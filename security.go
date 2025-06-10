package emit

import "strings"

// Default sensitive field patterns (case-insensitive)
var defaultSensitiveFields = []string{
	"password", "pwd", "pass", "secret", "key", "token", "auth",
	"credential", "cred", "private", "confidential", "sensitive",
	"api_key", "apikey", "access_token", "refresh_token", "jwt",
	"session", "cookie", "authorization", "bearer", "oauth",
	"client_secret", "private_key", "passphrase", "pin", "code",
}

// Default PII field patterns (case-insensitive)
var defaultPIIFields = []string{
	"email", "mail", "e_mail", "email_address", "emailaddress",
	"phone", "mobile", "telephone", "phone_number", "phonenumber", "tel",
	"ssn", "social_security", "social_security_number", "tax_id",
	"credit_card", "creditcard", "card_number", "cardnumber", "ccn",
	"passport", "passport_number", "license", "driver_license", "dl",
	"name", "first_name", "last_name", "full_name", "firstname", "lastname", "fullname",
	"address", "street", "city", "zip", "zipcode", "postal", "postal_code",
	"ip", "ip_address", "ipaddress", "user_agent", "useragent",
	"dob", "date_of_birth", "birthdate", "birthday", "birth_date",
	"iban", "account_number", "bank_account", "routing_number",
	"username", "user_name", "login", "userid",
}

// isPIIField checks if a field name matches PII patterns
func (l *Logger) isPIIField(fieldName string) bool {
	if l.piiMode == SHOW_PII {
		return false
	}

	lowerFieldName := strings.ToLower(fieldName)
	for _, pattern := range l.piiFields {
		if strings.Contains(lowerFieldName, pattern) {
			return true
		}
	}
	return false
}

// isSensitiveField checks if a field name matches sensitive patterns
func (l *Logger) isSensitiveField(fieldName string) bool {
	if l.sensitiveMode == SHOW_SENSITIVE {
		return false
	}

	lowerFieldName := strings.ToLower(fieldName)
	for _, pattern := range l.sensitiveFields {
		if strings.Contains(lowerFieldName, pattern) {
			return true
		}
	}
	return false
}

// maskSensitiveFields recursively masks sensitive and PII data in fields
func (l *Logger) maskSensitiveFields(fields map[string]any) map[string]any {
	if (l.sensitiveMode == SHOW_SENSITIVE && l.piiMode == SHOW_PII) || len(fields) == 0 {
		return fields
	}

	maskedFields := make(map[string]any)
	for key, value := range fields {
		// Check PII first (more specific), then sensitive data
		if l.isPIIField(key) {
			maskedFields[key] = l.piiMaskString
		} else if l.isSensitiveField(key) {
			maskedFields[key] = l.maskString
		} else {
			// Handle nested maps
			if nestedMap, ok := value.(map[string]any); ok {
				maskedFields[key] = l.maskSensitiveFields(nestedMap)
			} else {
				maskedFields[key] = value
			}
		}
	}
	return maskedFields
}
