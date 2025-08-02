package errors

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	err := New("Test error")
	if err == nil {
		t.Fatal("New() should not return nil")
	}

	if err.Message != "Test error" {
		t.Fatalf("Expected 'Test error', got '%s'", err.Message)
	}

	if err.Code != 500 {
		t.Fatalf("Expected code 500, got %d", err.Code)
	}

	if len(err.Stack) == 0 {
		t.Fatal("Stack should not be empty")
	}
}

func TestNewWithCode(t *testing.T) {
	err := NewWithCode(404, "Not found")
	if err == nil {
		t.Fatal("NewWithCode() should not return nil")
	}

	if err.Message != "Not found" {
		t.Fatalf("Expected 'Not found', got '%s'", err.Message)
	}

	if err.Code != 404 {
		t.Fatalf("Expected code 404, got %d", err.Code)
	}
}

func TestWrap(t *testing.T) {
	originalErr := New("Original error")
	wrappedErr := Wrap(originalErr, "Wrapped error")

	if wrappedErr == nil {
		t.Fatal("Wrap() should not return nil")
	}

	if wrappedErr.Message != "Wrapped error" {
		t.Fatalf("Expected 'Wrapped error', got '%s'", wrappedErr.Message)
	}

	if wrappedErr.Err != originalErr {
		t.Fatal("Wrapped error should contain original error")
	}
}

func TestWrapWithCode(t *testing.T) {
	originalErr := New("Original error")
	wrappedErr := WrapWithCode(originalErr, 422, "Validation error")

	if wrappedErr == nil {
		t.Fatal("WrapWithCode() should not return nil")
	}

	if wrappedErr.Message != "Validation error" {
		t.Fatalf("Expected 'Validation error', got '%s'", wrappedErr.Message)
	}

	if wrappedErr.Code != 422 {
		t.Fatalf("Expected code 422, got %d", wrappedErr.Code)
	}

	if wrappedErr.Err != originalErr {
		t.Fatal("Wrapped error should contain original error")
	}
}

func TestIsAppError(t *testing.T) {
	appErr := New("App error")
	regularErr := fmt.Errorf("Regular error")

	if !IsAppError(appErr) {
		t.Fatal("IsAppError() should return true for AppError")
	}

	if IsAppError(regularErr) {
		t.Fatal("IsAppError() should return false for regular error")
	}
}

func TestGetAppError(t *testing.T) {
	appErr := New("App error")
	regularErr := fmt.Errorf("Regular error")

	if result := GetAppError(appErr); result != appErr {
		t.Fatal("GetAppError() should return the same AppError")
	}

	if result := GetAppError(regularErr); result != nil {
		t.Fatal("GetAppError() should return nil for regular error")
	}
}

func TestErrorChaining(t *testing.T) {
	err := New("Base error").
		WithCode(400).
		WithMessage("Custom message").
		WithStack()

	if err.Code != 400 {
		t.Fatalf("Expected code 400, got %d", err.Code)
	}

	if err.Message != "Custom message" {
		t.Fatalf("Expected 'Custom message', got '%s'", err.Message)
	}

	if len(err.Stack) == 0 {
		t.Fatal("Stack should not be empty")
	}
}

func TestValidationError(t *testing.T) {
	validationErr := &ValidationError{
		Field:   "email",
		Message: "Invalid email format",
		Value:   "invalid-email",
	}

	expected := "Validation failed for field 'email': Invalid email format"
	if validationErr.Error() != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, validationErr.Error())
	}
}

func TestValidationErrors(t *testing.T) {
	errors := ValidationErrors{}

	// 添加验证错误
	errors.Add("name", "Name is required", "")
	errors.Add("email", "Invalid email format", "invalid-email")

	if !errors.HasErrors() {
		t.Fatal("HasErrors() should return true when there are errors")
	}

	if len(errors.GetErrors()) != 2 {
		t.Fatalf("Expected 2 errors, got %d", len(errors.GetErrors()))
	}

	// 测试按字段获取错误
	nameErrors := errors.GetErrorsByField("name")
	if len(nameErrors) != 1 {
		t.Fatalf("Expected 1 name error, got %d", len(nameErrors))
	}

	emailErrors := errors.GetErrorsByField("email")
	if len(emailErrors) != 1 {
		t.Fatalf("Expected 1 email error, got %d", len(emailErrors))
	}

	// 测试错误消息
	errorMessage := errors.Error()
	if errorMessage == "No validation errors" {
		t.Fatal("Error message should not be 'No validation errors' when there are errors")
	}
}

func TestValidationErrorsEmpty(t *testing.T) {
	errors := ValidationErrors{}

	if errors.HasErrors() {
		t.Fatal("HasErrors() should return false when there are no errors")
	}

	if len(errors.GetErrors()) != 0 {
		t.Fatalf("Expected 0 errors, got %d", len(errors.GetErrors()))
	}

	if errors.Error() != "No validation errors" {
		t.Fatalf("Expected 'No validation errors', got '%s'", errors.Error())
	}
}

// MockLogger 模拟日志记录器
type MockLogger struct {
	ErrorLogs   []map[string]interface{}
	WarningLogs []map[string]interface{}
	InfoLogs    []map[string]interface{}
	DebugLogs   []map[string]interface{}
}

func (m *MockLogger) Error(message string, context map[string]interface{}) {
	m.ErrorLogs = append(m.ErrorLogs, context)
}

func (m *MockLogger) Warning(message string, context map[string]interface{}) {
	m.WarningLogs = append(m.WarningLogs, context)
}

func (m *MockLogger) Info(message string, context map[string]interface{}) {
	m.InfoLogs = append(m.InfoLogs, context)
}

func (m *MockLogger) Debug(message string, context map[string]interface{}) {
	m.DebugLogs = append(m.DebugLogs, context)
}

func TestDefaultErrorHandler(t *testing.T) {
	logger := &MockLogger{}
	handler := NewDefaultErrorHandler(logger)

	// 测试处理应用错误
	appErr := New("Test app error")
	result := handler.Handle(appErr)

	if result != appErr {
		t.Fatal("Handler should return the same AppError")
	}

	if len(logger.ErrorLogs) != 1 {
		t.Fatalf("Expected 1 error log, got %d", len(logger.ErrorLogs))
	}

	// 测试处理普通错误
	regularErr := New("Regular error")
	result = handler.Handle(regularErr)

	if !IsAppError(result) {
		t.Fatal("Handler should wrap regular error as AppError")
	}

	if len(logger.ErrorLogs) != 2 {
		t.Fatalf("Expected 2 error logs, got %d", len(logger.ErrorLogs))
	}
}

func TestErrorMiddleware(t *testing.T) {
	logger := &MockLogger{}
	handler := NewDefaultErrorHandler(logger)
	middleware := NewErrorMiddleware(handler)

	// 测试正常执行
	next := func() interface{} {
		return "success"
	}

	result := middleware.Handle(nil, next)
	if result != "success" {
		t.Fatalf("Expected 'success', got '%v'", result)
	}

	// 测试错误处理
	nextWithError := func() interface{} {
		return New("Test error")
	}

	result = middleware.Handle(nil, nextWithError)
	if err, ok := result.(error); !ok || !IsAppError(err) {
		t.Fatal("Middleware should return AppError")
	}
}
