# 命令行工具指南

## 📖 概述

Laravel-Go Framework 提供了强大的命令行工具系统，包括 Artisan 命令、命令生成器、交互式命令和任务调度等功能。

## 🚀 快速开始

### 1. 基本命令

```go
// 创建命令
type HelloCommand struct {
    console.Command
}

func NewHelloCommand() *HelloCommand {
    return &HelloCommand{
        Command: console.Command{
            Name:        "hello",
            Description: "Say hello to the world",
            Arguments: []console.Argument{
                {
                    Name:        "name",
                    Description: "Your name",
                    Required:    false,
                    Default:     "World",
                },
            },
            Options: []console.Option{
                {
                    Name:        "greeting",
                    ShortName:   "g",
                    Description: "Greeting message",
                    Default:     "Hello",
                },
            },
        },
    }
}

// 执行命令
func (c *HelloCommand) Handle() error {
    name := c.Argument("name")
    greeting := c.Option("greeting")

    c.Info("%s, %s!", greeting, name)

    return nil
}

// 注册命令
func RegisterCommands() {
    app := console.NewApplication("Laravel-Go", "1.0.0")

    app.AddCommand(NewHelloCommand())
    app.AddCommand(NewMakeControllerCommand())
    app.AddCommand(NewMakeModelCommand())
    app.AddCommand(NewMigrateCommand())
    app.AddCommand(NewSeedCommand())

    app.Run()
}

// 主函数
func main() {
    RegisterCommands()
}
```

### 2. 运行命令

```bash
# 运行命令
go run cmd/artisan/main.go hello
go run cmd/artisan/main.go hello John
go run cmd/artisan/main.go hello John --greeting="Hi"

# 查看帮助
go run cmd/artisan/main.go --help
go run cmd/artisan/main.go hello --help
```

## 🔧 命令类型

### 1. 基础命令

```go
// 基础命令
type BasicCommand struct {
    console.Command
}

func NewBasicCommand() *BasicCommand {
    return &BasicCommand{
        Command: console.Command{
            Name:        "basic",
            Description: "A basic command example",
        },
    }
}

func (c *BasicCommand) Handle() error {
    c.Info("This is a basic command")
    c.Comment("This is a comment")
    c.Error("This is an error message")
    c.Warning("This is a warning message")

    return nil
}
```

### 2. 带参数的命令

```go
// 带参数的命令
type UserCommand struct {
    console.Command
}

func NewUserCommand() *UserCommand {
    return &UserCommand{
        Command: console.Command{
            Name:        "user",
            Description: "Manage users",
            Arguments: []console.Argument{
                {
                    Name:        "action",
                    Description: "Action to perform (create, update, delete)",
                    Required:    true,
                },
                {
                    Name:        "name",
                    Description: "User name",
                    Required:    false,
                },
                {
                    Name:        "email",
                    Description: "User email",
                    Required:    false,
                },
            },
        },
    }
}

func (c *UserCommand) Handle() error {
    action := c.Argument("action")
    name := c.Argument("name")
    email := c.Argument("email")

    switch action {
    case "create":
        return c.createUser(name, email)
    case "update":
        return c.updateUser(name, email)
    case "delete":
        return c.deleteUser(name)
    default:
        return fmt.Errorf("unknown action: %s", action)
    }
}

func (c *UserCommand) createUser(name, email string) error {
    if name == "" || email == "" {
        return errors.New("name and email are required for create action")
    }

    c.Info("Creating user: %s (%s)", name, email)

    // 创建用户的逻辑
    user := &User{
        Name:  name,
        Email: email,
    }

    if err := db.Create(user).Error; err != nil {
        return err
    }

    c.Info("User created successfully with ID: %d", user.ID)
    return nil
}

func (c *UserCommand) updateUser(name, email string) error {
    c.Info("Updating user: %s", name)
    // 更新用户的逻辑
    return nil
}

func (c *UserCommand) deleteUser(name string) error {
    c.Info("Deleting user: %s", name)
    // 删除用户的逻辑
    return nil
}
```

### 3. 带选项的命令

