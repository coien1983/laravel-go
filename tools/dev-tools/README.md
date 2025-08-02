# Laravel-Go Development Tools

This directory contains various development tools to enhance the Laravel-Go Framework development experience.

## Tools Overview

### 1. Code Generator (`code-generator/`)

A powerful code generation tool that creates Laravel-Go framework components from templates.

**Features:**
- Generate controllers with CRUD methods
- Generate models with ORM functionality
- Generate middleware classes
- Generate database migrations
- Generate test files

**Usage:**
```bash
cd tools/dev-tools/code-generator
go run main.go -command controller -name UserController
go run main.go -command model -name User -table users
go run main.go -command middleware -name AuthMiddleware
go run main.go -command migration -name CreateUsersTable -table users
go run main.go -command test -name UserTest -package tests
```

**Options:**
- `-command`: Type of component to generate (controller, model, middleware, migration, test)
- `-name`: Name of the component
- `-table`: Table name (for models and migrations)
- `-package`: Package name (for tests)
- `-project`: Project root directory (default: current directory)

### 2. Project Scaffold (`project-scaffold/`)

A project scaffolding tool that creates new Laravel-Go projects with the proper directory structure and initial files.

**Features:**
- Creates complete project structure
- Generates initial configuration files
- Sets up go.mod with proper dependencies
- Creates basic README and documentation
- Includes Artisan command-line tool

**Usage:**
```bash
cd tools/dev-tools/project-scaffold
go run main.go -name my-laravel-go-app
go run main.go -name my-app -path /path/to/project
```

**Options:**
- `-name`: Project name (required)
- `-path`: Project root directory (defaults to project name)

**Generated Structure:**
```
my-laravel-go-app/
├── app/
│   ├── controllers/     # HTTP controllers
│   ├── models/         # Database models
│   ├── middleware/     # HTTP middleware
│   └── services/       # Business logic services
├── config/             # Configuration files
├── database/
│   ├── migrations/     # Database migrations
│   └── seeders/        # Database seeders
├── routes/             # Route definitions
├── resources/
│   ├── views/          # Template views
│   └── assets/         # Static assets
├── storage/            # Application storage
├── tests/              # Test files
├── cmd/
│   └── artisan/        # Artisan command line tool
├── main.go             # Application entry point
├── go.mod              # Go module file
└── README.md           # Project documentation
```

### 3. Performance Analyzer (`performance-analyzer/`)

A real-time performance monitoring tool for Laravel-Go applications.

**Features:**
- Memory usage monitoring
- CPU usage tracking
- HTTP request metrics
- Goroutine monitoring
- Performance profiling
- JSON metrics export

**Usage:**
```bash
cd tools/dev-tools/performance-analyzer
go run main.go -port 8080 -interval 5s -profile
```

**Options:**
- `-port`: HTTP server port (default: 8080)
- `-interval`: Metrics collection interval (default: 5s)
- `-profile`: Enable CPU and memory profiling
- `-output`: Output file for metrics (default: performance_metrics.json)

**Test Endpoints:**
- `GET /` - Success response
- `GET /error` - Error response

**Metrics Collected:**
- Memory allocation and usage
- Garbage collection statistics
- HTTP request/response times
- Success/error rates
- Goroutine count

### 4. Debug Tool (`debug-tool/`)

A comprehensive debugging tool for Laravel-Go applications.

**Features:**
- Runtime information
- Memory statistics
- Goroutine analysis
- Stack traces
- CPU and memory profiling
- Garbage collection control
- Environment variable inspection

**Usage:**
```bash
cd tools/dev-tools/debug-tool
go run main.go -port 6060 -profile
```

**Options:**
- `-port`: Debug server port (default: 6060)
- `-profile`: Start profiling on startup

**Debug Endpoints:**
- `GET /debug` - General debug information
- `GET /debug/memory` - Memory profile download
- `GET /debug/cpu` - CPU profile download
- `GET /debug/goroutines` - Goroutine information
- `GET /debug/stack` - Stack traces
- `GET /debug/gc` - Force garbage collection
- `GET /debug/pprof/*` - Go pprof endpoints

## Installation and Setup

### Prerequisites
- Go 1.21 or later
- Git

### Building All Tools
```bash
cd tools/dev-tools

# Build code generator
cd code-generator
go build -o ../../bin/code-generator

# Build project scaffold
cd ../project-scaffold
go build -o ../../bin/project-scaffold

# Build performance analyzer
cd ../performance-analyzer
go build -o ../../bin/performance-analyzer

# Build debug tool
cd ../debug-tool
go build -o ../../bin/debug-tool
```

### Adding to PATH
Add the tools to your PATH for easy access:
```bash
export PATH=$PATH:/path/to/laravel-go/bin
```

## Integration with IDE

These tools can be integrated with your IDE:

### VS Code
Add to your VS Code settings:
```json
{
  "laravel-go.codeGeneratorPath": "/path/to/laravel-go/bin/code-generator",
  "laravel-go.projectScaffoldPath": "/path/to/laravel-go/bin/project-scaffold",
  "laravel-go.performanceAnalyzerPath": "/path/to/laravel-go/bin/performance-analyzer",
  "laravel-go.debugToolPath": "/path/to/laravel-go/bin/debug-tool"
}
```

### GoLand
Configure external tools in GoLand:
1. Go to `File` → `Settings` → `Tools` → `External Tools`
2. Add each tool with appropriate parameters
3. Assign keyboard shortcuts for quick access

## Best Practices

### Code Generation
- Use descriptive names for generated components
- Review generated code before committing
- Customize templates for your project needs
- Use consistent naming conventions

### Performance Analysis
- Run performance analysis during development
- Monitor memory usage patterns
- Track HTTP response times
- Use profiling for optimization

### Debugging
- Use debug tool during development
- Monitor goroutine leaks
- Check memory usage regularly
- Use profiling for performance issues

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

These tools are part of the Laravel-Go Framework project and follow the same license terms. 