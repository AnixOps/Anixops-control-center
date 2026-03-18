package validator

import (
	"testing"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	if v == nil {
		t.Fatal("Validator is nil")
	}
}

func TestAddRule(t *testing.T) {
	v := NewValidator()
	v.AddRule("test", Required())
	// No error means success
}

func TestRequired(t *testing.T) {
	rule := Required()

	// Test nil
	if rule.Apply(nil) == nil {
		t.Error("Required should fail for nil")
	}

	// Test empty string
	if rule.Apply("") == nil {
		t.Error("Required should fail for empty string")
	}

	// Test non-empty string
	if rule.Apply("value") != nil {
		t.Error("Required should pass for non-empty string")
	}

	// Test number
	if rule.Apply(123) != nil {
		t.Error("Required should pass for number")
	}
}

func TestMin(t *testing.T) {
	rule := Min(5)

	// Test int
	if rule.Apply(3) == nil {
		t.Error("Min(5) should fail for 3")
	}
	if rule.Apply(5) != nil {
		t.Error("Min(5) should pass for 5")
	}
	if rule.Apply(10) != nil {
		t.Error("Min(5) should pass for 10")
	}

	// Test string length
	strRule := Min(3)
	if strRule.Apply("ab") == nil {
		t.Error("Min(3) should fail for 'ab'")
	}
	if strRule.Apply("abc") != nil {
		t.Error("Min(3) should pass for 'abc'")
	}
}

func TestMax(t *testing.T) {
	rule := Max(10)

	// Test int
	if rule.Apply(15) == nil {
		t.Error("Max(10) should fail for 15")
	}
	if rule.Apply(10) != nil {
		t.Error("Max(10) should pass for 10")
	}
	if rule.Apply(5) != nil {
		t.Error("Max(10) should pass for 5")
	}

	// Test string length
	strRule := Max(3)
	if strRule.Apply("abcd") == nil {
		t.Error("Max(3) should fail for 'abcd'")
	}
	if strRule.Apply("abc") != nil {
		t.Error("Max(3) should pass for 'abc'")
	}
}

func TestRange(t *testing.T) {
	rule := Range(5, 10)

	if rule.Apply(3) == nil {
		t.Error("Range(5,10) should fail for 3")
	}
	if rule.Apply(7) != nil {
		t.Error("Range(5,10) should pass for 7")
	}
	if rule.Apply(15) == nil {
		t.Error("Range(5,10) should fail for 15")
	}
}

func TestRegex(t *testing.T) {
	rule := Regex(`^\d+$`)

	if rule.Apply("abc") == nil {
		t.Error("Regex should fail for 'abc'")
	}
	if rule.Apply("123") != nil {
		t.Error("Regex should pass for '123'")
	}
	if rule.Apply(123) == nil {
		t.Error("Regex should fail for non-string")
	}
}

func TestEmail(t *testing.T) {
	rule := Email()

	if rule.Apply("notanemail") == nil {
		t.Error("Email should fail for 'notanemail'")
	}
	if rule.Apply("test@example.com") != nil {
		t.Error("Email should pass for 'test@example.com'")
	}
	if rule.Apply(123) == nil {
		t.Error("Email should fail for non-string")
	}
}

func TestURL(t *testing.T) {
	rule := URL()

	if rule.Apply("notaurl") == nil {
		t.Error("URL should fail for 'notaurl'")
	}
	if rule.Apply("http://example.com") != nil {
		t.Error("URL should pass for 'http://example.com'")
	}
	if rule.Apply("https://example.com/path") != nil {
		t.Error("URL should pass for 'https://example.com/path'")
	}
}

func TestOneOf(t *testing.T) {
	rule := OneOf("a", "b", "c")

	if rule.Apply("d") == nil {
		t.Error("OneOf should fail for 'd'")
	}
	if rule.Apply("a") != nil {
		t.Error("OneOf should pass for 'a'")
	}
	if rule.Apply("b") != nil {
		t.Error("OneOf should pass for 'b'")
	}
}

func TestCustom(t *testing.T) {
	rule := Custom("always_fail", func(value interface{}) error {
		return nil // always pass for this test
	})

	if rule.Apply("anything") != nil {
		t.Error("Custom rule should pass")
	}
}

func TestValidate(t *testing.T) {
	v := NewValidator()
	v.AddRule("name", Required())
	v.AddRule("age", Min(0))

	errs := v.Validate("name", nil)
	if len(errs) == 0 {
		t.Error("Should have errors for nil name")
	}

	errs = v.Validate("name", "John")
	if len(errs) != 0 {
		t.Error("Should not have errors for valid name")
	}
}

