package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Generator struct {
	ProjectRoot string
	Templates   map[string]*template.Template
}

type ControllerData struct {
	Name        string
	PackageName string
	Methods     []string
}

type ModelData struct {
	Name        string
	PackageName string
	TableName   string
	Fields      []Field
}

type Field struct {
	Name     string
	Type     string
	Tag      string
	Required bool
}

type MiddlewareData struct {
	Name        string
	PackageName string
}

type MigrationData struct {
	Name        string
	PackageName string
	TableName   string
	Fields      []Field
}

func NewGenerator(projectRoot string) *Generator {
	gen := &Generator{
		ProjectRoot: projectRoot,
		Templates:   make(map[string]*template.Template),
	}
	gen.loadTemplates()
	return gen
}

func (g *Generator) loadTemplates() {
	// Controller template
	controllerTmpl := `package controllers

import (
	"net/http"
	"github.com/coien1983/laravel-go/framework/http"
)

type {{.Name}}Controller struct {
	http.Controller
}

func New{{.Name}}Controller() *{{.Name}}Controller {
	return &{{.Name}}Controller{}
}

func (c *{{.Name}}Controller) Index(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement index method
}

func (c *{{.Name}}Controller) Show(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement show method
}

func (c *{{.Name}}Controller) Store(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement store method
}

func (c *{{.Name}}Controller) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement update method
}

func (c *{{.Name}}Controller) Destroy(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement destroy method
}
`
	g.Templates["controller"] = template.Must(template.New("controller").Parse(controllerTmpl))

	// Model template
	modelTmpl := `package models

import (
	"time"
	"github.com/coien1983/laravel-go/framework/database"
)

type {{.Name}} struct {
	database.Model
	ID        uint      ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`{{.Tag}}`" + `
{{end}}}

func (m *{{.Name}}) TableName() string {
	return "{{.TableName}}"
}

func (m *{{.Name}}) Fillable() []string {
	return []string{
{{range .Fields}}		"{{.Name}}",
{{end}}	}
}

func (m *{{.Name}}) Hidden() []string {
	return []string{
		"password",
	}
}
`
	g.Templates["model"] = template.Must(template.New("model").Parse(modelTmpl))

	// Middleware template
	middlewareTmpl := `package middleware

import (
	"net/http"
	"github.com/coien1983/laravel-go/framework/http"
)

type {{.Name}}Middleware struct {
	http.Middleware
}

func New{{.Name}}Middleware() *{{.Name}}Middleware {
	return &{{.Name}}Middleware{}
}

func (m *{{.Name}}Middleware) Handle(w http.ResponseWriter, r *http.Request) bool {
	// TODO: Implement middleware logic
	return true
}

func (m *{{.Name}}Middleware) Terminate(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement termination logic
}
`
	g.Templates["middleware"] = template.Must(template.New("middleware").Parse(middlewareTmpl))

	// Migration template
	migrationTmpl := `package migrations

import (
	"github.com/coien1983/laravel-go/framework/database"
)

type {{.Name}}Migration struct {
	database.Migration
}

func (m *{{.Name}}Migration) Up() error {
	return DB.Schema.Create("{{.TableName}}", func(table *database.Blueprint) {
		table.ID()
{{range .Fields}}		table.{{.Type}}("{{.Name}}")
{{end}}		table.Timestamps()
	})
}

func (m *{{.Name}}Migration) Down() error {
	return DB.Schema.Drop("{{.TableName}}")
}
`
	g.Templates["migration"] = template.Must(template.New("migration").Parse(migrationTmpl))

	// Test template
	testTmpl := `package {{.PackageName}}

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/coien1983/laravel-go/framework/testing"
)

