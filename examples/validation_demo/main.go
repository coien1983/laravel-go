package main

import (
	"fmt"
	"strconv"

	"laravel-go/framework/errors"
	"laravel-go/framework/validation"
)

// User 用户结构体
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password"`
	Active   bool   `json:"active"`
}

// Order 订单结构体
type Order struct {
	OrderID  string   `json:"order_id"`
	UserID   int64    `json:"user_id"`
	Amount   float64  `json:"amount"`
	Products []string `json:"products"`
	Status   string   `json:"status"`
}

func main() {
	fmt.Println("🚀 Laravel-Go 验证系统演示")
	fmt.Println("==================================================")

	// 创建验证器
	validator := validation.NewValidator()

	// 演示用户验证
	fmt.Println("\n📝 用户验证演示:")
	demoUserValidation(validator)

	// 演示订单验证
	fmt.Println("\n📦 订单验证演示:")
	demoOrderValidation(validator)

	// 演示自定义规则
	fmt.Println("\n🔧 自定义规则演示:")
	demoCustomRules(validator)

	// 演示错误处理
	fmt.Println("\n❌ 错误处理演示:")
	demoErrorHandling(validator)

	fmt.Println("\n✅ 验证系统演示完成!")
}

func demoUserValidation(validator *validation.Validator) {
	// 有效用户数据
	validUser := map[string]interface{}{
		"name":     "John Doe",
		"email":    "john@example.com",
		"age":      25,
		"password": "secret123",
		"active":   true,
	}

	userRules := map[string]string{
		"name":     "required|string",
		"email":    "required|email",
		"age":      "required|int",
		"password": "required|string",
		"active":   "bool",
	}

	fmt.Println("  验证有效用户数据...")
	err := validator.Validate(validUser, userRules)
	if err != nil {
		fmt.Printf("   ❌ 验证失败: %v\n", err)
	} else {
		fmt.Println("   ✅ 验证通过")
	}

	// 无效用户数据
	invalidUser := map[string]interface{}{
		"name":     "",              // 违反required规则
		"email":    "invalid-email", // 违反email规则
		"age":      "not-a-number",  // 违反int规则
		"password": nil,             // 违反required规则
		"active":   "not-a-boolean", // 违反bool规则
	}

	fmt.Println("  验证无效用户数据...")
	err = validator.Validate(invalidUser, userRules)
	if err != nil {
		fmt.Printf("   ❌ 验证失败 (预期): %v\n", err)

		// 检查是否是验证错误
		if validationErrors, ok := err.(errors.ValidationErrors); ok {
			fmt.Println("   详细错误信息:")
			for _, validationErr := range validationErrors.GetErrors() {
				fmt.Printf("     - %s: %s (值: %v)\n",
					validationErr.Field,
					validationErr.Message,
					validationErr.Value)
			}
		}
	} else {
		fmt.Println("   ✅ 验证通过 (意外)")
	}
}

func demoOrderValidation(validator *validation.Validator) {
	// 有效订单数据
	validOrder := map[string]interface{}{
		"order_id": "ORD-2024-001",
		"user_id":  12345,
		"amount":   299.99,
		"products": []string{"iPhone 15", "AirPods Pro"},
		"status":   "pending",
	}

	orderRules := map[string]string{
		"order_id": "required|string",
		"user_id":  "required|int",
		"amount":   "required",
		"products": "required",
		"status":   "required|string",
	}

	fmt.Println("  验证有效订单数据...")
	err := validator.Validate(validOrder, orderRules)
	if err != nil {
		fmt.Printf("   ❌ 验证失败: %v\n", err)
	} else {
		fmt.Println("   ✅ 验证通过")
	}

	// 无效订单数据
	invalidOrder := map[string]interface{}{
		"order_id": "",             // 违反required规则
		"user_id":  "not-a-number", // 违反int规则
		"amount":   -100,           // 负数金额
		"products": []string{},     // 空数组
		"status":   nil,            // 违反required规则
	}

	fmt.Println("  验证无效订单数据...")
	err = validator.Validate(invalidOrder, orderRules)
	if err != nil {
		fmt.Printf("   ❌ 验证失败 (预期): %v\n", err)
	} else {
		fmt.Println("   ✅ 验证通过 (意外)")
	}
}

func demoCustomRules(validator *validation.Validator) {
	// 注册自定义规则
	validator.RegisterRule("min_length", validation.RuleFunc(func(value interface{}) error {
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

	validator.RegisterRule("positive_number", validation.RuleFunc(func(value interface{}) error {
		if value == nil {
			return nil
		}

		switch v := value.(type) {
		case int:
			if v <= 0 {
				return fmt.Errorf("field must be a positive number")
			}
		case float64:
			if v <= 0 {
				return fmt.Errorf("field must be a positive number")
			}
		case string:
			// 尝试转换为数字
			if num, err := strconv.Atoi(v); err == nil {
				if num <= 0 {
					return fmt.Errorf("field must be a positive number")
				}
			}
		}

		return nil
	}))

	// 使用自定义规则
	data := map[string]interface{}{
		"username": "john",
		"age":      25,
	}

	rules := map[string]string{
		"username": "required|min_length",
		"age":      "required|positive_number",
	}

	fmt.Println("  验证自定义规则...")
	err := validator.Validate(data, rules)
	if err != nil {
		fmt.Printf("   ❌ 验证失败: %v\n", err)
	} else {
		fmt.Println("   ✅ 验证通过")
	}

	// 测试自定义规则失败
	invalidData := map[string]interface{}{
		"username": "jo", // 违反min_length规则
		"age":      -5,   // 违反positive_number规则
	}

	fmt.Println("  验证自定义规则失败情况...")
	err = validator.Validate(invalidData, rules)
	if err != nil {
		fmt.Printf("   ❌ 验证失败 (预期): %v\n", err)
	} else {
		fmt.Println("   ✅ 验证通过 (意外)")
	}
}

func demoErrorHandling(validator *validation.Validator) {
	// 测试未知规则
	data := map[string]interface{}{
		"field": "value",
	}
	rules := map[string]string{
		"field": "unknown_rule",
	}

	fmt.Println("  测试未知规则...")
	err := validator.Validate(data, rules)
	if err != nil {
		fmt.Printf("   ❌ 验证失败 (预期): %v\n", err)
	} else {
		fmt.Println("   ✅ 验证通过 (意外)")
	}

	// 测试多个错误
	multiErrorData := map[string]interface{}{
		"name":  "",        // 违反required规则
		"email": "invalid", // 违反email规则
		"age":   "young",   // 违反int规则
	}
	multiErrorRules := map[string]string{
		"name":  "required",
		"email": "email",
		"age":   "int",
	}

	fmt.Println("  测试多个错误...")
	err = validator.Validate(multiErrorData, multiErrorRules)
	if err != nil {
		fmt.Printf("   ❌ 验证失败 (预期): %v\n", err)

		// 检查错误类型
		if validationErrors, ok := err.(errors.ValidationErrors); ok {
			fmt.Printf("   错误数量: %d\n", len(validationErrors.GetErrors()))

			// 按字段分组错误
			errorsByField := make(map[string][]string)
			for _, validationErr := range validationErrors.GetErrors() {
				errorsByField[validationErr.Field] = append(
					errorsByField[validationErr.Field],
					validationErr.Message)
			}

			fmt.Println("   按字段分组的错误:")
			for field, messages := range errorsByField {
				fmt.Printf("     %s: %v\n", field, messages)
			}
		}
	} else {
		fmt.Println("   ✅ 验证通过 (意外)")
	}
}