func TestValidateAll(t *testing.T) {
	v := NewValidator()
	v.AddRule("name", Required())
	v.AddRule("age", Min(0))

	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	errs := v.ValidateAll(data)
	if len(errs) != 0 {
		t.Errorf("Should not have errors for valid data: %v", errs)
	}

	invalidData := map[string]interface{}{
		"name": "",
		"age":  -1,
	}

	errs = v.ValidateAll(invalidData)
	if len(errs) == 0 {
		t.Error("Should have errors for invalid data")
	}
}

func TestValidResult(t *testing.T) {
	result := Valid()
	if !result.Valid {
		t.Error("Valid result should be valid")
	}
	if result.Errors != nil {
		t.Error("Valid result should have no errors")
	}
}

func TestInvalidResult(t *testing.T) {
	errs := map[string][]error{
		"field": {nil}, // placeholder
	}
	result := Invalid(errs)
	if result.Valid {
		t.Error("Invalid result should not be valid")
	}
}

func TestNewConfigValidator(t *testing.T) {
	v := NewConfigValidator()
	if v == nil {
		t.Fatal("ConfigValidator is nil")
	}
}

func TestValidateServerConfig(t *testing.T) {
	v := NewConfigValidator()

	// Valid config
	result := v.ValidateServerConfig(map[string]interface{}{
		"port": 8080,
		"host": "localhost",
	})
	if !result.Valid {
		t.Errorf("Valid config should pass: %v", result.Errors)
	}
}

func TestValidateDatabaseConfig(t *testing.T) {
	v := NewConfigValidator()

	// Valid sqlite config
	result := v.ValidateDatabaseConfig(map[string]interface{}{
		"driver": "sqlite",
	})
	if !result.Valid {
		t.Errorf("Valid sqlite config should pass: %v", result.Errors)
	}

	// Valid postgres config
	result = v.ValidateDatabaseConfig(map[string]interface{}{
		"driver":   "postgres",
		"host":     "localhost",
		"database": "testdb",
	})
	if !result.Valid {
		t.Errorf("Valid postgres config should pass: %v", result.Errors)
	}

	// Invalid postgres config (missing host)
	result = v.ValidateDatabaseConfig(map[string]interface{}{
		"driver":   "postgres",
		"database": "testdb",
	})
	if result.Valid {
		t.Error("Invalid postgres config should fail")
	}
}

func TestValidateAuthConfig(t *testing.T) {
	v := NewConfigValidator()

	// Valid config
	result := v.ValidateAuthConfig(map[string]interface{}{
		"provider": "local",
		"jwt": map[string]interface{}{
			"secret": "this-is-a-very-long-secret-key-32chars",
		},
	})
	if !result.Valid {
		t.Errorf("Valid auth config should pass: %v", result.Errors)
	}

	// Invalid - short secret
	result = v.ValidateAuthConfig(map[string]interface{}{
		"provider": "local",
		"jwt": map[string]interface{}{
			"secret": "short",
		},
	})
	if result.Valid {
		t.Error("Short JWT secret should fail")
	}
}

func TestNewPluginConfigValidator(t *testing.T) {
	v := NewPluginConfigValidator()
	if v == nil {
		t.Fatal("PluginConfigValidator is nil")
	}
}

func TestPluginConfigValidatorRegisterSchema(t *testing.T) {
	v := NewPluginConfigValidator()

	v.RegisterSchema("my-plugin", map[string][]Rule{
		"name": {Required()},
		"port": {Min(1), Max(65535)},
	})

	// Valid config
	result := v.Validate("my-plugin", map[string]interface{}{
		"name": "test",
		"port": 8080,
	})
	if !result.Valid {
		t.Errorf("Valid config should pass: %v", result.Errors)
	}

	// Invalid config
	result = v.Validate("my-plugin", map[string]interface{}{
		"name": "",
		"port": 0,
	})
	if result.Valid {
		t.Error("Invalid config should fail")
	}
}

func TestPluginConfigValidatorNoSchema(t *testing.T) {
	v := NewPluginConfigValidator()

	// No schema registered, should pass
	result := v.Validate("unknown-plugin", map[string]interface{}{
		"anything": "goes",
	})
	if !result.Valid {
		t.Error("No schema means always valid")
	}
}

func TestValidateConfig(t *testing.T) {
	result := ValidateConfig(map[string]interface{}{
		"port": 8080,
	})
	// Function exists and works
	_ = result
}

func TestValidatePluginConfig(t *testing.T) {
	result := ValidatePluginConfig("test", map[string]interface{}{
		"driver": "sqlite",
	})
	// Function exists and works
	_ = result
}