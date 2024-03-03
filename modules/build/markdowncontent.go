package build

/*
Build content written in markdown
*/

import (
	"fmt"
  "os/exec"
	"miu/modules/file"
	"miu/modules/markdown"
	"miu/modules/template"
	"os"
	"path/filepath"
	"strings"
)


type markdowncontentTemplateData struct {
  Content template.HTML
}

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
        newNode, error := BuildMarkdownContentFile(filePath, sourceDirectory, templateDirectory, targetDirectory)
        if error != nil {
          return nodes, fmt.Errorf("BuildMarkdownContentDirectory: %w", error)
        }

        nodes = append(nodes, newNode)
      } else if fileExtension == ".jpg" || fileExtension == ".png" {
        // Images
        destinationPath := file.ReplaceRootDirectory(filePath, targetDirectory)

        error := os.MkdirAll(filepath.Dir(destinationPath), 0755)
        if error != nil {
          return nil,fmt.Errorf("CopyFileToDestination: %w", error)
        }

        error = exec.Command(
          "convert",
          filePath,
          "-thumbnail", "1440x>",
          "-quality", "75",
          "-colorspace", "sRGB",
          destinationPath,
        ).Run()
        if error != nil {
          return nil, fmt.Errorf("Markdowncontent: %w", error)
        }
      } else {
        // Other files should be simply copied directly
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
func BuildMarkdownContentFile(filePath string, contentDirectory string, templateDirectory string, targetDirectory string) (SiteNode, error) {
  destinationPath := file.ReplaceRootDirectory(filePath, targetDirectory)

  if strings.Split(filePath, "/")[0] == contentDirectory {
    // When markdown is in the root of content directory, create a directory with its name.
    // e.g., "content/example.md" -> "public/example/index.html"
    destinationPath = filepath.Join(file.PathWithoutExtension(destinationPath), "index.html")
  } else {
    // When markdown is not in the root, use its directory name.
    // e.g. "content/example/anything.md" -> "public/example/index.html"
    destinationPath = filepath.Join(filepath.Dir(destinationPath), "index.html")
  }

  html, error := markdown.MarkdownToHtml(filePath)
  if error != nil {
    return SiteNode{}, fmt.Errorf("BuildMarkdownContentFile: %w", error)
  }
  markdowncontentTemplateData := markdowncontentTemplateData {
    Content: template.HTML(html),
  }

  // Only create a page if it has any content
  if html != "" {
    mainTemplatePath := filepath.Join(templateDirectory, "main.tmpl.html")
    markdowncontentTemplatePath := filepath.Join(templateDirectory, "markdowncontent.tmpl.html")
    template.RenderTemplateToFile([]string{mainTemplatePath, markdowncontentTemplatePath}, markdowncontentTemplateData, destinationPath)
  }

  return SiteNode{
    sourcePath: filePath,
    destinationPath: destinationPath,
  }, nil
}
