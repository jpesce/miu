package template

/*
Wrapper around "html/template" to render templates
*/

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

const TemplateDir = "layout/templates"

type HTML = template.HTML

// Render template to a given file
func RenderTemplateToFile(templateNames []string, data any, destinationPath string) error {
	templatePaths := templateNamesToPaths(templateNames)
	parsedTemplate, error := template.ParseFiles(templatePaths...)
	if error != nil {
		return fmt.Errorf("RenderTemplateToFile: %w", error)
	}

	error = os.MkdirAll(filepath.Dir(destinationPath), 0755)
	if error != nil {
		return fmt.Errorf("RenderTemplateToFile: %w", error)
	}

	file, error := os.Create(destinationPath)
	if error != nil {
		return fmt.Errorf("RenderTemplateToFile: %w", error)
	}
	defer file.Close()

	buffer := bufio.NewWriter(file)
	error = parsedTemplate.Execute(buffer, data)
	if error != nil {
		return fmt.Errorf("RenderTemplateToFile: %w", error)
	}
	buffer.Flush()

	return nil
}

// Render template to a string
func RenderTemplateToString(templateFullPath []string, data any) (string, error) {
	parsedTemplate, error := template.ParseFiles(templateFullPath...)
	if error != nil {
		return "", fmt.Errorf("RenderTemplateToFile: %w", error)
	}

	buffer := new(bytes.Buffer)
	error = parsedTemplate.Execute(buffer, data)
	if error != nil {
		return "", fmt.Errorf("RenderTemplateToFile: %w", error)
	}

	return buffer.String(), nil
}

// Given an array of template names, return their full paths
func templateNamesToPaths(templateNames []string) []string {
	templatePaths := []string{}
	for _, templateName := range templateNames {
		templatePaths = append(templatePaths, filepath.Join(TemplateDir, templateName+".tmpl.html"))
	}
	return templatePaths
}
