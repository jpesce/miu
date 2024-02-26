package markdown

/*
Utils for rendering and extracting information from markdown files
Currently using pandoc to render markdown to other formats and a very simple parser for frontmatter
data
*/

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Return compiled HTML string from Markdown file
func MarkdownToHtml(sourceFullPath string) (string, error) {
  output, error := exec.Command(
    "pandoc",
    sourceFullPath,
  ).Output()
  if error != nil {
    return "", fmt.Errorf("MarkdownToHtml: %w", error)
  }

  return string(output[:]), nil
}

// Get frontmatter from a file and return a map of strings for each field. Assumes a very simple
// format of key and value separated by a colon (:)
func GetFrontmatterFromFile(sourceFullPath string) (map[string]string, error) {
  metadata := make(map[string]string)

  sourceFile, error := os.Open(sourceFullPath)
  if error != nil {
    return nil, fmt.Errorf("GetFrontmatterFromFile: %w", error)
  }

  reader := bufio.NewReader(sourceFile)

  line, error := reader.ReadString('\n')
  if error != nil {
    return nil, fmt.Errorf("GetFrontmatterFromFile: %w", error)
  }

  // No frontmatter
  if line != "---\n" {
    return nil, nil
  }

  line, error = reader.ReadString('\n')
  if error != nil {
    return nil, fmt.Errorf("GetFrontmatterFromFile: %w", error)
  }

  for line != "---\n" {
    splittedString := strings.Split(line, ":")
    key := splittedString[0]
    value := strings.Trim(strings.Join(splittedString[1:], ":"), " ")
    value = strings.TrimSuffix(value, "\n")
    metadata[key] = value

    line, error = reader.ReadString('\n')
    if error != nil {
      return nil, fmt.Errorf("GetFrontmatterFromFile: %w", error)
    }
  }

  return metadata, nil
}
