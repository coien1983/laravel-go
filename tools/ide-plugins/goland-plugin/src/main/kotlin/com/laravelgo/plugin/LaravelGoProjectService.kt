package com.laravelgo.plugin

import com.intellij.openapi.project.Project
import com.intellij.openapi.components.Service
import com.intellij.openapi.components.ServiceManager
import com.intellij.openapi.vfs.VirtualFile
import java.io.File

@Service
class LaravelGoProjectService(private val project: Project) {
    
    companion object {
        fun getInstance(project: Project): LaravelGoProjectService {
            return project.getService(LaravelGoProjectService::class.java)
        }
    }
    
    private var isLaravelGoProject: Boolean = false
    private var artisanPath: String? = null
    
    fun initialize() {
        detectLaravelGoProject()
        if (isLaravelGoProject) {
            setupProjectConfiguration()
        }
    }
    
    private fun detectLaravelGoProject() {
        val projectRoot = project.baseDir
        val artisanFile = projectRoot.findChild("cmd")?.findChild("artisan")?.findChild("main.go")
        
        if (artisanFile != null) {
            isLaravelGoProject = true
            artisanPath = "cmd/artisan"
        }
    }
    
    private fun setupProjectConfiguration() {
        // Setup project-specific configurations
        // This could include setting up run configurations, code style settings, etc.
    }
    
    fun isLaravelGoProject(): Boolean = isLaravelGoProject
    
    fun getArtisanPath(): String? = artisanPath
    
    fun executeArtisanCommand(command: String, args: List<String> = emptyList()): Process? {
        if (!isLaravelGoProject || artisanPath == null) {
            return null
        }
        
        val projectRoot = project.basePath
        val artisanFullPath = File(projectRoot, artisanPath!!)
        
        if (!artisanFullPath.exists()) {
            return null
        }
        
        val commandList = mutableListOf("go", "run", artisanFullPath.absolutePath, command)
        commandList.addAll(args)
        
        return try {
            ProcessBuilder(commandList)
                .directory(File(projectRoot))
                .start()
        } catch (e: Exception) {
            null
        }
    }
} 