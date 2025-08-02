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

	// 注册所有命令
	app.AddCommand(console.NewMakeControllerCommand(generator))
	app.AddCommand(console.NewMakeModelCommand(generator))
	app.AddCommand(console.NewMakeMiddlewareCommand(generator))
	app.AddCommand(console.NewMakeMigrationCommand(generator))
	app.AddCommand(console.NewMakeTestCommand(generator))
	app.AddCommand(console.NewInitCommand(output))
	app.AddCommand(console.NewClearCacheCommand(output))
	app.AddCommand(console.NewRouteListCommand(output))

	// 注册新的部署命令
	app.AddCommand(console.NewMakeDockerCommand(generator))
	app.AddCommand(console.NewMakeK8sCommand(generator))

	// 运行应用
	if err := app.Run(os.Args[1:]); err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
}
