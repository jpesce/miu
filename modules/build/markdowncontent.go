package build

/*
Build content written in markdown
*/

import (
	"fmt"
	"miu/modules/file"
	"miu/modules/markdown"
	"miu/modules/template"
	"os"
	"path/filepath"
	"strings"
)

type SiteNode struct {
  sourcePath string
  destinationPath string
}

// Given a directory, build the files and output to the target directory. If there's another
// directory inside it, recursively call itself
func BuildMarkdownContentDirectory(sourceDirectory string, templateDirectory string, targetDirectory string) ([]SiteNode, error) {
  entries, error := os.ReadDir(sourceDirectory)
  if error != nil {
    return nil, fmt.Errorf("BuildMarkdownContentDirectory: %w", error)
  }

  var nodes []SiteNode
  for _, entry := range entries {
    entryName := entry.Name()
    entryPath := filepath.Join(sourceDirectory, entryName)

    // Ignore hidden (.*) files
    if string(entryName[0]) == "." { continue }

    // If it's a directory, build it recursively
    if entry.IsDir() {
      newNodes, error := BuildMarkdownContentDirectory(entryPath, templateDirectory, targetDirectory)
      nodes = append(nodes, newNodes...)
      if error != nil {
        return nodes, fmt.Errorf("BuildMarkdownContentDirectory: %w", error)
      }
    } else {
      fileName := entryName
      filePath := entryPath
      fileExtension := filepath.Ext(fileName)

      if fileExtension == ".md" {
        // Markdown files should be built transforming it to HTML and using the appropriate template
        newNode, error := BuildMarkdownContentFile(filePath, templateDirectory, targetDirectory)
        if error != nil {
          return nodes, fmt.Errorf("BuildMarkdownContentDirectory: %w", error)
        }

        nodes = append(nodes, newNode)
      } else {
        // Other files (images and other assets), should be simply copied directly and are not
        // considered a site node
        destinationPathParts := strings.Split(filePath, "/")
        destinationPathParts[0] = targetDirectory
        destinationPath := filepath.Join(destinationPathParts...)

        file.CopyFileToDestination(filePath, destinationPath)
      }
    }
  }

  return nodes, nil
}

// Compile markdown content file to full HTML page in the destination
func BuildMarkdownContentFile(filePath string, templateDirectory string, targetDirectory string) (SiteNode, error) {
  destinationPathParts := strings.Split(filePath, "/")
  destinationPathParts[0] = targetDirectory
  destinationPath := filepath.Join(destinationPathParts...)

  depth := len(destinationPathParts) - 1
  if depth == 1 {
    // When markdown is in the root of content directory, create a directory with its name.
    // E.g., "content/example.md" -> "public/example/index.html"
    destinationPath = filepath.Join(file.PathWithoutExtension(destinationPath), "index.html")
  } else {
    // When markdown is not in the root, use its directory name.
    // E.g. "content/example/anything.md" -> "public/example/index.html"
    destinationPath = filepath.Join(append(destinationPathParts[:depth], "index.html")...)
  }

  html, error := markdown.MarkdownToHtml(filePath)
  markdowncontentTemplateData := struct{
    Content template.HTML
  }{
    Content: template.HTML(html),
  }
  markdowncontentTemplatePath := filepath.Join(templateDirectory, "markdowncontent.tmpl.html")
  markdowncontentContent, error := template.RenderTemplateToString(markdowncontentTemplatePath, markdowncontentTemplateData)
  if error != nil {
    return SiteNode{}, fmt.Errorf("BuildMarkdownContentFile: %w", error)
  }

  if error != nil {
    return SiteNode{}, fmt.Errorf("BuildMarkdownContentFile: %w", error)
  }
  mainTemplateData := mainTemplateData {
    Content: template.HTML(markdowncontentContent),
  }
  mainTemplatePath := filepath.Join(templateDirectory, "main.tmpl.html")
  template.RenderTemplateToFile(mainTemplatePath, mainTemplateData, destinationPath)

  return SiteNode{
    sourcePath: filePath,
    destinationPath: destinationPath,
  }, nil
}
