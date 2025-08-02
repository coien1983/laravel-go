import * as vscode from 'vscode';

export class DebugConfigurationProvider implements vscode.DebugConfigurationProvider {
    provideDebugConfigurations(
        folder: vscode.WorkspaceFolder | undefined,
        token?: vscode.CancellationToken
    ): vscode.ProviderResult<vscode.DebugConfiguration[]> {
        return [
            {
                name: 'Launch Laravel-Go Application',
                type: 'go',
                request: 'launch',
                mode: 'auto',
                program: '${workspaceFolder}/cmd/artisan',
                args: ['serve'],
                env: {},
                showLog: true
            },
            {
                name: 'Debug Laravel-Go Tests',
                type: 'go',
                request: 'launch',
                mode: 'test',
                program: '${workspaceFolder}',
                args: ['-test.v'],
                env: {},
                showLog: true
            },
            {
                name: 'Debug Laravel-Go Artisan Command',
                type: 'go',
                request: 'launch',
                mode: 'auto',
                program: '${workspaceFolder}/cmd/artisan',
                args: ['${input:command}'],
                env: {},
                showLog: true
            }
        ];
    }

    resolveDebugConfiguration(
        folder: vscode.WorkspaceFolder | undefined,
        debugConfiguration: vscode.DebugConfiguration,
        token?: vscode.CancellationToken
    ): vscode.ProviderResult<vscode.DebugConfiguration> {
        // Add default environment variables for Laravel-Go
        if (!debugConfiguration.env) {
            debugConfiguration.env = {};
        }

        // Set default environment variables
        debugConfiguration.env['APP_ENV'] = debugConfiguration.env['APP_ENV'] || 'local';
        debugConfiguration.env['APP_DEBUG'] = debugConfiguration.env['APP_DEBUG'] || 'true';
        debugConfiguration.env['APP_KEY'] = debugConfiguration.env['APP_KEY'] || 'base64:your-secret-key-here';

        return debugConfiguration;
    }
} 