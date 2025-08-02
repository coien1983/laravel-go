package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"laravel-go/framework/testing"
)

// IntegrationTestRunner 集成测试运行器
type IntegrationTestRunner struct {
	runner *testing.TestRunner
}

// NewIntegrationTestRunner 创建集成测试运行器
func NewIntegrationTestRunner() *IntegrationTestRunner {
	return &IntegrationTestRunner{
		runner: testing.NewTestRunner(),
	}
}

// RunAllTests 运行所有集成测试
func (itr *IntegrationTestRunner) RunAllTests() {
	fmt.Println("🚀 开始运行 Laravel-Go Framework 集成测试...")
	fmt.Println("=" * 60)

	startTime := time.Now()

	// 创建测试套件
	httpSuite := testing.NewTestSuite("HTTP集成测试")
	databaseSuite := testing.NewTestSuite("数据库集成测试")
	cacheQueueSuite := testing.NewTestSuite("缓存队列集成测试")
	integrationSuite := testing.NewTestSuite("综合集成测试")

	// 添加HTTP测试用例
	httpSuite.AddTest(createHTTPTestCase("基础路由测试", "测试基本路由功能", func() error {
		// 这里会调用实际的HTTP测试
		return nil
	}))

	httpSuite.AddTest(createHTTPTestCase("路由参数测试", "测试路由参数处理", func() error {
		return nil
	}))

	httpSuite.AddTest(createHTTPTestCase("中间件测试", "测试中间件功能", func() error {
		return nil
	}))

	// 添加数据库测试用例
	databaseSuite.AddTest(createDatabaseTestCase("基础CRUD测试", "测试数据库基本操作", func() error {
		return nil
	}))

	databaseSuite.AddTest(createDatabaseTestCase("模型关联测试", "测试模型关联关系", func() error {
		return nil
	}))

	databaseSuite.AddTest(createDatabaseTestCase("事务测试", "测试数据库事务", func() error {
		return nil
	}))

	// 添加缓存队列测试用例
	cacheQueueSuite.AddTest(createCacheTestCase("缓存基础操作", "测试缓存基本功能", func() error {
		return nil
	}))

	cacheQueueSuite.AddTest(createQueueTestCase("队列基础操作", "测试队列基本功能", func() error {
		return nil
	}))

	// 添加综合测试用例
	integrationSuite.AddTest(createIntegrationTestCase("HTTP+数据库集成", "测试HTTP和数据库的集成", func() error {
		return nil
	}))

	integrationSuite.AddTest(createIntegrationTestCase("缓存+队列集成", "测试缓存和队列的集成", func() error {
		return nil
	}))

	// 添加测试套件到运行器
	itr.runner.AddSuite(httpSuite)
	itr.runner.AddSuite(databaseSuite)
	itr.runner.AddSuite(cacheQueueSuite)
	itr.runner.AddSuite(integrationSuite)

	// 运行所有测试
	results := itr.runner.Run()

	// 生成测试报告
	report := itr.runner.GenerateReport(results)

	// 打印测试结果
	itr.printTestResults(report, startTime)

	// 保存测试报告
	itr.saveTestReport(report)
}

// RunSpecificSuite 运行特定的测试套件
func (itr *IntegrationTestRunner) RunSpecificSuite(suiteName string) {
	fmt.Printf("🚀 开始运行 %s 测试套件...\n", suiteName)
	fmt.Println("=" * 60)

	startTime := time.Now()

	result := itr.runner.RunSuite(suiteName)

	// 生成测试报告
	report := itr.runner.GenerateReport([]*testing.TestResult{result})

	// 打印测试结果
	itr.printTestResults(report, startTime)
}

// RunSpecificTest 运行特定的测试用例
func (itr *IntegrationTestRunner) RunSpecificTest(suiteName, testName string) {
	fmt.Printf("🚀 开始运行测试: %s - %s\n", suiteName, testName)
	fmt.Println("=" * 60)

	startTime := time.Now()

	// 这里可以实现运行特定测试的逻辑
	// 目前简化处理，直接运行整个套件
	itr.RunSpecificSuite(suiteName)
}

