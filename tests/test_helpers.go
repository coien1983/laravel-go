package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"laravel-go/framework/database"
)

// TestHelper 测试辅助工具
type TestHelper struct {
	t *testing.T
}

// NewTestHelper 创建测试辅助工具
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// CreateTempFile 创建临时文件
func (th *TestHelper) CreateTempFile(content string) string {
	tmpfile, err := os.CreateTemp("", "test_*.tmp")
	require.NoError(th.t, err)
	defer tmpfile.Close()

	_, err = tmpfile.WriteString(content)
	require.NoError(th.t, err)

	return tmpfile.Name()
}

// CreateTempDir 创建临时目录
func (th *TestHelper) CreateTempDir() string {
	tmpdir, err := os.MkdirTemp("", "test_*")
	require.NoError(th.t, err)
	return tmpdir
}

// CleanupTempFile 清理临时文件
func (th *TestHelper) CleanupTempFile(filename string) {
	if filename != "" {
		os.Remove(filename)
	}
}

// CleanupTempDir 清理临时目录
func (th *TestHelper) CleanupTempDir(dirname string) {
	if dirname != "" {
		os.RemoveAll(dirname)
	}
}

// CreateTestDatabase 创建测试数据库
func (th *TestHelper) CreateTestDatabase() database.Connection {
	db, err := database.NewConnection(&database.Config{
		Driver:   "sqlite",
		Database: ":memory:",
	})
	require.NoError(th.t, err)
	return db
}

// CreateTestTable 创建测试表
func (th *TestHelper) CreateTestTable(db database.Connection, tableName, schema string) {
	_, err := db.Exec(schema)
	require.NoError(th.t, err)
}

// InsertTestData 插入测试数据
func (th *TestHelper) InsertTestData(db database.Connection, tableName string, data map[string]interface{}) {
	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	placeholders := make([]string, 0, len(data))

	for column, value := range data {
		columns = append(columns, column)
		values = append(values, value)
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	_, err := db.Exec(query, values...)
	require.NoError(th.t, err)
}

// CleanupTestData 清理测试数据
func (th *TestHelper) CleanupTestData(db database.Connection, tableName string) {
	_, err := db.Exec("DELETE FROM " + tableName)
	require.NoError(th.t, err)
}

// CreateHTTPRequest 创建HTTP请求
func (th *TestHelper) CreateHTTPRequest(method, url string, body interface{}) *http.Request {
	var reqBody io.Reader

	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = strings.NewReader(v)
		case []byte:
			reqBody = bytes.NewReader(v)
		default:
			jsonData, err := json.Marshal(body)
			require.NoError(th.t, err)
			reqBody = bytes.NewReader(jsonData)
		}
	}

	req, err := http.NewRequest(method, url, reqBody)
	require.NoError(th.t, err)

	if body != nil && reflect.TypeOf(body).Kind() != reflect.String && reflect.TypeOf(body).Kind() != reflect.Slice {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

// ExecuteHTTPRequest 执行HTTP请求
func (th *TestHelper) ExecuteHTTPRequest(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder
}

// AssertHTTPResponse 断言HTTP响应
func (th *TestHelper) AssertHTTPResponse(recorder *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	assert.Equal(th.t, expectedStatus, recorder.Code)

	if expectedBody != nil {
		switch v := expectedBody.(type) {
		case string:
			assert.Equal(th.t, v, recorder.Body.String())
		case []byte:
			assert.Equal(th.t, v, recorder.Body.Bytes())
		default:
			var expectedJSON, actualJSON interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &actualJSON)
			require.NoError(th.t, err)

			expectedBytes, err := json.Marshal(expectedBody)
			require.NoError(th.t, err)
			err = json.Unmarshal(expectedBytes, &expectedJSON)
			require.NoError(th.t, err)

			assert.Equal(th.t, expectedJSON, actualJSON)
		}
	}
}

