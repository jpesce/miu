package build

/*
Build content written in markdown
*/

import (
	"fmt"
	"miu/modules/cache"
	"miu/modules/file"
	"miu/modules/image"
	"miu/modules/markdown"
	"miu/modules/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type markdowncontentTemplateData struct {
	Title   string
	Content template.HTML
}

type SiteNode struct {
	sourcePath      string
	destinationPath string
}

// Given a directory, build the files and output to the target directory. If there's another directory inside it, recursively call itself
func BuildMarkdownContent(sourceDir string, targetDir string) ([]SiteNode, error) {
	waitGroup := sync.WaitGroup{}

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("BuildMarkdownContent: %w", err)
	}

	var (
		nodesMu sync.Mutex
		nodes   []SiteNode
	)
	for _, entry := range entries {
		entryName := entry.Name()
		entryPath := filepath.Join(sourceDir, entryName)

		// Ignore hidden (.*) files
		if string(entryName[0]) == "." {
			continue
		}

		// If it's a directory, build it recursively
		if entry.IsDir() {
			waitGroup.Add(1)
			newNodes := []SiteNode{}
			go func() {
				defer waitGroup.Done()
				newNodes, _ = BuildMarkdownContent(entryPath, targetDir) // TODO: properly deal with error

				nodesMu.Lock()
				nodes = append(nodes, newNodes...)
				nodesMu.Unlock()
			}()
		} else {
			fileName := entryName
			filePath := entryPath
			fileExtension := filepath.Ext(fileName)

			if fileExtension == ".md" {
				waitGroup.Add(1)
				newNode := SiteNode{}
				go func() {
					defer waitGroup.Done()
					newNode, _ = BuildMarkdownContentFile(filePath, sourceDir, targetDir) // TODO: properly deal with error

					nodesMu.Lock()
					nodes = append(nodes, newNode)
					nodesMu.Unlock()
				}()

				metadata, error := markdown.GetFrontmatterFromFile(filePath)
				if error != nil {
					return nil, fmt.Errorf("BuildMarkdownContent: %w", error)
				}

				if metadata["thumbnail"] != "" {
					waitGroup.Add(1)
					go func() {
						defer waitGroup.Done()
						_ = buildThumbnail(filePath, targetDir) // TODO: properly deal with error
					}()
				}
			} else if fileExtension == ".jpg" || fileExtension == ".png" {
				// Images
				destinationPath := file.ReplaceRootDir(filePath, targetDir)

				err = file.CopyFile(filePath, destinationPath)
				if err != nil {
					return nil, fmt.Errorf("BuildMarkdownContent: %w", err)
				}
			} else {
				// Other files should be simply copied directly
				destinationPath := file.ReplaceRootDir(filePath, targetDir)
				file.CopyFile(filePath, destinationPath)
			}
		}
	}

	waitGroup.Wait()
	return nodes, nil
}

// Compile markdown content file to full HTML page in the destination
func BuildMarkdownContentFile(filePath string, contentDir string, targetDir string) (SiteNode, error) {
	destinationPath := file.ReplaceRootDir(filePath, targetDir)

	if strings.Split(filePath, "/")[0] == contentDir {
		// When markdown is in the root of content directory, create a directory with its name.
		// e.g., "content/example.md" -> "public/example/index.html"
		destinationPath = filepath.Join(file.PathWithoutExtension(destinationPath), "index.html")
	} else {
		// When markdown is not in the root, use its directory name.
		// e.g. "content/example/anything.md" -> "public/example/index.html"
		destinationPath = filepath.Join(filepath.Dir(destinationPath), "index.html")
	}

	siteNode := SiteNode{
		sourcePath:      filePath,
		destinationPath: destinationPath,
	}

	shouldUseCache, error := cache.ShouldUseCache([]string{filePath}, destinationPath)
	if error != nil {
		return SiteNode{}, fmt.Errorf("BuildMarkdownContentFile: %w", error)
	}

	if shouldUseCache {
		error = file.CopyFile(cache.GetCachePath(destinationPath), destinationPath)
		if error != nil {
			return SiteNode{}, fmt.Errorf("BuildMarkdownContentFile: %w", error)
		}

		return siteNode, nil
	}

	metadata, error := markdown.GetFrontmatterFromFile(filePath)
	if error != nil {
		return SiteNode{}, fmt.Errorf("buildThumbnail: %w", error)
	}
	title := "|||||"
	if metadata["title"] != "" {
		title = metadata["title"]
	}

	html, error := markdown.MarkdownToHtml(filePath)
	if error != nil {
		return SiteNode{}, fmt.Errorf("BuildMarkdownContentFile: %w", error)
	}

	pageUrl := getPageUrl(destinationPath)
	html = addPrefixToSrc(html, pageUrl)

	markdowncontentTemplateData := markdowncontentTemplateData{
		Title:   title,
		Content: template.HTML(html),
	}

	// Only create a page if it has any content
	if html != "" {
		template.RenderTemplateToFile([]string{"main", "markdowncontent"}, markdowncontentTemplateData, cache.GetCachePath(destinationPath))

		error = file.CopyFile(cache.GetCachePath(destinationPath), destinationPath)
		if error != nil {
			return SiteNode{}, fmt.Errorf("BuildMarkdownContentFile: %w", error)
		}
	}

	return siteNode, nil
}

// Create optimized thumbnail file in the target directory
func buildThumbnail(filePath string, targetDir string) error {
	metadata, error := markdown.GetFrontmatterFromFile(filePath)
	if error != nil {
		return fmt.Errorf("buildThumbnail: %w", error)
	}

	width := 720
	if metadata["thumbnail-wide"] == "true" {
		width = 1440
	}

	imageSourcePath := filepath.Join(filepath.Dir(filePath), metadata["thumbnail"])
	imageDestinationPath := file.ReplaceRootDir(image.GetImageNameWithTag(imageSourcePath, "thumbnail"), targetDir)

	error = image.CompressImage(imageSourcePath, imageDestinationPath, width)
	if error != nil {
		return fmt.Errorf("buildThumbnail: %w", error)
	}

	return nil
}

// Get final URL for file
func getPageUrl(destinationPath string) string {
	pageUrl := filepath.Dir(destinationPath)
	return strings.Join(strings.Split(pageUrl, "/")[1:], "/")
}

// Add URL prefix to src attributes
func addPrefixToSrc(html string, pageUrl string) string {
	srcRegexp := regexp.MustCompile(`src="([^"]*)"`)
	return srcRegexp.ReplaceAllString(html, "src=\"/"+pageUrl+"/${1}\"")
}
