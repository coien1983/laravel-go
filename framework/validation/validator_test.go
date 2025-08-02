package validation

import (
	"fmt"
	"testing"
)

func TestNewValidator(t *testing.T) {
	validator := NewValidator()
	if validator == nil {
		t.Error("Expected validator to be created")
	}
}

func TestRequiredRule(t *testing.T) {
	validator := NewValidator()

	// 测试空值
	data := map[string]interface{}{
		"name": nil,
	}
	rules := map[string]string{
		"name": "required",
	}

	err := validator.Validate(data, rules)
	if err == nil {
		t.Error("Expected validation error for nil value")
	}

	// 测试空字符串
	data["name"] = ""
	err = validator.Validate(data, rules)
	if err == nil {
		t.Error("Expected validation error for empty string")
	}

	// 测试有效值
	data["name"] = "John Doe"
	err = validator.Validate(data, rules)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}
}

func TestStringRule(t *testing.T) {
	validator := NewValidator()

	data := map[string]interface{}{
		"name": "John Doe",
	}
	rules := map[string]string{
		"name": "string",
	}

	err := validator.Validate(data, rules)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}

	// 测试非字符串值
	data["name"] = 123
	err = validator.Validate(data, rules)
	if err == nil {
		t.Error("Expected validation error for non-string value")
	}
}

func TestIntRule(t *testing.T) {
	validator := NewValidator()

	// 测试整数
	data := map[string]interface{}{
		"age": 25,
	}
	rules := map[string]string{
		"age": "int",
	}

	err := validator.Validate(data, rules)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}

	// 测试字符串整数
	data["age"] = "30"
	err = validator.Validate(data, rules)
	if err != nil {
		t.Errorf("Expected no validation error for string integer, got: %v", err)
	}

	// 测试非整数值
	data["age"] = "not a number"
	err = validator.Validate(data, rules)
	if err == nil {
		t.Error("Expected validation error for non-integer value")
	}
}

func TestBoolRule(t *testing.T) {
	validator := NewValidator()

	// 测试布尔值
	data := map[string]interface{}{
		"active": true,
	}
	rules := map[string]string{
		"active": "bool",
	}

	err := validator.Validate(data, rules)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}

	// 测试字符串布尔值
	data["active"] = "true"
	err = validator.Validate(data, rules)
	if err != nil {
		t.Errorf("Expected no validation error for string boolean, got: %v", err)
	}

	// 测试非布尔值
	data["active"] = "not a boolean"
	err = validator.Validate(data, rules)
	if err == nil {
		t.Error("Expected validation error for non-boolean value")
	}
}

func TestEmailRule(t *testing.T) {
	validator := NewValidator()

	// 测试有效邮箱
	data := map[string]interface{}{
		"email": "john@example.com",
	}
	rules := map[string]string{
		"email": "email",
	}

	err := validator.Validate(data, rules)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}

	// 测试无效邮箱
	data["email"] = "invalid-email"
	err = validator.Validate(data, rules)
	if err == nil {
		t.Error("Expected validation error for invalid email")
	}

	// 测试空邮箱（应该通过）
	data["email"] = ""
	err = validator.Validate(data, rules)
	if err != nil {
		t.Errorf("Expected no validation error for empty email, got: %v", err)
	}
}

func TestMultipleRules(t *testing.T) {
	validator := NewValidator()

	data := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   25,
	}
	rules := map[string]string{
		"name":  "required|string",
		"email": "required|email",
		"age":   "required|int",
	}

	err := validator.Validate(data, rules)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}
}

func TestCustomRule(t *testing.T) {
	validator := NewValidator()

	// 注册自定义规则
	validator.RegisterRule("custom", RuleFunc(func(value interface{}) error {
		if value == nil {
			return nil
		}

		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field must be a string")
		}

		if len(str) < 3 {
			return fmt.Errorf("field must be at least 3 characters long")
		}

		return nil
	}))

	// 测试自定义规则
	data := map[string]interface{}{
		"name": "John",
	}
	rules := map[string]string{
		"name": "custom",
	}

	err := validator.Validate(data, rules)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}

	// 测试自定义规则失败
	data["name"] = "Jo"
	err = validator.Validate(data, rules)
	if err == nil {
		t.Error("Expected validation error for short string")
	}
}

func TestUnknownRule(t *testing.T) {
	validator := NewValidator()

	data := map[string]interface{}{
		"name": "John Doe",
	}
	rules := map[string]string{
		"name": "unknown_rule",
	}

	err := validator.Validate(data, rules)
	if err == nil {
		t.Error("Expected validation error for unknown rule")
	}
}

func TestNilValueHandling(t *testing.T) {
	validator := NewValidator()

	data := map[string]interface{}{
		"optional_field": nil,
	}
	rules := map[string]string{
		"optional_field": "string", // 非required规则应该允许nil值
	}

	err := validator.Validate(data, rules)
	if err != nil {
		t.Errorf("Expected no validation error for nil value with non-required rule, got: %v", err)
	}
}

func TestValidationErrors(t *testing.T) {
	validator := NewValidator()

	data := map[string]interface{}{
		"name":  "",              // 违反required规则
		"email": "invalid-email", // 违反email规则
		"age":   "not-a-number",  // 违反int规则
	}
	rules := map[string]string{
		"name":  "required",
		"email": "email",
		"age":   "int",
	}

	err := validator.Validate(data, rules)
	if err == nil {
		t.Error("Expected validation errors")
	}

	// 检查错误消息
	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}
}
