# Laravel-Go Framework GoLand Plugin

This plugin provides comprehensive IDE support for Laravel-Go Framework development in GoLand.

## Features

### 🚀 Code Generation
- **Generate Controller**: Create new controllers with standard CRUD methods
- **Generate Model**: Create new models with ORM functionality
- **Generate Middleware**: Create new middleware classes
- **Generate Migration**: Create new database migrations
- **Generate Test**: Create new test files

### 🔧 Development Tools
- **Run Migrations**: Execute database migrations
- **Run Seeders**: Execute database seeders
- **Run Tests**: Execute test suites
- **Serve Application**: Start the development server

### 💡 Code Completion
- Framework-specific code completions
- Controller method suggestions
- Model field suggestions
- Database query builder completions
- Route registration completions

### 📝 Live Templates
- Controller templates
- Model templates
- Middleware templates
- Route definitions
- Database queries
- Container bindings
- Configuration access

### 🐛 Debug Support
- Pre-configured debug configurations for Laravel-Go applications
- Support for debugging tests
- Environment variable management

### 🔍 Code Inspection
- Framework-specific code inspections
- Best practices validation
- Common error detection

## Installation

### From Source
1. Clone this repository
2. Navigate to the plugin directory: `cd tools/ide-plugins/goland-plugin`
3. Build the plugin: `./gradlew buildPlugin`
4. Install the plugin in GoLand:
   - Go to `File` → `Settings` → `Plugins`
   - Click `Install Plugin from Disk`
   - Select the generated JAR file from `build/distributions/`

### From JetBrains Marketplace
*Coming soon*

## Usage

### Project Detection
The plugin automatically detects Laravel-Go projects by looking for the `cmd/artisan/main.go` file. Once detected, all Laravel-Go features become available.

### Code Generation
Right-click in the Project view and select:
- `Generate Controller` - Create a new controller
- `Generate Model` - Create a new model
- `Generate Middleware` - Create a new middleware
- `Generate Migration` - Create a new migration
- `Generate Test` - Create a new test

### Artisan Commands
Access Laravel-Go Artisan commands through the Tools menu:
- `Tools` → `Run Migrations`
- `Tools` → `Run Seeders`
- `Tools` → `Run Tests`
- `Tools` → `Serve Application`

### Live Templates
Use the following abbreviations to trigger live templates:
- `lg-controller` - Create a new controller
- `lg-model` - Create a new model
- `lg-middleware` - Create a new middleware
- `lg-route` - Create a new route
- `lg-route-group` - Create a route group
- `lg-db-query` - Create a database query
- `lg-bind` - Bind to container
- `lg-make` - Resolve from container
- `lg-config` - Get configuration value
- `lg-test` - Create a test
- `lg-migration` - Create a migration

### Debugging
Use the provided debug configurations:
- **Laravel-Go Application**: Debug the main application
- **Laravel-Go Tests**: Debug test execution
- **Laravel-Go Artisan Command**: Debug specific artisan commands

## Development

### Building
```bash
./gradlew buildPlugin
```

### Testing
```bash
./gradlew test
```

### Running in Development Mode
```bash
./gradlew runIde
```

### Packaging
```bash
./gradlew buildPlugin
```

## Configuration

The plugin automatically configures itself based on the project structure. No manual configuration is required.

## Project Structure

```
src/main/kotlin/com/laravelgo/plugin/
├── LaravelGoPlugin.kt              # Main plugin class
├── LaravelGoProjectService.kt      # Project service
├── actions/                        # Action classes
│   ├── GenerateControllerAction.kt
│   ├── GenerateModelAction.kt
│   ├── GenerateMiddlewareAction.kt
│   ├── GenerateMigrationAction.kt
│   ├── GenerateTestAction.kt
│   ├── RunMigrateAction.kt
│   ├── RunSeedAction.kt
│   ├── RunTestAction.kt
│   └── RunServeAction.kt
├── dialogs/                        # Dialog classes
│   ├── GenerateControllerDialog.kt
│   ├── GenerateModelDialog.kt
│   ├── GenerateMiddlewareDialog.kt
│   ├── GenerateMigrationDialog.kt
│   └── GenerateTestDialog.kt
├── completion/                     # Code completion
│   └── LaravelGoCompletionContributor.kt
├── inspection/                     # Code inspection
│   └── LaravelGoInspection.kt
├── debug/                         # Debug configuration
│   ├── LaravelGoDebugConfigurationType.kt
│   └── LaravelGoRunConfigurationProducer.kt
├── project/                       # Project structure
│   └── LaravelGoProjectStructureDetector.kt
├── toolwindow/                    # Tool window
│   └── LaravelGoToolWindowFactory.kt
└── templates/                     # Live templates
    ├── LaravelGoTemplateContext.kt
    └── LaravelGoTemplatesProvider.kt
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This plugin is part of the Laravel-Go Framework project and follows the same license terms. 