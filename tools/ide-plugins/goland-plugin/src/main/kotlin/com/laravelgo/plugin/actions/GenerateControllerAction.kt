package com.laravelgo.plugin.actions

import com.intellij.openapi.actionSystem.AnAction
import com.intellij.openapi.actionSystem.AnActionEvent
import com.intellij.openapi.actionSystem.CommonDataKeys
import com.intellij.openapi.project.Project
import com.intellij.openapi.ui.Messages
import com.laravelgo.plugin.LaravelGoProjectService
import com.laravelgo.plugin.dialogs.GenerateControllerDialog

class GenerateControllerAction : AnAction() {
    
    override fun actionPerformed(e: AnActionEvent) {
        val project = e.getData(CommonDataKeys.PROJECT) ?: return
        val projectService = LaravelGoProjectService.getInstance(project)
        
        if (!projectService.isLaravelGoProject()) {
            Messages.showErrorDialog(
                project,
                "This is not a Laravel-Go project. Please open a Laravel-Go project to use this feature.",
                "Not a Laravel-Go Project"
            )
            return
        }
        
        val dialog = GenerateControllerDialog(project)
        if (dialog.showAndGet()) {
            val controllerName = dialog.getControllerName()
            if (controllerName.isNotEmpty()) {
                generateController(project, controllerName)
            }
        }
    }
    
    override fun update(e: AnActionEvent) {
        val project = e.getData(CommonDataKeys.PROJECT)
        val projectService = project?.let { LaravelGoProjectService.getInstance(it) }
        e.presentation.isEnabledAndVisible = projectService?.isLaravelGoProject() == true
    }
    
    private fun generateController(project: Project, controllerName: String) {
        val projectService = LaravelGoProjectService.getInstance(project)
        val process = projectService.executeArtisanCommand("make:controller", listOf(controllerName))
        
        if (process != null) {
            Messages.showInfoMessage(
                project,
                "Controller '$controllerName' generated successfully!",
                "Controller Generated"
            )
        } else {
            Messages.showErrorDialog(
                project,
                "Failed to generate controller. Please check if artisan is available.",
                "Generation Failed"
            )
        }
    }
} 