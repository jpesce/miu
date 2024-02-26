package template

/*
Wrapper around "html/template" to inject content in templates
*/

import (
  "html/template"
  "os"
  "fmt"
  "bufio"
  "bytes"
  "path/filepath"
)

type HTML = template.HTML

// Render template to a given file
func RenderTemplateToFile(templateFullPath string, content any, destinationPath string) (error) {
  parsedTemplate, error := template.ParseFiles(templateFullPath)
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
  error = parsedTemplate.Execute(buffer, content)
  if error != nil {
    return fmt.Errorf("RenderTemplateToFile: %w", error)
  }
  buffer.Flush()

  return nil
}

// Render template to a string and return it
func RenderTemplateToString(templateFullPath string, content any) (string, error) {
  parsedTemplate, error := template.ParseFiles(templateFullPath)
  if error != nil {
    return "", fmt.Errorf("RenderTemplateToFile: %w", error)
  }

  buffer := new(bytes.Buffer)
  error = parsedTemplate.Execute(buffer, content)
  if error != nil {
    return "", fmt.Errorf("RenderTemplateToFile: %w", error)
  }

  return buffer.String(), nil
}
