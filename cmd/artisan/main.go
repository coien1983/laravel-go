package main

import (
	"laravel-go/framework/console"
	"os"
)

func main() {
	// 创建应用
	app := console.NewApplication("Laravel-Go Artisan", "1.0.0")
	output := console.NewConsoleOutput()
	generator := console.NewGenerator(output)

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
	// 部署配置命令
	// =============================================================================
	app.AddCommand(console.NewMakeDockerCommand(generator))
	app.AddCommand(console.NewMakeK8sCommand(generator))

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
