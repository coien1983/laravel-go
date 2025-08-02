package main

import (
	"os"

	"github.com/coien1983/laravel-go/framework/console"
)

func main() {
	// 创建应用
	app := console.NewApplication("Laravel-Go Artisan", "1.0.0")
	output := console.NewConsoleOutput()
	generator := console.NewGenerator(output)
	goZeroGenerator := console.NewGoZeroGenerator(output)

	// =============================================================================
	// 项目初始化命令
	// =============================================================================
	app.AddCommand(console.NewInitCommand(output))

	// =============================================================================
	// 代码生成命令
	// =============================================================================
	app.AddCommand(console.NewMakeControllerCommand(generator))
	app.AddCommand(console.NewMakeModelCommand(generator))
	app.AddCommand(console.NewMakeMiddlewareCommand(generator))
	app.AddCommand(console.NewMakeMigrationCommand(generator))
	app.AddCommand(console.NewMakeTestCommand(generator))

	// =============================================================================
	// go-zero 增强命令 (类似 goctl 功能)
	// =============================================================================
	app.AddCommand(console.NewGoZeroMakeRpcCommand(goZeroGenerator))
	app.AddCommand(console.NewGoZeroMakeApiCommand(goZeroGenerator))

	// =============================================================================
	// 项目维护命令
	// =============================================================================
	app.AddCommand(console.NewClearCacheCommand(output))
	app.AddCommand(console.NewRouteListCommand(output))

	// =============================================================================
	// 快速生成命令
	// =============================================================================
	app.AddCommand(console.NewMakeApiCommand(generator))
	app.AddCommand(console.NewMakeCrudCommand(generator))

	// =============================================================================
	// 模块管理命令
	// =============================================================================
	app.AddCommand(console.NewAddModuleCommand(generator))
	app.AddCommand(console.NewAddServiceCommand(generator))
	app.AddCommand(console.NewAddRepositoryCommand(generator))
	app.AddCommand(console.NewAddValidatorCommand(generator))
	app.AddCommand(console.NewAddEventCommand(generator))

	// =============================================================================
	// 项目信息命令
	// =============================================================================
	app.AddCommand(console.NewProjectInfoCommand(output))
	app.AddCommand(console.NewVersionCommand(output))

	// 运行应用
	if err := app.Run(os.Args[1:]); err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
}