// AssertJSONResponse 断言JSON响应
func (th *TestHelper) AssertJSONResponse(recorder *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	assert.Equal(th.t, expectedStatus, recorder.Code)
	assert.Equal(th.t, "application/json", recorder.Header().Get("Content-Type"))

	var actualBody interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &actualBody)
	require.NoError(th.t, err)

	var expectedJSON interface{}
	expectedBytes, err := json.Marshal(expectedBody)
	require.NoError(th.t, err)
	err = json.Unmarshal(expectedBytes, &expectedJSON)
	require.NoError(th.t, err)

	assert.Equal(th.t, expectedJSON, actualBody)
}

// WaitForCondition 等待条件满足
func (th *TestHelper) WaitForCondition(condition func() bool, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// RetryOperation 重试操作
func (th *TestHelper) RetryOperation(operation func() error, maxRetries int, delay time.Duration) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := operation(); err == nil {
			return nil
		} else {
			lastErr = err
			if i < maxRetries-1 {
				time.Sleep(delay)
			}
		}
	}
	return lastErr
}

// GenerateRandomString 生成随机字符串
func (th *TestHelper) GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// GenerateRandomEmail 生成随机邮箱
func (th *TestHelper) GenerateRandomEmail() string {
	return fmt.Sprintf("test_%s@example.com", th.GenerateRandomString(8))
}

// GenerateRandomInt 生成随机整数
func (th *TestHelper) GenerateRandomInt(min, max int) int {
	return int(time.Now().UnixNano())%(max-min+1) + min
}

// AssertFileExists 断言文件存在
func (th *TestHelper) AssertFileExists(filename string) {
	_, err := os.Stat(filename)
	assert.NoError(th.t, err)
}

// AssertFileNotExists 断言文件不存在
func (th *TestHelper) AssertFileNotExists(filename string) {
	_, err := os.Stat(filename)
	assert.True(th.t, os.IsNotExist(err))
}

// AssertDirExists 断言目录存在
func (th *TestHelper) AssertDirExists(dirname string) {
	info, err := os.Stat(dirname)
	assert.NoError(th.t, err)
	assert.True(th.t, info.IsDir())
}

// AssertFileContent 断言文件内容
func (th *TestHelper) AssertFileContent(filename, expectedContent string) {
	content, err := os.ReadFile(filename)
	require.NoError(th.t, err)
	assert.Equal(th.t, expectedContent, string(content))
}

// CreateTestDataFile 创建测试数据文件
func (th *TestHelper) CreateTestDataFile(filename string, data interface{}) {
	file, err := os.Create(filename)
	require.NoError(th.t, err)
	defer file.Close()

	switch v := data.(type) {
	case string:
		_, err = file.WriteString(v)
	case []byte:
		_, err = file.Write(v)
	default:
		jsonData, err := json.Marshal(data)
		require.NoError(th.t, err)
		_, err = file.Write(jsonData)
	}
	require.NoError(th.t, err)
}

// LoadTestDataFile 加载测试数据文件
func (th *TestHelper) LoadTestDataFile(filename string) []byte {
	data, err := os.ReadFile(filename)
	require.NoError(th.t, err)
	return data
}

// GetTestDataPath 获取测试数据路径
func (th *TestHelper) GetTestDataPath(filename string) string {
	return filepath.Join("testdata", filename)
}

// EnsureTestDataDir 确保测试数据目录存在
func (th *TestHelper) EnsureTestDataDir() {
	err := os.MkdirAll("testdata", 0755)
	require.NoError(th.t, err)
}

// AssertStructEqual 断言结构体相等
func (th *TestHelper) AssertStructEqual(expected, actual interface{}) {
	assert.Equal(th.t, expected, actual)
}

// AssertStructNotEqual 断言结构体不相等
func (th *TestHelper) AssertStructNotEqual(expected, actual interface{}) {
	assert.NotEqual(th.t, expected, actual)
}

// AssertMapContains 断言Map包含指定键值对
func (th *TestHelper) AssertMapContains(m map[string]interface{}, key string, value interface{}) {
	assert.Contains(th.t, m, key)
	assert.Equal(th.t, value, m[key])
}

