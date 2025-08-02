package com.laravelgo.plugin

import com.intellij.openapi.project.Project
import com.intellij.openapi.startup.StartupActivity

class LaravelGoPlugin : StartupActivity {
    override fun runActivity(project: Project) {
        // Initialize Laravel-Go plugin for the project
        LaravelGoProjectService.getInstance(project).initialize()
    }
} 