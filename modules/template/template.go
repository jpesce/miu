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

type HTML = template.HTML

// Render template to a given file
func RenderTemplateToFile(templateFullPath []string, data any, destinationPath string) error {
	parsedTemplate, error := template.ParseFiles(templateFullPath...)
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
