import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';

export class CommandExecutor {
    private artisanPath: string;
    private projectRoot: string;

    constructor(artisanPath: string, projectRoot: string) {
        this.artisanPath = artisanPath;
        this.projectRoot = projectRoot;
    }

    private async getProjectRoot(): Promise<string> {
        if (vscode.workspace.workspaceFolders && vscode.workspace.workspaceFolders.length > 0) {
            return vscode.workspace.workspaceFolders[0].uri.fsPath;
        }
        return this.projectRoot;
    }

    private async executeCommand(command: string, args: string[] = []): Promise<void> {
        const projectRoot = await this.getProjectRoot();
        const artisanFullPath = path.join(projectRoot, this.artisanPath);

        // Check if artisan exists
        if (!fs.existsSync(artisanFullPath)) {
            vscode.window.showErrorMessage(`Artisan not found at: ${artisanFullPath}`);
            return;
        }

        const terminal = vscode.window.createTerminal('Laravel-Go Artisan');
        const fullCommand = `go run ${artisanFullPath} ${command} ${args.join(' ')}`;
        
        terminal.sendText(fullCommand);
        terminal.show();
    }

    async generateController(): Promise<void> {
        const name = await vscode.window.showInputBox({
            prompt: 'Enter controller name',
            placeHolder: 'UserController'
        });

        if (name) {
            await this.executeCommand('make:controller', [name]);
        }
    }

    async generateModel(): Promise<void> {
        const name = await vscode.window.showInputBox({
            prompt: 'Enter model name',
            placeHolder: 'User'
        });

        if (name) {
            await this.executeCommand('make:model', [name]);
        }
    }

    async generateMiddleware(): Promise<void> {
        const name = await vscode.window.showInputBox({
            prompt: 'Enter middleware name',
            placeHolder: 'AuthMiddleware'
        });

        if (name) {
            await this.executeCommand('make:middleware', [name]);
        }
    }

    async generateMigration(): Promise<void> {
        const name = await vscode.window.showInputBox({
            prompt: 'Enter migration name',
            placeHolder: 'create_users_table'
        });

        if (name) {
            await this.executeCommand('make:migration', [name]);
        }
    }

    async generateTest(): Promise<void> {
        const name = await vscode.window.showInputBox({
            prompt: 'Enter test name',
            placeHolder: 'UserTest'
        });

        if (name) {
            await this.executeCommand('make:test', [name]);
        }
    }

    async runMigrate(): Promise<void> {
        await this.executeCommand('migrate');
    }

    async runSeed(): Promise<void> {
        await this.executeCommand('db:seed');
    }

    async runTest(): Promise<void> {
        await this.executeCommand('test');
    }

    async runServe(): Promise<void> {
        await this.executeCommand('serve');
    }
} 