```go
// 带选项的命令
type DatabaseCommand struct {
    console.Command
}

func NewDatabaseCommand() *DatabaseCommand {
    return &DatabaseCommand{
        Command: console.Command{
            Name:        "db",
            Description: "Database operations",
            Arguments: []console.Argument{
                {
                    Name:        "operation",
                    Description: "Operation to perform (migrate, seed, backup)",
                    Required:    true,
                },
            },
            Options: []console.Option{
                {
                    Name:        "connection",
                    ShortName:   "c",
                    Description: "Database connection",
                    Default:     "default",
                },
                {
                    Name:        "force",
                    ShortName:   "f",
                    Description: "Force operation",
                    Default:     false,
                },
                {
                    Name:        "verbose",
                    ShortName:   "v",
                    Description: "Verbose output",
                    Default:     false,
                },
            },
        },
    }
}

func (c *DatabaseCommand) Handle() error {
    operation := c.Argument("operation")
    connection := c.Option("connection")
    force := c.OptionBool("force")
    verbose := c.OptionBool("verbose")

    if verbose {
        c.Info("Operation: %s", operation)
        c.Info("Connection: %s", connection)
        c.Info("Force: %v", force)
    }

    switch operation {
    case "migrate":
        return c.migrate(connection, force)
    case "seed":
        return c.seed(connection, force)
    case "backup":
        return c.backup(connection)
    default:
        return fmt.Errorf("unknown operation: %s", operation)
    }
}

func (c *DatabaseCommand) migrate(connection string, force bool) error {
    c.Info("Running migrations on connection: %s", connection)

    if force {
        c.Warning("Force mode enabled - this may cause data loss")
    }

    // 执行迁移逻辑
    return nil
}

func (c *DatabaseCommand) seed(connection string, force bool) error {
    c.Info("Running seeders on connection: %s", connection)

    // 执行数据填充逻辑
    return nil
}

func (c *DatabaseCommand) backup(connection string) error {
    c.Info("Creating backup for connection: %s", connection)

    // 执行备份逻辑
    return nil
}
```

## 🎯 交互式命令

### 1. 交互式输入

```go
// 交互式命令
type InteractiveCommand struct {
    console.Command
}

func NewInteractiveCommand() *InteractiveCommand {
    return &InteractiveCommand{
        Command: console.Command{
            Name:        "interactive",
            Description: "Interactive command example",
        },
    }
}

func (c *InteractiveCommand) Handle() error {
    // 询问用户输入
    name, err := c.Ask("What is your name?")
    if err != nil {
        return err
    }

    // 询问密码（隐藏输入）
    password, err := c.AskHidden("Enter your password:")
    if err != nil {
        return err
    }

    // 确认操作
    confirmed, err := c.Confirm("Do you want to proceed?")
    if err != nil {
        return err
    }

    if !confirmed {
        c.Info("Operation cancelled")
        return nil
    }

    // 选择选项
    choice, err := c.Choice("Select your favorite color:", []string{"Red", "Green", "Blue"})
    if err != nil {
        return err
    }

    c.Info("Hello, %s!", name)
    c.Info("Your favorite color is: %s", choice)

    return nil
}
```

### 2. 进度条

```go
// 带进度条的命令
type ProgressCommand struct {
    console.Command
}

func NewProgressCommand() *ProgressCommand {
    return &ProgressCommand{
        Command: console.Command{
            Name:        "progress",
            Description: "Command with progress bar",
        },
    }
}

func (c *ProgressCommand) Handle() error {
    c.Info("Starting operation...")

    // 创建进度条
    progress := c.ProgressBar(100, "Processing")

    for i := 0; i <= 100; i++ {
        // 模拟工作
        time.Sleep(time.Millisecond * 50)

        // 更新进度
        progress.Set(i)

        if i%20 == 0 {
            c.Info("Completed %d%%", i)
        }
    }

    progress.Finish()
    c.Info("Operation completed!")

    return nil
}
```

### 3. 表格输出

