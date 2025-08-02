package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"laravel-go/framework/errors"
)

// Validator 验证器
type Validator struct {
	rules map[string]Rule
}

// Rule 验证规则接口
type Rule interface {
	Validate(value interface{}) error
}

// RuleFunc 验证规则函数
type RuleFunc func(value interface{}) error

// Validate 实现Rule接口
func (f RuleFunc) Validate(value interface{}) error {
	return f(value)
}

// NewValidator 创建新的验证器
func NewValidator() *Validator {
	v := &Validator{
		rules: make(map[string]Rule),
	}
	
	// 注册默认规则
	v.registerDefaultRules()
	
	return v
}

// RegisterRule 注册验证规则
func (v *Validator) RegisterRule(name string, rule Rule) {
	v.rules[name] = rule
}

// Validate 验证数据
func (v *Validator) Validate(data map[string]interface{}, rules map[string]string) error {
	var validationErrors errors.ValidationErrors
	
	for field, ruleString := range rules {
		value, _ := data[field]
		
		// 解析规则
		ruleParts := strings.Split(ruleString, "|")
		
		for _, rulePart := range ruleParts {
			ruleName := rulePart
			
			// 检查是否有参数
			if strings.Contains(rulePart, ":") {
				parts := strings.SplitN(rulePart, ":", 2)
				ruleName = parts[0]
				// 暂时不使用参数
			}
			
			// 获取规则
			rule, exists := v.rules[ruleName]
			if !exists {
				validationErrors.Add(field, fmt.Sprintf("Unknown validation rule: %s", ruleName), value)
				continue
			}
			
			// 执行验证
			if err := rule.Validate(value); err != nil {
				validationErrors.Add(field, err.Error(), value)
			}
		}
	}
	
	if validationErrors.HasErrors() {
		return validationErrors
	}
	
	return nil
}

// registerDefaultRules 注册默认规则
func (v *Validator) registerDefaultRules() {
	// required 规则
	v.RegisterRule("required", RuleFunc(func(value interface{}) error {
		if value == nil {
			return fmt.Errorf("field is required")
		}
		
		switch v := value.(type) {
		case string:
			if strings.TrimSpace(v) == "" {
				return fmt.Errorf("field is required")
			}
		case []string:
			if len(v) == 0 {
				return fmt.Errorf("field is required")
			}
		}
		
		return nil
	}))
	
	// string 规则
	v.RegisterRule("string", RuleFunc(func(value interface{}) error {
		if value == nil {
			return nil
		}
		
		if reflect.TypeOf(value).Kind() != reflect.String {
			return fmt.Errorf("field must be a string")
		}
		
		return nil
	}))
	
	// int 规则
	v.RegisterRule("int", RuleFunc(func(value interface{}) error {
		if value == nil {
			return nil
		}
		
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return nil
		case string:
			if _, err := strconv.Atoi(value.(string)); err != nil {
				return fmt.Errorf("field must be an integer")
			}
			return nil
		default:
			return fmt.Errorf("field must be an integer")
		}
	}))
	
	// bool 规则
	v.RegisterRule("bool", RuleFunc(func(value interface{}) error {
		if value == nil {
			return nil
		}
		
		switch value.(type) {
		case bool:
			return nil
		case string:
			if _, err := strconv.ParseBool(value.(string)); err != nil {
				return fmt.Errorf("field must be a boolean")
			}
			return nil
		default:
			return fmt.Errorf("field must be a boolean")
		}
	}))
	
	// email 规则
	v.RegisterRule("email", RuleFunc(func(value interface{}) error {
		if value == nil {
			return nil
		}
		
		email, ok := value.(string)
		if !ok {
			return fmt.Errorf("field must be a string")
		}
		
		if email == "" {
			return nil
		}
		
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(email) {
			return fmt.Errorf("field must be a valid email address")
		}
		
		return nil
	}))
	
	// min 规则
	v.RegisterRule("min", RuleFunc(func(value interface{}) error {
		if value == nil {
			return nil
		}
		
		// 这里需要参数，暂时跳过
		return nil
	}))
	
	// max 规则
	v.RegisterRule("max", RuleFunc(func(value interface{}) error {
		if value == nil {
			return nil
		}
		
		// 这里需要参数，暂时跳过
		return nil
	}))
	
	// unique 规则
	v.RegisterRule("unique", RuleFunc(func(value interface{}) error {
		if value == nil {
			return nil
		}
		
		// 这里需要数据库查询，暂时跳过
		return nil
	}))
} 