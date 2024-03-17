package image

/*
Image processor module
*/

import (
	"fmt"
	"miu/modules/cache"
	"miu/modules/file"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Compress image to the desired maximum width and create necessary directories along the way
func CompressImage(source string, destination string, maxWidth int) error {
	shouldUseCache, error := cache.ShouldUseCache([]string{source}, destination)
	if error != nil {
		return fmt.Errorf("CompressImage: %w", error)
	}

	cacheDestination := cache.GetCachePath(destination)

	if shouldUseCache {
		error = file.CopyFile(cacheDestination, destination)
		if error != nil {
			return fmt.Errorf("CompressImage: %w", error)
		}

		return nil
	}

	error = os.MkdirAll(filepath.Dir(cacheDestination), 0755)
	if error != nil {
		return fmt.Errorf("CompressImage: %w", error)
	}

	// If the flag `nocompress` is set, copy images directly to destination
	if strings.Contains(os.Getenv("MIU_FLAGS"), "nocompress") {
		error = exec.Command(
			"cp",
			source,
			cacheDestination,
		).Run()
		if error != nil {
			return fmt.Errorf("CompressImage: %w", error)
		}

		error = file.CopyFile(cacheDestination, destination)
		if error != nil {
			return fmt.Errorf("CompressImage: %w", error)
		}

		return nil
	}

	error = exec.Command(
		"convert",
		source,
		"-thumbnail", strconv.Itoa(maxWidth)+"x>",
		cacheDestination,
	).Run()
	if error != nil {
		return fmt.Errorf("CompressImage: %w", error)
	}

	error = exec.Command(
		"cwebp",
		"-preset", "photo",
		"-metadata", "all",
		"-q", "85",
		cacheDestination,
		"-o", cacheDestination,
	).Run()
	if error != nil {
		return fmt.Errorf("CompressImage: %w", error)
	}

	error = file.CopyFile(cacheDestination, destination)
	if error != nil {
		return fmt.Errorf("CompressImage: %w", error)
	}

	return nil
}

/* File name + @ + tag + extension
 * eg. image.jpg -> image@tag.jpg
 */
func GetImageNameWithTag(imagePath string, tag string) string {
	return file.PathWithoutExtension(imagePath) + "@" + tag + filepath.Ext(imagePath)
}
