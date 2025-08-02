package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type Scaffold struct {
	ProjectName string
	ProjectRoot string
	Templates   map[string]*template.Template
}

type ProjectData struct {
	Name        string
	ModuleName  string
	Description string
	Author      string
	Version     string
}

func NewScaffold(projectName, projectRoot string) *Scaffold {
	scaffold := &Scaffold{
		ProjectName: projectName,
		ProjectRoot: projectRoot,
		Templates:   make(map[string]*template.Template),
	}
	scaffold.loadTemplates()
	return scaffold
}

func (s *Scaffold) loadTemplates() {
	// go.mod template
	goModTmpl := `module {{.ModuleName}}

go 1.21

require (
	laravel-go/framework v0.1.0
)

require (
	github.com/go-sql-driver/mysql v1.7.1
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.17
)
`
	s.Templates["go.mod"] = template.Must(template.New("go.mod").Parse(goModTmpl))

	// main.go template
	mainTmpl := `package main

import (
	"log"
	"laravel-go/framework"
	"laravel-go/framework/http"
	"laravel-go/framework/console"
)

func main() {
	// Create application
	app := framework.NewApplication()

	// Register routes
	registerRoutes(app)

	// Start server
	log.Fatal(app.Run(":8080"))
}

func registerRoutes(app *framework.Application) {
	// Register your routes here
	app.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Laravel-Go Framework!"))
	})
}
`
	s.Templates["main.go"] = template.Must(template.New("main.go").Parse(mainTmpl))

	// README.md template
	readmeTmpl := `# {{.Name}}

{{.Description}}

## Installation

1. Clone the repository
2. Install dependencies: ` + "`go mod tidy`" + `
3. Run the application: ` + "`go run main.go`" + `

## Development

### Running the Application
` + "```bash" + `
go run main.go
` + "```" + `

### Running Tests
` + "```bash" + `
go test ./...
` + "```" + `

### Code Generation
` + "```bash" + `
go run cmd/artisan/main.go make:controller UserController
go run cmd/artisan/main.go make:model User
go run cmd/artisan/main.go make:middleware AuthMiddleware
` + "```" + `

## Project Structure

```
{{.Name}}/
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
└── README.md           # This file
```

## License

This project is licensed under the MIT License.`
	s.Templates["README.md"] = template.Must(template.New("README.md").Parse(readmeTmpl))

	// .gitignore template
	gitignoreTmpl := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Application specific
storage/logs/*
!storage/logs/.gitkeep
storage/cache/*
!storage/cache/.gitkeep
.env
.env.local
.env.production

# Database
*.db
*.sqlite
*.sqlite3

# Build artifacts
build/
dist/
`
	s.Templates[".gitignore"] = template.Must(template.New(".gitignore").Parse(gitignoreTmpl))

	// .env template
	envTmpl := `APP_NAME={{.Name}}
APP_ENV=local
APP_DEBUG=true
APP_KEY=base64:your-secret-key-here

DB_CONNECTION=sqlite
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE={{.Name}}
DB_USERNAME=root
DB_PASSWORD=

CACHE_DRIVER=file
QUEUE_CONNECTION=sync
SESSION_DRIVER=file
SESSION_LIFETIME=120`
	s.Templates[".env"] = template.Must(template.New(".env").Parse(envTmpl))
}

func (s *Scaffold) CreateProject() error {
	// Create project directory
	if err := os.MkdirAll(s.ProjectRoot, 0755); err != nil {
		return err
	}

	// Create directory structure
	directories := []string{
		"app/controllers",
		"app/models",
		"app/middleware",
		"app/services",
		"config",
		"database/migrations",
		"database/seeders",
		"routes",
		"resources/views",
		"resources/assets",
		"storage/logs",
		"storage/cache",
		"tests",
		"cmd/artisan",
	}

	for _, dir := range directories {
		path := filepath.Join(s.ProjectRoot, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	// Create files
	data := ProjectData{
		Name:        s.ProjectName,
		ModuleName:  s.ProjectName,
		Description: "A Laravel-Go Framework application",
		Author:      "Your Name",
		Version:     "1.0.0",
	}

	files := map[string]string{
		"go.mod":     "go.mod",
		"main.go":    "main.go",
		"README.md":  "README.md",
		".gitignore": ".gitignore",
		".env":       ".env",
	}

	for templateName, filename := range files {
		filepath := filepath.Join(s.ProjectRoot, filename)
		if err := s.writeTemplate(templateName, filepath, data); err != nil {
			return err
		}
	}

	// Create placeholder files
	placeholderFiles := map[string]string{
		"storage/logs/.gitkeep":   "",
		"storage/cache/.gitkeep":  "",
		"app/controllers/.gitkeep": "",
		"app/models/.gitkeep":     "",
		"app/middleware/.gitkeep": "",
		"app/services/.gitkeep":   "",
		"config/.gitkeep":         "",
		"database/migrations/.gitkeep": "",
		"database/seeders/.gitkeep":    "",
		"routes/.gitkeep":         "",
		"resources/views/.gitkeep": "",
		"resources/assets/.gitkeep": "",
		"tests/.gitkeep":          "",
	}

	for filepath, content := range placeholderFiles {
		fullPath := filepath.Join(s.ProjectRoot, filepath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	// Copy artisan command
	if err := s.copyArtisanCommand(); err != nil {
		return err
	}

	return nil
}

func (s *Scaffold) copyArtisanCommand() error {
	artisanContent := `package main

import (
	"laravel-go/framework/console"
	"os"
)

func main() {
	// Create application
	app := console.NewApplication("Laravel-Go Artisan", "1.0.0")
	output := console.NewConsoleOutput()
	generator := console.NewGenerator(output)

	// Register all commands
	app.AddCommand(console.NewMakeControllerCommand(generator))
	app.AddCommand(console.NewMakeModelCommand(generator))
	app.AddCommand(console.NewMakeMiddlewareCommand(generator))
	app.AddCommand(console.NewMakeMigrationCommand(generator))
	app.AddCommand(console.NewMakeTestCommand(generator))
	app.AddCommand(console.NewInitCommand(output))
	app.AddCommand(console.NewClearCacheCommand(output))
	app.AddCommand(console.NewRouteListCommand(output))

	// Run application
	if err := app.Run(os.Args[1:]); err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
}
`

	artisanPath := filepath.Join(s.ProjectRoot, "cmd", "artisan", "main.go")
	return os.WriteFile(artisanPath, []byte(artisanContent), 0644)
}

func (s *Scaffold) writeTemplate(templateName, filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, exists := s.Templates[templateName]
	if !exists {
		return fmt.Errorf("template %s not found", templateName)
	}

	return tmpl.Execute(file, data)
}

func main() {
	var (
		projectName = flag.String("name", "", "Project name")
		projectRoot = flag.String("path", "", "Project root directory (defaults to project name)")
	)
	flag.Parse()

	if *projectName == "" {
		fmt.Println("Usage: project-scaffold -name <project-name> [-path <project-path>]")
		os.Exit(1)
	}

	if *projectRoot == "" {
		*projectRoot = *projectName
	}

	scaffold := NewScaffold(*projectName, *projectRoot)

	if err := scaffold.CreateProject(); err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Project '%s' created successfully at '%s'\n", *projectName, *projectRoot)
	fmt.Println("\nNext steps:")
	fmt.Printf("1. cd %s\n", *projectRoot)
	fmt.Println("2. go mod tidy")
	fmt.Println("3. go run main.go")
} 