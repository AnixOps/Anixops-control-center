package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Rule represents a validation rule
type Rule struct {
	Name    string
	Apply   func(value interface{}) error
	Message string
}

// Validator validates configuration
type Validator struct {
	rules map[string][]Rule
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		rules: make(map[string][]Rule),
	}
}

// AddRule adds a validation rule for a field
func (v *Validator) AddRule(field string, rule Rule) {
	v.rules[field] = append(v.rules[field], rule)
}

// Required adds a required rule
func Required() Rule {
	return Rule{
		Name: "required",
		Apply: func(value interface{}) error {
			if value == nil || value == "" {
				return fmt.Errorf("is required")
			}
			rv := reflect.ValueOf(value)
			if rv.Kind() == reflect.Ptr && rv.IsNil() {
				return fmt.Errorf("is required")
			}
			if rv.Kind() == reflect.String && rv.String() == "" {
				return fmt.Errorf("is required")
			}
			return nil
		},
		Message: "field is required",
	}
}

// Min adds a minimum value rule
func Min(min interface{}) Rule {
	return Rule{
		Name: "min",
		Apply: func(value interface{}) error {
			rv := reflect.ValueOf(value)
			minv := reflect.ValueOf(min)

			switch rv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if rv.Int() < minv.Int() {
					return fmt.Errorf("must be at least %v", min)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if rv.Uint() < minv.Uint() {
					return fmt.Errorf("must be at least %v", min)
				}
			case reflect.Float32, reflect.Float64:
				if rv.Float() < minv.Float() {
					return fmt.Errorf("must be at least %v", min)
				}
			case reflect.String:
				if len(rv.String()) < int(minv.Int()) {
					return fmt.Errorf("must be at least %d characters", min)
				}
			case reflect.Slice, reflect.Array:
				if rv.Len() < int(minv.Int()) {
					return fmt.Errorf("must have at least %d items", min)
				}
			}
			return nil
		},
		Message: fmt.Sprintf("minimum value is %v", min),
	}
}

// Max adds a maximum value rule
func Max(max interface{}) Rule {
	return Rule{
		Name: "max",
		Apply: func(value interface{}) error {
			rv := reflect.ValueOf(value)
			maxv := reflect.ValueOf(max)

			switch rv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if rv.Int() > maxv.Int() {
					return fmt.Errorf("must be at most %v", max)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if rv.Uint() > maxv.Uint() {
					return fmt.Errorf("must be at most %v", max)
				}
			case reflect.Float32, reflect.Float64:
				if rv.Float() > maxv.Float() {
					return fmt.Errorf("must be at most %v", max)
				}
			case reflect.String:
				if len(rv.String()) > int(maxv.Int()) {
					return fmt.Errorf("must be at most %d characters", max)
				}
			case reflect.Slice, reflect.Array:
				if rv.Len() > int(maxv.Int()) {
					return fmt.Errorf("must have at most %d items", max)
				}
			}
			return nil
		},
		Message: fmt.Sprintf("maximum value is %v", max),
	}
}

// Range adds a range rule
func Range(min, max interface{}) Rule {
	return Rule{
		Name: "range",
		Apply: func(value interface{}) error {
			if err := Min(min).Apply(value); err != nil {
				return err
			}
			if err := Max(max).Apply(value); err != nil {
				return err
			}
			return nil
		},
		Message: fmt.Sprintf("must be between %v and %v", min, max),
	}
}

// Regex adds a regex pattern rule
func Regex(pattern string) Rule {
	re := regexp.MustCompile(pattern)
	return Rule{
		Name: "regex",
		Apply: func(value interface{}) error {
			s, ok := value.(string)
			if !ok {
				return fmt.Errorf("must be a string")
			}
			if !re.MatchString(s) {
				return fmt.Errorf("must match pattern %s", pattern)
			}
			return nil
		},
		Message: fmt.Sprintf("must match pattern %s", pattern),
	}
}

// Email adds an email validation rule
func Email() Rule {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return Rule{
		Name: "email",
		Apply: func(value interface{}) error {
			s, ok := value.(string)
			if !ok {
				return fmt.Errorf("must be a string")
			}
			if !emailRegex.MatchString(s) {
				return fmt.Errorf("must be a valid email address")
			}
			return nil
		},
		Message: "must be a valid email address",
	}
}

// URL adds a URL validation rule
func URL() Rule {
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return Rule{
		Name: "url",
		Apply: func(value interface{}) error {
			s, ok := value.(string)
			if !ok {
				return fmt.Errorf("must be a string")
			}
			if !urlRegex.MatchString(s) {
				return fmt.Errorf("must be a valid URL")
			}
			return nil
		},
		Message: "must be a valid URL",
	}
}

// OneOf adds a one-of rule
func OneOf(values ...string) Rule {
	return Rule{
		Name: "one_of",
		Apply: func(value interface{}) error {
			s, ok := value.(string)
			if !ok {
				return fmt.Errorf("must be a string")
			}
			for _, v := range values {
				if s == v {
					return nil
				}
			}
			return fmt.Errorf("must be one of: %s", strings.Join(values, ", "))
		},
		Message: fmt.Sprintf("must be one of: %s", strings.Join(values, ", ")),
	}
}

// Custom adds a custom validation rule
func Custom(name string, fn func(value interface{}) error) Rule {
	return Rule{
		Name: name,
		Apply: fn,
		Message: name + " validation failed",
	}
}

