package build

/*
Build the site's index
*/

import (
	"fmt"
	"log"
	"miu/modules/file"
	"miu/modules/image"
	"miu/modules/markdown"
	"miu/modules/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Structs
type indexTemplateData struct {
	Posts []indexTemplatePostData
}

type indexTemplatePostData struct {
	Title         string
	Url           string
	Thumbnail     string
	ThumbnailWide string
	PublishedDate time.Time
	Caption       string
}

// Public functions
func BuildIndex(nodes []SiteNode, targetDir string) error {
	postsData, error := frontmatterToPostsDataStruct(nodes)
	if error != nil {
		return fmt.Errorf("BuildIndex: %w", error)
	}

	sort.Slice(postsData, func(i, j int) bool {
		return postsData[i].PublishedDate.After(postsData[j].PublishedDate)
	})
	// for _, post := range postsData {
	// 	fmt.Println(post.Title, post.PublishedDate)
	// }

	indexTemplateData := indexTemplateData{
		Posts: postsData,
	}
	targetFile := filepath.Join(targetDir, "index.html")
	template.RenderTemplateToFile([]string{"main", "index"}, indexTemplateData, targetFile)

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
		if metadata["showInMainPage"] == "false" {
			continue
		}

		// Title
		postsData[i].Title = metadata["title"]
		if postsData[i].Title == "" {
			postsData[i].Title = file.FileNameWithoutExtension(node.sourcePath)
		}

		// Published date
		postsData[i].PublishedDate, error = time.Parse("2006-01-02", metadata["date"])
		if error != nil {
			postStat, error := os.Stat(node.sourcePath)
			if error != nil {
				return nil, fmt.Errorf("BuildIndex: %w", error)
			}
			postsData[i].PublishedDate = postStat.ModTime()
			if strings.Contains(os.Getenv("MIU_FLAGS"), "verbose") {
				log.Println("INFO: Unknown published date for", node.sourcePath, "using last modified date")
			}
		}

		// URL is file path minus first dir (public/) and the filename (index.html)
		destinationPathParts := strings.Split(node.destinationPath, "/")
		destinationUrl := filepath.Join(destinationPathParts[1 : len(destinationPathParts)-1]...)
		postsData[i].Url = destinationUrl

		// Thumbnail wide
		postsData[i].ThumbnailWide = metadata["thumbnail-wide"]

		// Thumbnail
		if metadata["thumbnail"] != "" {
			thumbnailPath := filepath.Join(destinationUrl, metadata["thumbnail"])
			postsData[i].Thumbnail = image.GetImageNameWithTag(thumbnailPath, "thumbnail")
		}

		// Description
		if metadata["description"] != "" {
			year := strings.Split(metadata["date"], "-")[0]
			postsData[i].Caption = fmt.Sprintf("%s (%s) /  %s", postsData[i].Title, year, metadata["description"])
		}
	}

	return postsData, nil
}
