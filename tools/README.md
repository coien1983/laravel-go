# Laravel-Go Framework Tools

This directory contains tools and plugins to enhance the Laravel-Go Framework development experience.

## Directory Structure

```
tools/
├── ide-plugins/           # IDE plugins and extensions
│   ├── vscode-extension/  # Visual Studio Code extension
│   └── goland-plugin/     # GoLand plugin
└── dev-tools/             # Development tools
    ├── code-generator/    # Code generation tool
    ├── project-scaffold/  # Project scaffolding tool
    ├── performance-analyzer/ # Performance monitoring tool
    └── debug-tool/        # Debugging tool
```

## IDE Plugins

### Visual Studio Code Extension

A comprehensive VS Code extension that provides:
- Code generation commands
- Code completion for Laravel-Go framework
- Hover information and documentation
- Debug configurations
- Code snippets and templates

**Installation:**
1. Navigate to `tools/ide-plugins/vscode-extension/`
2. Run `npm install`
3. Run `npm run compile`
4. Package with `vsce package`
5. Install the VSIX file in VS Code

**Features:**
- Generate controllers, models, middleware, migrations, and tests
- Run Laravel-Go Artisan commands
- Framework-specific code completions
- Debug support for Laravel-Go applications

### GoLand Plugin

A GoLand plugin that provides:
- Code generation actions
- Live templates
- Code completion
- Debug configurations
- Project structure detection

**Installation:**
1. Navigate to `tools/ide-plugins/goland-plugin/`
2. Run `./gradlew buildPlugin`
3. Install the generated JAR file in GoLand

**Features:**
- Right-click menu actions for code generation
- Live templates for Laravel-Go components
- Framework-specific inspections
- Debug configurations for Laravel-Go applications

## Development Tools

### Code Generator

A command-line tool for generating Laravel-Go framework components.

**Usage:**
```bash
cd tools/dev-tools/code-generator
go run main.go -command controller -name UserController
```

**Supported Commands:**
- `controller` - Generate HTTP controllers
- `model` - Generate database models
- `middleware` - Generate HTTP middleware
- `migration` - Generate database migrations
- `test` - Generate test files

### Project Scaffold

A tool for creating new Laravel-Go projects with proper structure.

**Usage:**
```bash
cd tools/dev-tools/project-scaffold
go run main.go -name my-laravel-go-app
```

**Features:**
- Creates complete project directory structure
- Generates initial configuration files
- Sets up go.mod with dependencies
- Includes Artisan command-line tool

### Performance Analyzer

A real-time performance monitoring tool.

**Usage:**
```bash
cd tools/dev-tools/performance-analyzer
go run main.go -port 8080 -interval 5s
```

**Features:**
- Memory usage monitoring
- CPU usage tracking
- HTTP request metrics
- Performance profiling
- JSON metrics export

### Debug Tool

A comprehensive debugging tool for Laravel-Go applications.

**Usage:**
```bash
cd tools/dev-tools/debug-tool
go run main.go -port 6060
```

**Features:**
- Runtime information
- Memory statistics
- Goroutine analysis
- Stack traces
- CPU and memory profiling

## Quick Start

### 1. Install IDE Plugin

Choose your preferred IDE and install the corresponding plugin:

**VS Code:**
```bash
cd tools/ide-plugins/vscode-extension
npm install && npm run compile
```

**GoLand:**
```bash
cd tools/ide-plugins/goland-plugin
./gradlew buildPlugin
```

### 2. Build Development Tools

```bash
cd tools/dev-tools

# Build all tools
for tool in code-generator project-scaffold performance-analyzer debug-tool; do
    cd $tool
    go build -o ../../bin/$tool
    cd ..
done
```

### 3. Create a New Project

```bash
# Using project scaffold
./bin/project-scaffold -name my-app

# Navigate to project
cd my-app

# Generate components
./bin/code-generator -command controller -name UserController
./bin/code-generator -command model -name User
```

### 4. Monitor Performance

```bash
# Start performance analyzer
./bin/performance-analyzer -port 8080

# Start debug tool
./bin/debug-tool -port 6060
```

## Configuration

### IDE Settings

**VS Code:**
```json
{
  "laravelGo.artisanPath": "cmd/artisan",
  "laravelGo.projectRoot": ".",
  "laravelGo.enableCodeCompletion": true,
  "laravelGo.enableDebugging": true
}
```

**GoLand:**
The plugin automatically detects Laravel-Go projects and configures itself.

### Environment Variables

```bash
# Development tools paths
export LARAVEL_GO_CODE_GENERATOR_PATH="/path/to/bin/code-generator"
export LARAVEL_GO_PROJECT_SCAFFOLD_PATH="/path/to/bin/project-scaffold"
export LARAVEL_GO_PERFORMANCE_ANALYZER_PATH="/path/to/bin/performance-analyzer"
export LARAVEL_GO_DEBUG_TOOL_PATH="/path/to/bin/debug-tool"
```

## Integration Examples

### VS Code Tasks

Add to `.vscode/tasks.json`:
```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Generate Controller",
      "type": "shell",
      "command": "code-generator",
      "args": ["-command", "controller", "-name", "${input:controllerName}"],
      "group": "build"
    },
    {
      "label": "Generate Model",
      "type": "shell",
      "command": "code-generator",
      "args": ["-command", "model", "-name", "${input:modelName}"],
      "group": "build"
    }
  ]
}
```

### GoLand External Tools

Configure external tools in GoLand:
1. `File` → `Settings` → `Tools` → `External Tools`
2. Add tools with appropriate parameters
3. Assign keyboard shortcuts

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Documentation

- [IDE Plugins Documentation](ide-plugins/README.md)
- [Development Tools Documentation](dev-tools/README.md)
- [VS Code Extension Documentation](ide-plugins/vscode-extension/README.md)
- [GoLand Plugin Documentation](ide-plugins/goland-plugin/README.md)

## License

These tools are part of the Laravel-Go Framework project and follow the same license terms. 