// Validate validates a value against all rules for a field
func (v *Validator) Validate(field string, value interface{}) []error {
	rules, exists := v.rules[field]
	if !exists {
		return nil
	}

	var errs []error
	for _, rule := range rules {
		if err := rule.Apply(value); err != nil {
			errs = append(errs, fmt.Errorf("%s: %s", field, err.Error()))
		}
	}
	return errs
}

// ValidateAll validates all fields
func (v *Validator) ValidateAll(data map[string]interface{}) map[string][]error {
	result := make(map[string][]error)
	for field, value := range data {
		if errs := v.Validate(field, value); len(errs) > 0 {
			result[field] = errs
		}
	}
	return result
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	Valid  bool
	Errors map[string][]error
}

// Valid returns a successful validation result
func Valid() ValidationResult {
	return ValidationResult{Valid: true, Errors: nil}
}

// Invalid returns a failed validation result
func Invalid(errors map[string][]error) ValidationResult {
	return ValidationResult{Valid: false, Errors: errors}
}

// ConfigValidator validates configuration structures
type ConfigValidator struct {
	validator *Validator
}

// NewConfigValidator creates a new config validator
func NewConfigValidator() *ConfigValidator {
	v := NewValidator()

	// Common config validations
	v.AddRule("server.port", Min(1))
	v.AddRule("server.port", Max(65535))
	v.AddRule("database.driver", OneOf("sqlite", "postgres", "mysql"))
	v.AddRule("auth.provider", OneOf("local", "oauth", "ldap", "saml"))
	v.AddRule("logging.level", OneOf("debug", "info", "warn", "error"))
	v.AddRule("logging.format", OneOf("json", "text"))

	return &ConfigValidator{validator: v}
}

// ValidateServerConfig validates server configuration
func (cv *ConfigValidator) ValidateServerConfig(config map[string]interface{}) ValidationResult {
	errs := make(map[string][]error)

	if port, ok := config["port"]; ok {
		if e := cv.validator.Validate("server.port", port); len(e) > 0 {
			errs["port"] = e
		}
	}

	if host, ok := config["host"]; ok {
		if host == "" {
			errs["host"] = []error{fmt.Errorf("host cannot be empty")}
		}
	}

	if len(errs) == 0 {
		return Valid()
	}
	return Invalid(errs)
}

// ValidateDatabaseConfig validates database configuration
func (cv *ConfigValidator) ValidateDatabaseConfig(config map[string]interface{}) ValidationResult {
	errs := make(map[string][]error)

	if driver, ok := config["driver"]; ok {
		if e := cv.validator.Validate("database.driver", driver); len(e) > 0 {
			errs["driver"] = e
		}
	}

	if driver, ok := config["driver"].(string); ok && driver == "postgres" {
		if _, ok := config["host"]; !ok {
			errs["host"] = []error{fmt.Errorf("host is required for postgres")}
		}
		if _, ok := config["database"]; !ok {
			errs["database"] = []error{fmt.Errorf("database name is required for postgres")}
		}
	}

	if len(errs) == 0 {
		return Valid()
	}
	return Invalid(errs)
}

// ValidateAuthConfig validates authentication configuration
func (cv *ConfigValidator) ValidateAuthConfig(config map[string]interface{}) ValidationResult {
	errs := make(map[string][]error)

	if provider, ok := config["provider"]; ok {
		if e := cv.validator.Validate("auth.provider", provider); len(e) > 0 {
			errs["provider"] = e
		}
	}

	if jwt, ok := config["jwt"].(map[string]interface{}); ok {
		if secret, ok := jwt["secret"].(string); ok {
			if len(secret) < 32 {
				errs["jwt.secret"] = []error{fmt.Errorf("jwt secret must be at least 32 characters")}
			}
		}
	}

	if len(errs) == 0 {
		return Valid()
	}
	return Invalid(errs)
}

// PluginConfigValidator validates plugin configurations
type PluginConfigValidator struct {
	schemas map[string]map[string][]Rule
}

// NewPluginConfigValidator creates a new plugin config validator
func NewPluginConfigValidator() *PluginConfigValidator {
	return &PluginConfigValidator{
		schemas: make(map[string]map[string][]Rule),
	}
}

// RegisterSchema registers a validation schema for a plugin
func (v *PluginConfigValidator) RegisterSchema(pluginName string, schema map[string][]Rule) {
	v.schemas[pluginName] = schema
}

// Validate validates a plugin configuration
func (v *PluginConfigValidator) Validate(pluginName string, config map[string]interface{}) ValidationResult {
	schema, exists := v.schemas[pluginName]
	if !exists {
		return Valid() // No schema, assume valid
	}

	errs := make(map[string][]error)
	for field, rules := range schema {
		value := config[field]
		for _, rule := range rules {
			if err := rule.Apply(value); err != nil {
				errs[field] = append(errs[field], fmt.Errorf("%s: %s", field, err.Error()))
			}
		}
	}

	if len(errs) == 0 {
		return Valid()
	}
	return Invalid(errs)
}

// Default validator instance
var defaultValidator = NewConfigValidator()

// ValidateConfig validates configuration using the default validator
func ValidateConfig(config map[string]interface{}) ValidationResult {
	return defaultValidator.ValidateServerConfig(config)
}

// ValidatePluginConfig validates plugin configuration
func ValidatePluginConfig(pluginName string, config map[string]interface{}) ValidationResult {
	return defaultValidator.ValidateDatabaseConfig(config)
}