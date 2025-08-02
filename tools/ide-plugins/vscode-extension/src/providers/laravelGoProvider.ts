import * as vscode from 'vscode';

export class LaravelGoProvider implements vscode.HoverProvider {
    provideHover(
        document: vscode.TextDocument,
        position: vscode.Position,
        token: vscode.CancellationToken
    ): vscode.ProviderResult<vscode.Hover> {
        const wordRange = document.getWordRangeAtPosition(position);
        if (!wordRange) {
            return null;
        }

        const word = document.getText(wordRange);
        const hoverInfo = this.getHoverInfo(word);

        if (hoverInfo) {
            return new vscode.Hover(hoverInfo);
        }

        return null;
    }

    private getHoverInfo(word: string): vscode.MarkdownString | null {
        const hoverData: { [key: string]: string } = {
            'app': 'Laravel-Go Application Container\n\nAccess the application container instance to resolve dependencies and manage application state.',
            'config': 'Configuration Manager\n\nAccess configuration values from various sources (files, environment variables, etc.).',
            'Route': 'Route Manager\n\nRegister and manage HTTP routes for your application.',
            'Controller': 'Base Controller\n\nBase controller class that provides common functionality for all controllers.',
            'Model': 'Base Model\n\nBase model class that provides ORM functionality and database interactions.',
            'Middleware': 'Base Middleware\n\nBase middleware class for filtering HTTP requests.',
            'DB': 'Database Manager\n\nExecute database queries and manage database connections.',
            'Make': 'Resolve a dependency from the container\n\n```go\napp.Make("interface") -> concrete\n```',
            'Bind': 'Bind an interface to a concrete implementation\n\n```go\napp.Bind("interface", concrete)\n```',
            'Get': 'Get a configuration value\n\n```go\nconfig.Get("app.name")\n```',
            'Table': 'Begin a fluent query against a database table\n\n```go\nDB.Table("users").Select("*").Get()\n```',
            'Where': 'Add a where clause to the query\n\n```go\nDB.Table("users").Where("id", 1).Get()\n```',
            'Select': 'Select specific columns from the table\n\n```go\nDB.Table("users").Select("id", "name").Get()\n```',
            'OrderBy': 'Add an order by clause to the query\n\n```go\nDB.Table("users").OrderBy("created_at", "desc").Get()\n```',
            'Limit': 'Limit the number of results returned\n\n```go\nDB.Table("users").Limit(10).Get()\n```',
            'Index': 'Display a listing of the resource\n\n```go\nfunc (c *UserController) Index(w http.ResponseWriter, r *http.Request)\n```',
            'Show': 'Display the specified resource\n\n```go\nfunc (c *UserController) Show(w http.ResponseWriter, r *http.Request)\n```',
            'Store': 'Store a newly created resource\n\n```go\nfunc (c *UserController) Store(w http.ResponseWriter, r *http.Request)\n```',
            'Update': 'Update the specified resource\n\n```go\nfunc (c *UserController) Update(w http.ResponseWriter, r *http.Request)\n```',
            'Destroy': 'Remove the specified resource\n\n```go\nfunc (c *UserController) Destroy(w http.ResponseWriter, r *http.Request)\n```',
            'Handle': 'Handle the incoming request\n\n```go\nfunc (m *AuthMiddleware) Handle(w http.ResponseWriter, r *http.Request) bool\n```',
            'Terminate': 'Handle the response after it is sent to the browser\n\n```go\nfunc (m *LogMiddleware) Terminate(w http.ResponseWriter, r *http.Request)\n```',
            'TableName': 'Specify the table name for the model\n\n```go\ntype User struct {\n    Model\n    TableName string\n}\n```',
            'Fillable': 'Specify which fields can be mass assigned\n\n```go\ntype User struct {\n    Model\n    Fillable []string\n}\n```',
            'Hidden': 'Specify which fields should be hidden from JSON output\n\n```go\ntype User struct {\n    Model\n    Hidden []string\n}\n```',
            'Timestamps': 'Enable or disable timestamps\n\n```go\ntype User struct {\n    Model\n    Timestamps bool\n}\n```'
        };

        const info = hoverData[word];
        if (info) {
            return new vscode.MarkdownString(info);
        }

        return null;
    }
} 