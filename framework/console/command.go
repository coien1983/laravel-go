package console

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

// Command 命令接口
type Command interface {
	// GetName 获取命令名称
	GetName() string
	// GetDescription 获取命令描述
	GetDescription() string
	// GetSignature 获取命令签名
	GetSignature() string
	// GetArguments 获取命令参数
	GetArguments() []Argument
	// GetOptions 获取命令选项
	GetOptions() []Option
	// Execute 执行命令
	Execute(input Input) error
}

// Argument 参数定义
type Argument struct {
	Name        string
	Description string
	Required    bool
	Default     interface{}
}

// Option 选项定义
type Option struct {
	Name        string
	ShortName   string
	Description string
	Required    bool
	Default     interface{}
	Type        string // string, bool, int, etc.
}

// Input 输入接口
type Input interface {
	// GetArgument 获取参数值
	GetArgument(name string) interface{}
	// GetOption 获取选项值
	GetOption(name string) interface{}
	// HasOption 检查是否有选项
	HasOption(name string) bool
	// GetArguments 获取所有参数
	GetArguments() map[string]interface{}
	// GetOptions 获取所有选项
	GetOptions() map[string]interface{}
}

// Output 输出接口
type Output interface {
	// Write 写入内容
	Write(content string)
	// WriteLine 写入一行
	WriteLine(content string)
	// Error 输出错误信息
	Error(message string)
	// Success 输出成功信息
	Success(message string)
	// Warning 输出警告信息
	Warning(message string)
	// Info 输出信息
	Info(message string)
	// Table 输出表格
	Table(headers []string, rows [][]string)
}

// Application 命令行应用
type Application struct {
	name        string
	version     string
	commands    map[string]Command
	output      Output
	interactive bool
}

// NewApplication 创建新的命令行应用
func NewApplication(name, version string) *Application {
	return &Application{
		name:     name,
		version:  version,
		commands: make(map[string]Command),
		output:   NewConsoleOutput(),
	}
}

// AddCommand 添加命令
func (app *Application) AddCommand(command Command) {
	app.commands[command.GetName()] = command
}

// GetCommand 获取命令
func (app *Application) GetCommand(name string) (Command, bool) {
	command, exists := app.commands[name]
	return command, exists
}

// GetCommands 获取所有命令
func (app *Application) GetCommands() map[string]Command {
	return app.commands
}

// SetInteractive 设置交互模式
func (app *Application) SetInteractive(interactive bool) {
	app.interactive = interactive
}

// Run 运行应用
func (app *Application) Run(args []string) error {
	if len(args) < 1 {
		return app.showHelp()
	}

	commandName := args[0]

	// 特殊命令处理
	switch commandName {
	case "help", "--help", "-h":
		return app.showHelp()
	case "version", "--version", "-v":
		return app.showVersion()
	case "list":
		return app.listCommands()
	}

	// 获取命令
	command, exists := app.GetCommand(commandName)
	if !exists {
		app.output.Error(fmt.Sprintf("Command '%s' not found", commandName))
		return app.showHelp()
	}

	// 解析输入
	input, err := app.parseInput(args[1:], command)
	if err != nil {
		app.output.Error(fmt.Sprintf("Error parsing input: %v", err))
		return err
	}

	// 执行命令
	return command.Execute(input)
}

// showHelp 显示帮助信息
func (app *Application) showHelp() error {
	app.output.WriteLine(fmt.Sprintf("Usage: %s <command> [options]", app.name))
	app.output.WriteLine("")
	app.output.WriteLine("Available commands:")
	app.output.WriteLine("")

	// 获取所有命令名称并排序
	var commandNames []string
	for name := range app.commands {
		commandNames = append(commandNames, name)
	}
	sort.Strings(commandNames)

	// 显示命令列表
	for _, name := range commandNames {
		command := app.commands[name]
		app.output.WriteLine(fmt.Sprintf("  %-20s %s", name, command.GetDescription()))
	}

	app.output.WriteLine("")
	app.output.WriteLine("For more information about a command, run:")
	app.output.WriteLine(fmt.Sprintf("  %s help <command>", app.name))

	return nil
}

// showVersion 显示版本信息
func (app *Application) showVersion() error {
	app.output.WriteLine(fmt.Sprintf("%s version %s", app.name, app.version))
	return nil
}

// listCommands 列出所有命令
func (app *Application) listCommands() error {
	app.output.WriteLine(fmt.Sprintf("%s version %s", app.name, app.version))
	app.output.WriteLine("")
	app.output.WriteLine("Available commands:")
	app.output.WriteLine("")

	// 获取所有命令名称并排序
	var commandNames []string
	for name := range app.commands {
		commandNames = append(commandNames, name)
	}
	sort.Strings(commandNames)

	// 创建表格输出
	var headers []string
	var rows [][]string

	headers = []string{"Command", "Description", "Signature"}

	for _, name := range commandNames {
		command := app.commands[name]
		rows = append(rows, []string{
			name,
			command.GetDescription(),
			command.GetSignature(),
		})
	}

	app.output.Table(headers, rows)
	return nil
}

