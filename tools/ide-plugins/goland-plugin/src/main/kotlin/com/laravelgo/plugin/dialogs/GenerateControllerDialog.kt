package com.laravelgo.plugin.dialogs

import com.intellij.openapi.project.Project
import com.intellij.openapi.ui.DialogWrapper
import com.intellij.openapi.ui.ValidationInfo
import com.intellij.ui.components.JBLabel
import com.intellij.ui.components.JBTextField
import com.intellij.util.ui.FormBuilder
import java.awt.Dimension
import javax.swing.JComponent
import javax.swing.JPanel

class GenerateControllerDialog(project: Project) : DialogWrapper(project) {
    
    private val controllerNameField = JBTextField()
    
    init {
        title = "Generate Laravel-Go Controller"
        init()
    }
    
    override fun createCenterPanel(): JComponent {
        val dialogPanel = FormBuilder.createFormBuilder()
            .addLabeledComponent(JBLabel("Controller Name: "), controllerNameField, 1, false)
            .addComponentFillVertically(JPanel(), 0)
            .panel
        
        dialogPanel.preferredSize = Dimension(400, 100)
        return dialogPanel
    }
    
    override fun doValidate(): ValidationInfo? {
        val controllerName = controllerNameField.text.trim()
        
        if (controllerName.isEmpty()) {
            return ValidationInfo("Controller name cannot be empty", controllerNameField)
        }
        
        if (!controllerName.matches(Regex("^[A-Z][a-zA-Z0-9]*Controller$"))) {
            return ValidationInfo("Controller name should be in PascalCase and end with 'Controller'", controllerNameField)
        }
        
        return null
    }
    
    fun getControllerName(): String {
        return controllerNameField.text.trim()
    }
} 