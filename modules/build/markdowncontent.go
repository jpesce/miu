package build

/*
Build content written in markdown
*/

import (
  "miu/modules/file"
  "miu/modules/markdown"
  "miu/modules/template"
  "miu/modules/image"
	"fmt"
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
        newNode, error := BuildMarkdownContentFile(filePath, sourceDirectory, templateDirectory, targetDirectory)
        if error != nil {
          return nodes, fmt.Errorf("BuildMarkdownContentDirectory: %w", error)
        }

        metadata, error := markdown.GetFrontmatterFromFile(filePath)
        if error != nil {
          return nil, fmt.Errorf("BuildMarkdownContentDirectory: %w", error)
        }

        if(metadata["thumbnail"] != "") {
          error := buildThumbnail(filePath, targetDirectory)
          if error != nil {
            return nil, fmt.Errorf("BuildMarkdownContentDirectory: %w", error)
          }
        }

        nodes = append(nodes, newNode)
      } else if fileExtension == ".jpg" || fileExtension == ".png" {
        // Images
        destinationPath := file.ReplaceRootDirectory(filePath, targetDirectory)

        error = file.CopyFileToDestination(filePath, destinationPath)
        if error != nil {
          return nil, fmt.Errorf("BuildMarkdownContentDirectory: %w", error)
        }
      } else {
        // Other files should be simply copied directly
        destinationPath := file.ReplaceRootDirectory(filePath, targetDirectory)
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

/* Create optimized thumbnail file in the target directory */
func buildThumbnail(filePath string, targetDirectory string) error {
  metadata, error := markdown.GetFrontmatterFromFile(filePath)
  if error != nil {
    return fmt.Errorf("buildThumbnail: %w", error)
  }

  width := 720
  if(metadata["thumbnail-wide"] == "true") { width = 1440; }

  imageSourcePath := filepath.Join(filepath.Dir(filePath), metadata["thumbnail"])
  imageDestinationPath := file.ReplaceRootDirectory(image.GetImageNameWithTag(imageSourcePath, "thumbnail"), targetDirectory)

  error = image.CompressImage(imageSourcePath, imageDestinationPath, width)
  if error != nil {
    return fmt.Errorf("buildThumbnail: %w", error)
  }

  return nil
}
