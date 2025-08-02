package main

import (
	"laravel-go/framework/console"
	"os"
)

func main() {
	// 创建命令行应用
	app := console.NewApplication("laravel-go", "1.0.0")
	output := console.NewConsoleOutput()
	generator := console.NewGenerator(output)

	// 注册代码生成命令
	app.AddCommand(console.NewMakeControllerCommand(generator))
	app.AddCommand(console.NewMakeModelCommand(generator))
	app.AddCommand(console.NewMakeMiddlewareCommand(generator))
	app.AddCommand(console.NewMakeMigrationCommand(generator))
	app.AddCommand(console.NewMakeTestCommand(generator))

	// 注册项目管理命令
	app.AddCommand(console.NewInitCommand(output))
	app.AddCommand(console.NewClearCacheCommand(output))
	app.AddCommand(console.NewRouteListCommand(output))

	// 运行应用
	if err := app.Run(os.Args); err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
}