```go
// 表格输出命令
type TableCommand struct {
    console.Command
}

func NewTableCommand() *TableCommand {
    return &TableCommand{
        Command: console.Command{
            Name:        "table",
            Description: "Command with table output",
        },
    }
}

func (c *TableCommand) Handle() error {
    // 创建表格
    table := c.Table([]string{"ID", "Name", "Email", "Status"})

    // 添加数据
    table.AddRow([]string{"1", "John Doe", "john@example.com", "Active"})
    table.AddRow([]string{"2", "Jane Smith", "jane@example.com", "Inactive"})
    table.AddRow([]string{"3", "Bob Johnson", "bob@example.com", "Active"})

    // 渲染表格
    table.Render()

    return nil
}
```

## 🔨 命令生成器

### 1. 控制器生成器

```go
// 控制器生成器命令
type MakeControllerCommand struct {
    console.Command
}

func NewMakeControllerCommand() *MakeControllerCommand {
    return &MakeControllerCommand{
        Command: console.Command{
            Name:        "make:controller",
            Description: "Create a new controller",
            Arguments: []console.Argument{
                {
                    Name:        "name",
                    Description: "Controller name",
                    Required:    true,
                },
            },
            Options: []console.Option{
                {
                    Name:        "resource",
                    ShortName:   "r",
                    Description: "Generate resource controller",
                    Default:     false,
                },
                {
                    Name:        "api",
                    ShortName:   "a",
                    Description: "Generate API controller",
                    Default:     false,
                },
            },
        },
    }
}

func (c *MakeControllerCommand) Handle() error {
    name := c.Argument("name")
    isResource := c.OptionBool("resource")
    isAPI := c.OptionBool("api")

    // 生成控制器文件名
    fileName := fmt.Sprintf("app/Http/Controllers/%sController.go", name)

    // 生成控制器内容
    content := c.generateControllerContent(name, isResource, isAPI)

    // 写入文件
    if err := c.writeFile(fileName, content); err != nil {
        return err
    }

    c.Info("Controller created successfully: %s", fileName)

    return nil
}

func (c *MakeControllerCommand) generateControllerContent(name string, isResource, isAPI bool) string {
    var content strings.Builder

    content.WriteString("package controllers\n\n")
    content.WriteString("import (\n")
    content.WriteString("    \"laravel-go/framework/http\"\n")
    content.WriteString(")\n\n")

    content.WriteString(fmt.Sprintf("type %sController struct {\n", name))
    content.WriteString("    http.Controller\n")
    content.WriteString("}\n\n")

    if isResource {
        content.WriteString(c.generateResourceMethods(name))
    } else {
        content.WriteString(c.generateBasicMethods(name))
    }

    return content.String()
}

func (c *MakeControllerCommand) generateResourceMethods(name string) string {
    var content strings.Builder

    methods := []string{"Index", "Show", "Store", "Update", "Delete"}

    for _, method := range methods {
        content.WriteString(fmt.Sprintf("func (c *%sController) %s(request http.Request) http.Response {\n", name, method))
        content.WriteString("    // TODO: Implement method\n")
        content.WriteString("    return c.Json(map[string]string{\n")
        content.WriteString(fmt.Sprintf("        \"message\": \"%s method called\",\n", method))
        content.WriteString("    })\n")
        content.WriteString("}\n\n")
    }

    return content.String()
}

func (c *MakeControllerCommand) writeFile(fileName, content string) error {
    // 确保目录存在
    dir := filepath.Dir(fileName)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    // 写入文件
    return ioutil.WriteFile(fileName, []byte(content), 0644)
}
```

### 2. 模型生成器

