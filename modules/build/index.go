package build

/*
Build the site's index
*/

import (
  "miu/modules/file"
  "miu/modules/markdown"
  "miu/modules/template"
  "fmt"
  "strings"
  "path/filepath"
  "slices"
)

// Structs
type indexTemplateData struct {
  Posts []indexTemplatePostData
}

type indexTemplatePostData struct {
  Title string
  Url string
  Thumbnail string
  ThumbnailWide string
  PublishedDate string
  Caption string
}

// Public functions
func BuildIndex(nodes []SiteNode, templateDir string, targetDir string) (error) {
  postsData, error := frontmatterToPostsDataStruct(nodes)
  if error != nil {
    return fmt.Errorf("BuildIndex: %w", error)
  }

  slices.Reverse(postsData)

  indexTemplateData := indexTemplateData {
    Posts: postsData,
  }
  indexTemplatePath := filepath.Join(templateDir, "index.tmpl.html")
  indexContent, error := template.RenderTemplateToString(indexTemplatePath, indexTemplateData)
  if error != nil {
    return fmt.Errorf("BuildIndex: %w", error)
  }

  mainTemplateData := mainTemplateData {
    Content: template.HTML(indexContent),
  }
  mainTemplatePath := filepath.Join(templateDir, "main.tmpl.html")
  targetFile := filepath.Join(targetDir, "index.html")
  template.RenderTemplateToFile(mainTemplatePath, mainTemplateData, targetFile)

  return nil
}

// Private functions
func frontmatterToPostsDataStruct(nodes []SiteNode) ([]indexTemplatePostData, error) {
  postsData := make([]indexTemplatePostData, len(nodes))

  for i, node := range nodes {
    metadata, error := markdown.GetFrontmatterFromFile(node.sourcePath)
    if error != nil {
      return nil, fmt.Errorf("BuildIndex: %w", error)
    }

    // Ignore certain content
    if metadata["showInMainPage"] == "false" { continue }

    postsData[i].Title = metadata["title"]
    if postsData[i].Title == "" {
      postsData[i].Title = file.FileNameWithoutExtension(node.sourcePath)
    }


    // URL is file path minus first dir (public/) and the filename (index.html)
    destinationPathParts := strings.Split(node.destinationPath, "/")
    destinationUrl := filepath.Join(destinationPathParts[1:len(destinationPathParts)-1]...)
    postsData[i].Url = destinationUrl

    postsData[i].ThumbnailWide = metadata["thumbnail-wide"]

    if metadata["thumbnail"] != "" {
      postsData[i].Thumbnail = filepath.Join(destinationUrl, metadata["thumbnail"])
    }


    if metadata["description"] != "" {
      year := strings.Split(metadata["date"], "-")[0]
      postsData[i].Caption = fmt.Sprintf("%s (%s) /  %s", postsData[i].Title, year, metadata["description"])
    }
  }

  return postsData, nil
}
