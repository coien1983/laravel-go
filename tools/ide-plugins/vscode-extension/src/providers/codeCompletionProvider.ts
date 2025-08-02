import * as vscode from 'vscode';

export class CodeCompletionProvider implements vscode.CompletionItemProvider {
    provideCompletionItems(
        document: vscode.TextDocument,
        position: vscode.Position,
        token: vscode.CancellationToken,
        context: vscode.CompletionContext
    ): vscode.ProviderResult<vscode.CompletionItem[] | vscode.CompletionList<vscode.CompletionItem>> {
        const items: vscode.CompletionItem[] = [];

        // Laravel-Go Framework specific completions
        items.push(...this.getFrameworkCompletions());
        items.push(...this.getControllerCompletions());
        items.push(...this.getModelCompletions());
        items.push(...this.getMiddlewareCompletions());
        items.push(...this.getRouteCompletions());
        items.push(...this.getDatabaseCompletions());

        return items;
    }

    private getFrameworkCompletions(): vscode.CompletionItem[] {
        const items: vscode.CompletionItem[] = [];

        // Container completions
        const containerItem = new vscode.CompletionItem('app', vscode.CompletionItemKind.Class);
        containerItem.detail = 'Laravel-Go Application Container';
        containerItem.documentation = 'Access the application container instance';
        items.push(containerItem);

        const makeItem = new vscode.CompletionItem('app.Make', vscode.CompletionItemKind.Method);
        makeItem.detail = 'Resolve a dependency from the container';
        makeItem.documentation = 'Resolve a dependency from the application container';
        items.push(makeItem);

        const bindItem = new vscode.CompletionItem('app.Bind', vscode.CompletionItemKind.Method);
        bindItem.detail = 'Bind an interface to a concrete implementation';
        bindItem.documentation = 'Bind an interface to a concrete implementation in the container';
        items.push(bindItem);

        // Config completions
        const configItem = new vscode.CompletionItem('config', vscode.CompletionItemKind.Class);
        configItem.detail = 'Configuration Manager';
        configItem.documentation = 'Access configuration values';
        items.push(configItem);

        const getItem = new vscode.CompletionItem('config.Get', vscode.CompletionItemKind.Method);
        getItem.detail = 'Get a configuration value';
        getItem.documentation = 'Get a configuration value by key';
        items.push(getItem);

        return items;
    }

    private getControllerCompletions(): vscode.CompletionItem[] {
        const items: vscode.CompletionItem[] = [];

        const controllerItem = new vscode.CompletionItem('Controller', vscode.CompletionItemKind.Class);
        controllerItem.detail = 'Base Controller';
        controllerItem.documentation = 'Base controller class for Laravel-Go';
        items.push(controllerItem);

        const indexItem = new vscode.CompletionItem('Index', vscode.CompletionItemKind.Method);
        indexItem.detail = 'Index method';
        indexItem.documentation = 'Display a listing of the resource';
        items.push(indexItem);

        const showItem = new vscode.CompletionItem('Show', vscode.CompletionItemKind.Method);
        showItem.detail = 'Show method';
        showItem.documentation = 'Display the specified resource';
        items.push(showItem);

        const storeItem = new vscode.CompletionItem('Store', vscode.CompletionItemKind.Method);
        storeItem.detail = 'Store method';
        storeItem.documentation = 'Store a newly created resource';
        items.push(storeItem);

        const updateItem = new vscode.CompletionItem('Update', vscode.CompletionItemKind.Method);
        updateItem.detail = 'Update method';
        updateItem.documentation = 'Update the specified resource';
        items.push(updateItem);

        const destroyItem = new vscode.CompletionItem('Destroy', vscode.CompletionItemKind.Method);
        destroyItem.detail = 'Destroy method';
        destroyItem.documentation = 'Remove the specified resource';
        items.push(destroyItem);

        return items;
    }

    private getModelCompletions(): vscode.CompletionItem[] {
        const items: vscode.CompletionItem[] = [];

        const modelItem = new vscode.CompletionItem('Model', vscode.CompletionItemKind.Class);
        modelItem.detail = 'Base Model';
        modelItem.documentation = 'Base model class for Laravel-Go';
        items.push(modelItem);

        const tableItem = new vscode.CompletionItem('TableName', vscode.CompletionItemKind.Field);
        tableItem.detail = 'Table name';
        tableItem.documentation = 'Specify the table name for the model';
        items.push(tableItem);

        const fillableItem = new vscode.CompletionItem('Fillable', vscode.CompletionItemKind.Field);
        fillableItem.detail = 'Fillable fields';
        fillableItem.documentation = 'Specify which fields can be mass assigned';
        items.push(fillableItem);

        const hiddenItem = new vscode.CompletionItem('Hidden', vscode.CompletionItemKind.Field);
        hiddenItem.detail = 'Hidden fields';
        hiddenItem.documentation = 'Specify which fields should be hidden from JSON output';
        items.push(hiddenItem);

        const timestampsItem = new vscode.CompletionItem('Timestamps', vscode.CompletionItemKind.Field);
        timestampsItem.detail = 'Timestamps';
        timestampsItem.documentation = 'Enable or disable timestamps';
        items.push(timestampsItem);

        return items;
    }