```go
// 模型生成器命令
type MakeModelCommand struct {
    console.Command
}

func NewMakeModelCommand() *MakeModelCommand {
    return &MakeModelCommand{
        Command: console.Command{
            Name:        "make:model",
            Description: "Create a new model",
            Arguments: []console.Argument{
                {
                    Name:        "name",
                    Description: "Model name",
                    Required:    true,
                },
            },
            Options: []console.Option{
                {
                    Name:        "migration",
                    ShortName:   "m",
                    Description: "Create migration file",
                    Default:     false,
                },
                {
                    Name:        "factory",
                    ShortName:   "f",
                    Description: "Create factory file",
                    Default:     false,
                },
            },
        },
    }
}

func (c *MakeModelCommand) Handle() error {
    name := c.Argument("name")
    createMigration := c.OptionBool("migration")
    createFactory := c.OptionBool("factory")

    // 生成模型文件
    modelFileName := fmt.Sprintf("app/Models/%s.go", name)
    modelContent := c.generateModelContent(name)

    if err := c.writeFile(modelFileName, modelContent); err != nil {
        return err
    }

    c.Info("Model created successfully: %s", modelFileName)

    // 生成迁移文件
    if createMigration {
        migrationFileName := fmt.Sprintf("database/migrations/create_%s_table.go", strings.ToLower(name))
        migrationContent := c.generateMigrationContent(name)

        if err := c.writeFile(migrationFileName, migrationContent); err != nil {
            return err
        }

        c.Info("Migration created successfully: %s", migrationFileName)
    }

    // 生成工厂文件
    if createFactory {
        factoryFileName := fmt.Sprintf("database/factories/%sFactory.go", name)
        factoryContent := c.generateFactoryContent(name)

        if err := c.writeFile(factoryFileName, factoryContent); err != nil {
            return err
        }

        c.Info("Factory created successfully: %s", factoryFileName)
    }

    return nil
}

func (c *MakeModelCommand) generateModelContent(name string) string {
    var content strings.Builder

    content.WriteString("package models\n\n")
    content.WriteString("import (\n")
    content.WriteString("    \"laravel-go/framework/database\"\n")
    content.WriteString("    \"time\"\n")
    content.WriteString(")\n\n")

    content.WriteString(fmt.Sprintf("type %s struct {\n", name))
    content.WriteString("    database.Model\n")
    content.WriteString("    ID        uint      `json:\"id\" gorm:\"primaryKey\"`\n")
    content.WriteString("    CreatedAt time.Time `json:\"created_at\"`\n")
    content.WriteString("    UpdatedAt time.Time `json:\"updated_at\"`\n")
    content.WriteString("}\n")

    return content.String()
}
```

## 📅 任务调度

### 1. 调度器命令

```go
// 调度器命令
type ScheduleCommand struct {
    console.Command
}

func NewScheduleCommand() *ScheduleCommand {
    return &ScheduleCommand{
        Command: console.Command{
            Name:        "schedule",
            Description: "Run scheduled tasks",
        },
    }
}

func (c *ScheduleCommand) Handle() error {
    c.Info("Running scheduled tasks...")

    // 获取所有调度任务
    tasks := c.getScheduledTasks()

    for _, task := range tasks {
        if c.shouldRunTask(task) {
            c.Info("Running task: %s", task.Name)

            if err := c.runTask(task); err != nil {
                c.Error("Task failed: %s - %v", task.Name, err)
            } else {
                c.Info("Task completed: %s", task.Name)
            }
        }
    }

    c.Info("All scheduled tasks completed")
    return nil
}

type ScheduledTask struct {
    Name     string
    Schedule string
    Command  string
    LastRun  time.Time
}

func (c *ScheduleCommand) getScheduledTasks() []ScheduledTask {
    return []ScheduledTask{
        {
            Name:     "Daily Backup",
            Schedule: "0 2 * * *", // 每天凌晨2点
            Command:  "db:backup",
            LastRun:  time.Now().Add(-24 * time.Hour),
        },
        {
            Name:     "Clean Logs",
            Schedule: "0 3 * * *", // 每天凌晨3点
            Command:  "logs:clean",
            LastRun:  time.Now().Add(-24 * time.Hour),
        },
    }
}

func (c *ScheduleCommand) shouldRunTask(task ScheduledTask) bool {
    // 解析 Cron 表达式并检查是否应该运行
    // 这里简化处理，实际应该使用 Cron 解析库
    return time.Since(task.LastRun) > time.Hour*24
}