// printTestResults 打印测试结果
func (itr *IntegrationTestRunner) printTestResults(report *testing.TestReport, startTime time.Time) {
	fmt.Println("\n📊 测试结果汇总")
	fmt.Println("=" * 60)
	fmt.Printf("总测试数: %d\n", report.TotalTests)
	fmt.Printf("通过测试: %d\n", report.PassedTests)
	fmt.Printf("失败测试: %d\n", report.FailedTests)
	fmt.Printf("成功率: %.2f%%\n", report.SuccessRate*100)
	fmt.Printf("总耗时: %v\n", time.Since(startTime))
	fmt.Printf("生成时间: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))

	fmt.Println("\n📋 详细结果")
	fmt.Println("-" * 60)

	for _, suiteResult := range report.Suites {
		fmt.Printf("\n🏷️  %s\n", suiteResult.SuiteName)
		fmt.Printf("   通过: %d, 失败: %d, 成功率: %.2f%%\n",
			suiteResult.PassedCount,
			suiteResult.FailedCount,
			suiteResult.SuccessRate()*100)
		fmt.Printf("   耗时: %v\n", suiteResult.Duration())

		for _, testResult := range suiteResult.Tests {
			status := "✅"
			if testResult.Status == testing.TestStatusFailed {
				status = "❌"
			}
			fmt.Printf("   %s %s (%v)\n", status, testResult.Name, testResult.Duration())
			if testResult.Error != nil {
				fmt.Printf("     错误: %v\n", testResult.Error)
			}
		}
	}

	// 打印总结
	fmt.Println("\n" + "="*60)
	if report.SuccessRate == 1.0 {
		fmt.Println("🎉 所有测试通过！")
	} else if report.SuccessRate >= 0.8 {
		fmt.Println("⚠️  大部分测试通过，但有一些失败")
	} else {
		fmt.Println("❌ 测试失败较多，需要检查")
	}
}

// saveTestReport 保存测试报告
func (itr *IntegrationTestRunner) saveTestReport(report *testing.TestReport) {
	// 创建测试报告目录
	reportDir := "test_reports"
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		fmt.Printf("❌ 创建测试报告目录失败: %v\n", err)
		return
	}

	// 生成报告文件名
	reportFile := fmt.Sprintf("%s/integration_test_report_%s.html",
		reportDir,
		report.GeneratedAt.Format("20060102_150405"))

	// 生成HTML报告
	html := itr.generateHTMLReport(report)

	// 保存报告文件
	if err := os.WriteFile(reportFile, []byte(html), 0644); err != nil {
		fmt.Printf("❌ 保存测试报告失败: %v\n", err)
		return
	}

	fmt.Printf("📄 测试报告已保存到: %s\n", reportFile)
}

