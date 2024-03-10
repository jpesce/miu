package main

import (
	"log"
	"miu/modules/build"
	"miu/modules/file"
	"path/filepath"
	"time"
)

func main() {
	start := time.Now()
	directory := map[string]string{
		"content":     "content",
		"templates":   "layout/templates",
		"stylesheets": "layout/style",
		"images":      "layout/images",
		"target":      "public",
	}

	log.Println("Building stylesheets")
	error := file.CopyDirectoryRecursively(directory["stylesheets"], filepath.Join(directory["target"], "style"))
	if error != nil {
		log.Fatal(error)
	}

	log.Println("Building global images")
	error = file.CopyDirectoryRecursively(directory["images"], filepath.Join(directory["target"], "images"))
	if error != nil {
		log.Fatal(error)
	}

	log.Println("Building markdown content")
	nodes, error := build.BuildMarkdownContentDirectory(directory["content"], directory["templates"], directory["target"])
	if error != nil {
		log.Fatal(error)
	}

	log.Println("Building index")
	error = build.BuildIndex(nodes, directory["templates"], directory["target"])
	if error != nil {
		log.Fatal(error)
	}

	elapsed := time.Since(start)
	log.Println("Built in ", elapsed)
}