func (c *ScheduleCommand) runTask(task ScheduledTask) error {
    // 执行任务命令
    // 这里应该调用相应的命令处理器
    return nil
}
```

### 2. 队列工作进程

```go
// 队列工作进程命令
type QueueWorkCommand struct {
    console.Command
}

func NewQueueWorkCommand() *QueueWorkCommand {
    return &QueueWorkCommand{
        Command: console.Command{
            Name:        "queue:work",
            Description: "Start queue worker",
            Options: []console.Option{
                {
                    Name:        "queue",
                    ShortName:   "q",
                    Description: "Queue name",
                    Default:     "default",
                },
                {
                    Name:        "workers",
                    ShortName:   "w",
                    Description: "Number of workers",
                    Default:     1,
                },
                {
                    Name:        "timeout",
                    ShortName:   "t",
                    Description: "Job timeout in seconds",
                    Default:     60,
                },
            },
        },
    }
}

func (c *QueueWorkCommand) Handle() error {
    queueName := c.Option("queue")
    workers := c.OptionInt("workers")
    timeout := c.OptionInt("timeout")

    c.Info("Starting queue worker for queue: %s", queueName)
    c.Info("Workers: %d, Timeout: %d seconds", workers, timeout)

    // 启动队列工作进程
    worker := queue.NewWorker(queueName, workers)
    worker.SetTimeout(time.Duration(timeout) * time.Second)

    // 设置信号处理
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        c.Info("Shutting down queue worker...")
        worker.Stop()
    }()

    return worker.Start()
}
```

## 🔧 高级功能

### 1. 命令组

```go
// 命令组
type UserGroupCommand struct {
    console.Command
}

func NewUserGroupCommand() *UserGroupCommand {
    return &UserGroupCommand{
        Command: console.Command{
            Name:        "user",
            Description: "User management commands",
            Commands: []console.Command{
                *NewUserCreateCommand(),
                *NewUserUpdateCommand(),
                *NewUserDeleteCommand(),
            },
        },
    }
}

// 子命令
type UserCreateCommand struct {
    console.Command
}

func NewUserCreateCommand() *UserCreateCommand {
    return &UserCreateCommand{
        Command: console.Command{
            Name:        "create",
            Description: "Create a new user",
            Arguments: []console.Argument{
                {
                    Name:        "name",
                    Description: "User name",
                    Required:    true,
                },
                {
                    Name:        "email",
                    Description: "User email",
                    Required:    true,
                },
            },
        },
    }
}

func (c *UserCreateCommand) Handle() error {
    name := c.Argument("name")
    email := c.Argument("email")

    c.Info("Creating user: %s (%s)", name, email)

    // 创建用户的逻辑
    return nil
}
```

### 2. 命令中间件

```go
// 命令中间件
type CommandMiddleware interface {
    Before(command console.Command) error
    After(command console.Command, err error) error
}

// 日志中间件
type LoggingMiddleware struct{}

func (m *LoggingMiddleware) Before(command console.Command) error {
    log.Printf("Executing command: %s", command.Name)
    return nil
}

func (m *LoggingMiddleware) After(command console.Command, err error) error {
    if err != nil {
        log.Printf("Command failed: %s - %v", command.Name, err)
    } else {
        log.Printf("Command completed: %s", command.Name)
    }
    return err
}

// 使用中间件
func RegisterCommands() {
    app := console.NewApplication("Laravel-Go", "1.0.0")

    // 添加中间件
    app.Use(&LoggingMiddleware{})

    // 添加命令
    app.AddCommand(NewHelloCommand())
    app.AddCommand(NewUserGroupCommand())

    app.Run()
}
```

## 📚 总结

Laravel-Go Framework 的命令行工具系统提供了：

1. **基础命令**: 简单的命令执行
2. **参数和选项**: 灵活的命令参数处理
3. **交互式命令**: 用户交互和进度显示
4. **命令生成器**: 自动生成代码文件
5. **任务调度**: 定时任务执行
6. **队列工作**: 后台任务处理
7. **命令组**: 组织相关命令
8. **中间件**: 命令执行前后的处理

通过合理使用命令行工具系统，可以构建强大的开发和管理工具。