    private getMiddlewareCompletions(): vscode.CompletionItem[] {
        const items: vscode.CompletionItem[] = [];

        const middlewareItem = new vscode.CompletionItem('Middleware', vscode.CompletionItemKind.Class);
        middlewareItem.detail = 'Base Middleware';
        middlewareItem.documentation = 'Base middleware class for Laravel-Go';
        items.push(middlewareItem);

        const handleItem = new vscode.CompletionItem('Handle', vscode.CompletionItemKind.Method);
        handleItem.detail = 'Handle method';
        handleItem.documentation = 'Handle the incoming request';
        items.push(handleItem);

        const terminateItem = new vscode.CompletionItem('Terminate', vscode.CompletionItemKind.Method);
        terminateItem.detail = 'Terminate method';
        terminateItem.documentation = 'Handle the response after it is sent to the browser';
        items.push(terminateItem);

        return items;
    }

    private getRouteCompletions(): vscode.CompletionItem[] {
        const items: vscode.CompletionItem[] = [];

        const routeItem = new vscode.CompletionItem('Route', vscode.CompletionItemKind.Class);
        routeItem.detail = 'Route Manager';
        routeItem.documentation = 'Route management for Laravel-Go';
        items.push(routeItem);

        const getItem = new vscode.CompletionItem('Route.Get', vscode.CompletionItemKind.Method);
        getItem.detail = 'GET route';
        getItem.documentation = 'Register a GET route';
        items.push(getItem);

        const postItem = new vscode.CompletionItem('Route.Post', vscode.CompletionItemKind.Method);
        postItem.detail = 'POST route';
        postItem.documentation = 'Register a POST route';
        items.push(postItem);

        const putItem = new vscode.CompletionItem('Route.Put', vscode.CompletionItemKind.Method);
        putItem.detail = 'PUT route';
        putItem.documentation = 'Register a PUT route';
        items.push(putItem);

        const deleteItem = new vscode.CompletionItem('Route.Delete', vscode.CompletionItemKind.Method);
        deleteItem.detail = 'DELETE route';
        deleteItem.documentation = 'Register a DELETE route';
        items.push(deleteItem);

        const patchItem = new vscode.CompletionItem('Route.Patch', vscode.CompletionItemKind.Method);
        patchItem.detail = 'PATCH route';
        patchItem.documentation = 'Register a PATCH route';
        items.push(patchItem);

        return items;
    }

    private getDatabaseCompletions(): vscode.CompletionItem[] {
        const items: vscode.CompletionItem[] = [];

        const dbItem = new vscode.CompletionItem('DB', vscode.CompletionItemKind.Class);
        dbItem.detail = 'Database Manager';
        dbItem.documentation = 'Database management for Laravel-Go';
        items.push(dbItem);

        const tableItem = new vscode.CompletionItem('DB.Table', vscode.CompletionItemKind.Method);
        tableItem.detail = 'Table query';
        tableItem.documentation = 'Begin a fluent query against a database table';
        items.push(tableItem);

        const selectItem = new vscode.CompletionItem('Select', vscode.CompletionItemKind.Method);
        selectItem.detail = 'Select columns';
        selectItem.documentation = 'Select specific columns from the table';
        items.push(selectItem);

        const whereItem = new vscode.CompletionItem('Where', vscode.CompletionItemKind.Method);
        whereItem.detail = 'Where clause';
        whereItem.documentation = 'Add a where clause to the query';
        items.push(whereItem);

        const orderByItem = new vscode.CompletionItem('OrderBy', vscode.CompletionItemKind.Method);
        orderByItem.detail = 'Order by';
        orderByItem.documentation = 'Add an order by clause to the query';
        items.push(orderByItem);

        const limitItem = new vscode.CompletionItem('Limit', vscode.CompletionItemKind.Method);
        limitItem.detail = 'Limit results';
        limitItem.documentation = 'Limit the number of results returned';
        items.push(limitItem);

        return items;
    }
} 