func Test{{.Name}}(t *testing.T) {
	// TODO: Implement test
}
`
	g.Templates["test"] = template.Must(template.New("test").Parse(testTmpl))
}

func (g *Generator) GenerateController(name string) error {
	data := ControllerData{
		Name:        name,
		PackageName: "controllers",
	}

	dir := filepath.Join(g.ProjectRoot, "app", "controllers")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filename := filepath.Join(dir, strings.ToLower(name)+"_controller.go")
	return g.writeTemplate("controller", filename, data)
}

func (g *Generator) GenerateModel(name, tableName string, fields []Field) error {
	data := ModelData{
		Name:        name,
		PackageName: "models",
		TableName:   tableName,
		Fields:      fields,
	}

	dir := filepath.Join(g.ProjectRoot, "app", "models")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filename := filepath.Join(dir, strings.ToLower(name)+".go")
	return g.writeTemplate("model", filename, data)
}

func (g *Generator) GenerateMiddleware(name string) error {
	data := MiddlewareData{
		Name:        name,
		PackageName: "middleware",
	}

	dir := filepath.Join(g.ProjectRoot, "app", "middleware")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filename := filepath.Join(dir, strings.ToLower(name)+"_middleware.go")
	return g.writeTemplate("middleware", filename, data)
}

func (g *Generator) GenerateMigration(name, tableName string, fields []Field) error {
	data := MigrationData{
		Name:        name,
		PackageName: "migrations",
		TableName:   tableName,
		Fields:      fields,
	}

	dir := filepath.Join(g.ProjectRoot, "database", "migrations")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filename := filepath.Join(dir, strings.ToLower(name)+"_migration.go")
	return g.writeTemplate("migration", filename, data)
}

func (g *Generator) GenerateTest(name, packageName string) error {
	data := struct {
		Name        string
		PackageName string
	}{
		Name:        name,
		PackageName: packageName,
	}

	dir := filepath.Join(g.ProjectRoot, "tests")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filename := filepath.Join(dir, strings.ToLower(name)+"_test.go")
	return g.writeTemplate("test", filename, data)
}

func (g *Generator) writeTemplate(templateName, filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, exists := g.Templates[templateName]
	if !exists {
		return fmt.Errorf("template %s not found", templateName)
	}

	return tmpl.Execute(file, data)
}

func main() {
	var (
		projectRoot = flag.String("project", ".", "Project root directory")
		command     = flag.String("command", "", "Command to execute (controller, model, middleware, migration, test)")
		name        = flag.String("name", "", "Name for the generated file")
		tableName   = flag.String("table", "", "Table name (for models and migrations)")
		packageName = flag.String("package", "", "Package name (for tests)")
	)
	flag.Parse()

	if *command == "" || *name == "" {
		fmt.Println("Usage: code-generator -command <command> -name <name> [options]")
		fmt.Println("Commands: controller, model, middleware, migration, test")
		os.Exit(1)
	}

	generator := NewGenerator(*projectRoot)

	switch *command {
	case "controller":
		if err := generator.GenerateController(*name); err != nil {
			fmt.Printf("Error generating controller: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Controller %s generated successfully\n", *name)

	case "model":
		if *tableName == "" {
			*tableName = strings.ToLower(*name) + "s"
		}
		fields := []Field{
			{Name: "Name", Type: "string", Tag: `json:"name"`, Required: true},
			{Name: "Email", Type: "string", Tag: `json:"email"`, Required: true},
		}
		if err := generator.GenerateModel(*name, *tableName, fields); err != nil {
			fmt.Printf("Error generating model: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Model %s generated successfully\n", *name)

	case "middleware":
		if err := generator.GenerateMiddleware(*name); err != nil {
			fmt.Printf("Error generating middleware: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Middleware %s generated successfully\n", *name)

	case "migration":
		if *tableName == "" {
			*tableName = strings.ToLower(*name) + "s"
		}
		fields := []Field{
			{Name: "name", Type: "String", Tag: "", Required: true},
			{Name: "email", Type: "String", Tag: "", Required: true},
		}
		if err := generator.GenerateMigration(*name, *tableName, fields); err != nil {
			fmt.Printf("Error generating migration: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Migration %s generated successfully\n", *name)

	case "test":
		if *packageName == "" {
			*packageName = "tests"
		}
		if err := generator.GenerateTest(*name, *packageName); err != nil {
			fmt.Printf("Error generating test: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Test %s generated successfully\n", *name)

	default:
		fmt.Printf("Unknown command: %s\n", *command)
		os.Exit(1)
	}
}
