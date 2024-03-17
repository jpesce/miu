package cache

import (
	"errors"
	"fmt"
	"miu/modules/file"
	"os"
	"path/filepath"
	"time"
)

const CacheDir = "cache"

// Given a file path, get the corresponding file in cache
func GetCachePath(path string) string {
	return filepath.Join(CacheDir, path)
}

// Look at how fresh the cache is and determine if it should be used or not
func ShouldUseCache(dependeciesPaths []string, destinationPath string) (bool, error) {
	latestChange, error := getLatestChange(dependeciesPaths)
	if error != nil {
		return false, fmt.Errorf("ShouldUseCache: %w", error)
	}

	shouldUseCache := false
	cachePath := GetCachePath(destinationPath)

	cacheStats, error := os.Stat(cachePath)
	if error == nil && !latestChange.After(cacheStats.ModTime()) {
		shouldUseCache = true
	} else if error != nil && !errors.Is(error, os.ErrNotExist) {
		return false, fmt.Errorf("ShouldUseCache: %w", error)
	}
	return shouldUseCache, nil
}

func ClearAllCache() error {
	return file.RemoveDirContents(CacheDir)
}

// Given a slice of files, return the last modifed date among them
func getLatestChange(paths []string) (time.Time, error) {
	latestChangeStat, error := os.Stat(paths[0])
	if error != nil {
		return time.Time{}, fmt.Errorf("getLatestChange: %w", error)
	}
	latestChange := latestChangeStat.ModTime()

	for _, dependency := range paths {
		dependencyStat, error := os.Stat(dependency)
		if error != nil {
			return time.Time{}, fmt.Errorf("getLatestChange: %w", error)
		}
		dependencyLatestChange := dependencyStat.ModTime()
		if dependencyLatestChange.After(latestChange) {
			latestChange = dependencyLatestChange
		}
	}
	return latestChange, nil
}