// AssertSliceContains 断言Slice包含指定元素
func (th *TestHelper) AssertSliceContains(slice interface{}, element interface{}) {
	assert.Contains(th.t, slice, element)
}

// AssertSliceNotContains 断言Slice不包含指定元素
func (th *TestHelper) AssertSliceNotContains(slice interface{}, element interface{}) {
	assert.NotContains(th.t, slice, element)
}

// AssertSliceLength 断言Slice长度
func (th *TestHelper) AssertSliceLength(slice interface{}, expectedLength int) {
	value := reflect.ValueOf(slice)
	assert.Equal(th.t, expectedLength, value.Len())
}

// AssertMapLength 断言Map长度
func (th *TestHelper) AssertMapLength(m map[string]interface{}, expectedLength int) {
	assert.Equal(th.t, expectedLength, len(m))
}

// AssertTimeEqual 断言时间相等（忽略纳秒）
func (th *TestHelper) AssertTimeEqual(expected, actual time.Time) {
	assert.Equal(th.t, expected.Truncate(time.Second), actual.Truncate(time.Second))
}

// AssertTimeNear 断言时间接近（允许误差）
func (th *TestHelper) AssertTimeNear(expected, actual time.Time, tolerance time.Duration) {
	diff := expected.Sub(actual)
	if diff < 0 {
		diff = -diff
	}
	assert.LessOrEqual(th.t, diff, tolerance)
}

// GetFunctionName 获取函数名
func (th *TestHelper) GetFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

// BenchmarkOperation 基准测试操作
func (th *TestHelper) BenchmarkOperation(operation func(), iterations int) time.Duration {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		operation()
	}
	return time.Since(start)
}

// AssertPerformance 断言性能
func (th *TestHelper) AssertPerformance(operation func(), maxDuration time.Duration) {
	duration := th.BenchmarkOperation(operation, 1)
	assert.LessOrEqual(th.t, duration, maxDuration)
}

// CreateMockData 创建模拟数据
func (th *TestHelper) CreateMockData(dataType string, count int) interface{} {
	switch dataType {
	case "users":
		users := make([]map[string]interface{}, count)
		for i := 0; i < count; i++ {
			users[i] = map[string]interface{}{
				"id":       i + 1,
				"name":     fmt.Sprintf("User %d", i+1),
				"email":    fmt.Sprintf("user%d@example.com", i+1),
				"password": "password123",
			}
		}
		return users
	case "posts":
		posts := make([]map[string]interface{}, count)
		for i := 0; i < count; i++ {
			posts[i] = map[string]interface{}{
				"id":      i + 1,
				"title":   fmt.Sprintf("Post %d", i+1),
				"content": fmt.Sprintf("Content for post %d", i+1),
				"user_id": (i % 5) + 1,
			}
		}
		return posts
	default:
		return nil
	}
}

// AssertDatabaseRecord 断言数据库记录
func (th *TestHelper) AssertDatabaseRecord(db database.Connection, tableName string, conditions map[string]interface{}) {
	whereClause := make([]string, 0, len(conditions))
	values := make([]interface{}, 0, len(conditions))

	for column, value := range conditions {
		whereClause = append(whereClause, column+" = ?")
		values = append(values, value)
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", tableName, strings.Join(whereClause, " AND "))
	var count int
	err := db.Query(query, values...).Scan(&count)
	require.NoError(th.t, err)
	assert.Greater(th.t, count, 0)
}

// AssertDatabaseRecordNotExists 断言数据库记录不存在
func (th *TestHelper) AssertDatabaseRecordNotExists(db database.Connection, tableName string, conditions map[string]interface{}) {
	whereClause := make([]string, 0, len(conditions))
	values := make([]interface{}, 0, len(conditions))

	for column, value := range conditions {
		whereClause = append(whereClause, column+" = ?")
		values = append(values, value)
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", tableName, strings.Join(whereClause, " AND "))
	var count int
	err := db.Query(query, values...).Scan(&count)
	require.NoError(th.t, err)
	assert.Equal(th.t, 0, count)
}
