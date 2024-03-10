package image

/*
Image processor module
*/

import (
	"fmt"
	"miu/modules/file"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// Compress image to the desired maximum width and create necessary directories along the way
func CompressImage(imageSourcePath string, imageDestinationPath string, maxWidth int) error {
	error := os.MkdirAll(filepath.Dir(imageDestinationPath), 0755)
	if error != nil {
		return fmt.Errorf("CompressImage: %w", error)
	}

	// error = exec.Command(
	//   "cp",
	//   imageSourcePath,
	//   imageDestinationPath,
	// ).Run()
	// if error != nil {
	//   return fmt.Errorf("CompressImage: %w", error)
	// }

	error = exec.Command(
		"convert",
		imageSourcePath,
		"-thumbnail", strconv.Itoa(maxWidth)+"x>",
		imageDestinationPath,
	).Run()
	if error != nil {
		return fmt.Errorf("CompressImage: %w", error)
	}

	error = exec.Command(
		"cwebp",
		"-preset", "photo",
		"-metadata", "all",
		"-q", "85",
		imageDestinationPath,
		"-o", imageDestinationPath,
	).Run()
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
