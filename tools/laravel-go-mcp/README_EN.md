# Laravel-Go MCP Service

A comprehensive Model Context Protocol (MCP) service for Laravel-Go framework management, providing complete project lifecycle management capabilities.

## Features

- 🚀 **Project Initialization**: Quickly create new Laravel-Go projects
- 📝 **Code Generation**: Automatically generate controllers, models, services, and other modules
- 🔨 **Project Building**: Automated build and compilation
- 🧪 **Test Execution**: Run unit tests and integration tests
- 🚀 **Project Deployment**: Multi-environment deployment support
- 📊 **Performance Monitoring**: Real-time performance data collection
- 🔍 **Code Analysis**: Code quality checks and statistics
- ⚡ **Performance Optimization**: Automatic performance optimization suggestions

## Quick Start

### 1. Start MCP Server

```bash
cd tools/laravel-go-mcp
go run main.go
```

The server will start on `http://localhost:8080`.

### 2. Use Client Example

```bash
# Run client demonstration
go run client_example.go
```

## API Interfaces

### 1. Initialize Project

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "name": "my-api",
    "description": "My API Project",
    "version": "1.0.0",
    "modules": ["user", "product", "order"],
    "database": "mysql",
    "cache": "redis",
    "queue": "redis"
  }
}
```

### 2. Generate Module

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "generate",
  "params": {
    "type": "api",
    "name": "category"
  }
}
```

### 3. Build Project

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "build",
  "params": {}
}
```

### 4. Run Tests

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "test",
  "params": {}
}
```

### 5. Deploy Project

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "deploy",
  "params": {
    "environment": "production"
  }
}
```

### 6. Performance Monitoring

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "monitor",
  "params": {}
}
```

### 7. Code Analysis

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 7,
  "method": "analyze",
  "params": {}
}
```

### 8. Performance Optimization

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 8,
  "method": "optimize",
  "params": {}
}
```

### 9. Get Project Information

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 9,
  "method": "info",
  "params": {}
}
```

## Response Format

All interfaces return responses in standard JSON-RPC 2.0 format:

### Success Response

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "success": true,
    "message": "Operation successful",
    "data": {}
  }
}
```

### Error Response

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32601,
    "message": "Method not found",
    "data": null
  }
}
```

## Error Codes

| Code   | Description      |
| ------ | ---------------- |
| -32700 | Parse error      |
| -32600 | Invalid request  |
| -32601 | Method not found |
| -32602 | Invalid params   |
| -32603 | Internal error   |
| -32000 | Server error     |

## Project Structure

```
laravel-go-mcp/
├── main.go              # MCP server main file
├── client_example.go    # Client example
└── README.md           # Documentation
```

## Environment Variables

| Variable | Default | Description     |
| -------- | ------- | --------------- |
| MCP_PORT | 8080    | MCP server port |

## Usage Examples

### Go Client

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    client := NewMCPClientExample("http://localhost:8080")

    // Initialize project
    resp, err := client.Initialize(&ClientInitializeRequest{
        Name:        "my-api",
        Description: "My API Project",
        Version:     "1.0.0",
        Modules:     []string{"user", "product"},
        Database:    "mysql",
        Cache:       "redis",
        Queue:       "redis",
    })

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Initialization result: %v\n", resp["result"])
}
```

### cURL Examples

```bash
# Initialize project
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "name": "my-api",
      "description": "My API Project",
      "version": "1.0.0",
      "modules": ["user", "product"],
      "database": "mysql",
      "cache": "redis",
      "queue": "redis"
    }
  }'

# Generate module
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "generate",
    "params": {
      "type": "api",
      "name": "category"
    }
  }'
```

## Development Guide

### Adding New MCP Methods

1. Add new handler method in `main.go`
2. Add corresponding client method in `client_example.go`
3. Update documentation

### Extending Features

- Support more module types
- Add database migration functionality
- Integrate CI/CD workflows
- Add more performance monitoring metrics

## License

MIT License