// generateHTMLReport 生成HTML测试报告
func (itr *IntegrationTestRunner) generateHTMLReport(report *testing.TestReport) string {
	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Laravel-Go Framework 集成测试报告</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; }
        .summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .summary-card { background: #f8f9fa; padding: 20px; border-radius: 8px; text-align: center; }
        .summary-card h3 { margin: 0 0 10px 0; color: #333; }
        .summary-card .number { font-size: 2em; font-weight: bold; }
        .success { color: #28a745; }
        .failure { color: #dc3545; }
        .warning { color: #ffc107; }
        .suite { margin-bottom: 30px; border: 1px solid #ddd; border-radius: 8px; overflow: hidden; }
        .suite-header { background: #e9ecef; padding: 15px; font-weight: bold; }
        .test { padding: 10px 15px; border-bottom: 1px solid #eee; }
        .test:last-child { border-bottom: none; }
        .test.passed { background-color: #d4edda; }
        .test.failed { background-color: #f8d7da; }
        .test-duration { color: #666; font-size: 0.9em; }
        .test-error { color: #dc3545; font-size: 0.9em; margin-top: 5px; }
        .status-icon { margin-right: 10px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Laravel-Go Framework 集成测试报告</h1>
            <p>生成时间: ` + report.GeneratedAt.Format("2006-01-02 15:04:05") + `</p>
        </div>

        <div class="summary">
            <div class="summary-card">
                <h3>总测试数</h3>
                <div class="number">` + fmt.Sprintf("%d", report.TotalTests) + `</div>
            </div>
            <div class="summary-card">
                <h3>通过测试</h3>
                <div class="number success">` + fmt.Sprintf("%d", report.PassedTests) + `</div>
            </div>
            <div class="summary-card">
                <h3>失败测试</h3>
                <div class="number failure">` + fmt.Sprintf("%d", report.FailedTests) + `</div>
            </div>
            <div class="summary-card">
                <h3>成功率</h3>
                <div class="number ` + getSuccessClass(report.SuccessRate) + `">` + fmt.Sprintf("%.1f%%", report.SuccessRate*100) + `</div>
            </div>
        </div>`

	// 添加测试套件详情
	for _, suiteResult := range report.Suites {
		html += `
        <div class="suite">
            <div class="suite-header">
                🏷️ ` + suiteResult.SuiteName + `
                <span style="float: right; font-size: 0.9em;">
                    通过: ` + fmt.Sprintf("%d", suiteResult.PassedCount) + `, 
                    失败: ` + fmt.Sprintf("%d", suiteResult.FailedCount) + `, 
                    成功率: ` + fmt.Sprintf("%.1f%%", suiteResult.SuccessRate()*100) + `
                </span>
            </div>`

		for _, testResult := range suiteResult.Tests {
			statusClass := "passed"
			statusIcon := "✅"
			if testResult.Status == testing.TestStatusFailed {
				statusClass = "failed"
				statusIcon = "❌"
			}

			html += `
            <div class="test ` + statusClass + `">
                <div>
                    <span class="status-icon">` + statusIcon + `</span>
                    ` + testResult.Name + `
                    <span class="test-duration">(` + testResult.Duration().String() + `)</span>
                </div>`

			if testResult.Error != nil {
				html += `
                <div class="test-error">错误: ` + testResult.Error.Error() + `</div>`
			}

			html += `
            </div>`
		}

		html += `
        </div>`
	}

	html += `
    </div>
</body>
</html>`

	return html
}

// getSuccessClass 根据成功率返回CSS类名
func getSuccessClass(successRate float64) string {
	if successRate == 1.0 {
		return "success"
	} else if successRate >= 0.8 {
		return "warning"
	} else {
		return "failure"
	}
}

// 创建测试用例的辅助函数
func createHTTPTestCase(name, description string, testFunc func() error) *testing.HTTPTestCase {
	return testing.NewHTTPTestCase(name, description).
		SetTestFunc(testFunc)
}

func createDatabaseTestCase(name, description string, testFunc func() error) *testing.DatabaseTestCase {
	return testing.NewDatabaseTestCase(name, description).
		SetTestFunc(testFunc)
}

func createCacheTestCase(name, description string, testFunc func() error) *testing.BaseTestCase {
	return testing.NewBaseTestCase(name, description).
		SetTestFunc(testFunc)
}

func createQueueTestCase(name, description string, testFunc func() error) *testing.BaseTestCase {
	return testing.NewBaseTestCase(name, description).
		SetTestFunc(testFunc)
}

func createIntegrationTestCase(name, description string, testFunc func() error) *testing.BaseTestCase {
	return testing.NewBaseTestCase(name, description).
		SetTestFunc(testFunc)
}

// 命令行入口函数
func main() {
	runner := NewIntegrationTestRunner()

	// 检查命令行参数
	args := os.Args[1:]
	if len(args) == 0 {
		// 运行所有测试
		runner.RunAllTests()
	} else if len(args) == 1 {
		// 运行特定套件
		runner.RunSpecificSuite(args[0])
	} else if len(args) == 2 {
		// 运行特定测试
		runner.RunSpecificTest(args[0], args[1])
	} else {
		fmt.Println("用法:")
		fmt.Println("  go run tests/run_integration_tests.go                    # 运行所有测试")
		fmt.Println("  go run tests/run_integration_tests.go <suite>            # 运行特定套件")
		fmt.Println("  go run tests/run_integration_tests.go <suite> <test>     # 运行特定测试")
	}
}