// parseInput 解析输入
func (app *Application) parseInput(args []string, command Command) (Input, error) {
	// 创建标志集
	flagSet := flag.NewFlagSet(command.GetName(), flag.ExitOnError)

	// 添加选项
	options := make(map[string]interface{})
	for _, opt := range command.GetOptions() {
		switch opt.Type {
		case "string":
			var value string
			defaultValue := ""
			if opt.Default != nil {
				defaultValue = opt.Default.(string)
			}
			if opt.ShortName != "" {
				flagSet.StringVar(&value, opt.ShortName, defaultValue, opt.Description)
			}
			flagSet.StringVar(&value, opt.Name, defaultValue, opt.Description)
			options[opt.Name] = &value
		case "bool":
			var value bool
			defaultValue := false
			if opt.Default != nil {
				defaultValue = opt.Default.(bool)
			}
			if opt.ShortName != "" {
				flagSet.BoolVar(&value, opt.ShortName, defaultValue, opt.Description)
			}
			flagSet.BoolVar(&value, opt.Name, defaultValue, opt.Description)
			options[opt.Name] = &value
		case "int":
			var value int
			defaultValue := 0
			if opt.Default != nil {
				defaultValue = opt.Default.(int)
			}
			if opt.ShortName != "" {
				flagSet.IntVar(&value, opt.ShortName, defaultValue, opt.Description)
			}
			flagSet.IntVar(&value, opt.Name, defaultValue, opt.Description)
			options[opt.Name] = &value
		}
	}

	// 解析标志
	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	// 获取参数
	arguments := make(map[string]interface{})
	args = flagSet.Args()
	for i, arg := range command.GetArguments() {
		if i < len(args) {
			arguments[arg.Name] = args[i]
		} else if arg.Required {
			return nil, fmt.Errorf("required argument '%s' not provided", arg.Name)
		} else if arg.Default != nil {
			arguments[arg.Name] = arg.Default
		}
	}

	// 获取选项值
	optionValues := make(map[string]interface{})
	for name, ptr := range options {
		switch v := ptr.(type) {
		case *string:
			optionValues[name] = *v
		case *bool:
			optionValues[name] = *v
		case *int:
			optionValues[name] = *v
		}
	}

	return &ConsoleInput{
		arguments: arguments,
		options:   optionValues,
	}, nil
}

// ConsoleInput 控制台输入实现
type ConsoleInput struct {
	arguments map[string]interface{}
	options   map[string]interface{}
}

// GetArgument 获取参数值
func (input *ConsoleInput) GetArgument(name string) interface{} {
	return input.arguments[name]
}

// GetOption 获取选项值
func (input *ConsoleInput) GetOption(name string) interface{} {
	return input.options[name]
}

// HasOption 检查是否有选项
func (input *ConsoleInput) HasOption(name string) bool {
	_, exists := input.options[name]
	return exists
}

// GetArguments 获取所有参数
func (input *ConsoleInput) GetArguments() map[string]interface{} {
	return input.arguments
}

// GetOptions 获取所有选项
func (input *ConsoleInput) GetOptions() map[string]interface{} {
	return input.options
}

// ConsoleOutput 控制台输出实现
type ConsoleOutput struct {
	writer *tabwriter.Writer
}

// NewConsoleOutput 创建新的控制台输出
func NewConsoleOutput() *ConsoleOutput {
	return &ConsoleOutput{
		writer: tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0),
	}
}

// Write 写入内容
func (output *ConsoleOutput) Write(content string) {
	fmt.Print(content)
}

// WriteLine 写入一行
func (output *ConsoleOutput) WriteLine(content string) {
	fmt.Println(content)
}

// Error 输出错误信息
func (output *ConsoleOutput) Error(message string) {
	fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", message)
}

// Success 输出成功信息
func (output *ConsoleOutput) Success(message string) {
	fmt.Printf("\033[32m%s\033[0m\n", message)
}

// Warning 输出警告信息
func (output *ConsoleOutput) Warning(message string) {
	fmt.Printf("\033[33m%s\033[0m\n", message)
}

// Info 输出信息
func (output *ConsoleOutput) Info(message string) {
	fmt.Printf("\033[34m%s\033[0m\n", message)
}

// Table 输出表格
func (output *ConsoleOutput) Table(headers []string, rows [][]string) {
	// 写入表头
	for _, header := range headers {
		fmt.Fprintf(output.writer, "%s\t", header)
	}
	fmt.Fprintln(output.writer)

	// 写入分隔线
	for range headers {
		fmt.Fprintf(output.writer, "----\t")
	}
	fmt.Fprintln(output.writer)

	// 写入数据行
	for _, row := range rows {
		for _, cell := range row {
			fmt.Fprintf(output.writer, "%s\t", cell)
		}
		fmt.Fprintln(output.writer)
	}

	output.writer.Flush()
}
