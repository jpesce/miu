package main

import (
	"log"
	"miu/modules/build"
	"miu/modules/cache"
	"miu/modules/file"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	dir := map[string]string{
		"content":     "content",
		"stylesheets": "layout/style",
		"images":      "layout/images",
		"target":      "public",
	}

	if strings.Contains(os.Getenv("MIU_FLAGS"), "nocache") {
		cache.ClearAllCache()
	}

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		log.Println("Building stylesheets")
		error := file.CopyDir(dir["stylesheets"], filepath.Join(dir["target"], "style"))
		if error != nil {
			log.Fatal(error)
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		log.Println("Building global images")
		error := file.CopyDir(dir["images"], filepath.Join(dir["target"], "images"))
		if error != nil {
			log.Fatal(error)
		}
	}()

	nodes := []build.SiteNode{}
	waitGroup.Add(1)
	go func(nodes *[]build.SiteNode) {
		log.Println("Building markdown content")
		defer waitGroup.Done()

		error := error(nil)
		*nodes, error = build.BuildMarkdownContent(dir["content"], dir["target"])
		if error != nil {
			log.Fatal(error)
		}
	}(&nodes)

	waitGroup.Wait()

	log.Println("Building index")
	error := build.BuildIndex(nodes, dir["target"])
	if error != nil {
		log.Fatal(error)
	}

	elapsed := time.Since(start)
	log.Println("Built in ", elapsed)
}
