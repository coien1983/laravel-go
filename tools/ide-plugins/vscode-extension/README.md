# Laravel-Go Framework VS Code Extension

This extension provides comprehensive IDE support for Laravel-Go Framework development in Visual Studio Code.

## Features

### 🚀 Code Generation Commands
- **Generate Controller**: Create new controllers with standard CRUD methods
- **Generate Model**: Create new models with ORM functionality
- **Generate Middleware**: Create new middleware classes
- **Generate Migration**: Create new database migrations
- **Generate Test**: Create new test files

### 🔧 Development Commands
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

### 📚 Hover Information
- Detailed documentation for framework components
- Code examples and usage patterns
- Parameter descriptions and return types

### 🐛 Debug Support
- Pre-configured debug configurations for Laravel-Go applications
- Support for debugging tests
- Environment variable management

### 📝 Code Snippets
- Controller templates
- Model templates
- Middleware templates
- Route definitions
- Database queries
- Container bindings
- Configuration access

## Installation

1. Clone this repository
2. Navigate to the extension directory: `cd tools/ide-plugins/vscode-extension`
3. Install dependencies: `npm install`
4. Build the extension: `npm run compile`
5. Package the extension: `vsce package`
6. Install the VSIX file in VS Code

## Configuration

The extension can be configured through VS Code settings:

```json
{
  "laravelGo.artisanPath": "cmd/artisan",
  "laravelGo.projectRoot": ".",
  "laravelGo.enableCodeCompletion": true,
  "laravelGo.enableDebugging": true
}
```

## Usage

### Command Palette
Access all Laravel-Go commands through the command palette (`Ctrl+Shift+P` or `Cmd+Shift+P`):
- `Laravel-Go: Generate Controller`
- `Laravel-Go: Generate Model`
- `Lavel-Go: Generate Middleware`
- `Laravel-Go: Generate Migration`
- `Laravel-Go: Generate Test`
- `Laravel-Go: Run Migrations`
- `Laravel-Go: Run Seeders`
- `Laravel-Go: Run Tests`
- `Laravel-Go: Serve Application`

### Code Snippets
Use the following prefixes to trigger code snippets:
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
- **Launch Laravel-Go Application**: Debug the main application
- **Debug Laravel-Go Tests**: Debug test execution
- **Debug Laravel-Go Artisan Command**: Debug specific artisan commands

## Development

### Building
```bash
npm install
npm run compile
```

### Testing
```bash
npm run test
```

### Packaging
```bash
npm install -g vsce
vsce package
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This extension is part of the Laravel-Go Framework project and follows the same license terms. 