package file

/*
Utils to deal with files
*/

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Replace root directory of path with a new one
// e.g. "path/to/file.go" "newPath" -> "newPath/to/file.go"
func ReplaceRootDir(path string, newRoot string) string {
	destinationPathParts := strings.Split(path, string(os.PathSeparator))
	destinationPathParts[0] = newRoot
	return filepath.Join(destinationPathParts...)
}

// Given the path to a file, return the filename without the file's extension.
// e.g. "path/to/file.go" -> "file"
// "file.go" -> "file"
func FileNameWithoutExtension(path string) string {
	fileName := filepath.Base(path)
	return PathWithoutExtension(fileName)
}

// Given a path, return the path without the file's extension.
// e.g. "path/to/file.go" -> "path/to/file"
func PathWithoutExtension(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

// Copy file to destination and create any necessary directories along the way
func CopyFile(sourceFile string, destinationFile string) error {
	destinationDir := filepath.Dir(destinationFile)

	error := os.MkdirAll(destinationDir, 0755)
	if error != nil {
		return fmt.Errorf("CopyFileToDestination: %w", error)
	}

	input, error := os.ReadFile(sourceFile)
	if error != nil {
		return fmt.Errorf("CopyFileToDestination: %w", error)
	}

	error = os.WriteFile(destinationFile, input, 0644)
	if error != nil {
		return fmt.Errorf("CopyFileToDestination: %w", error)
	}

	return nil
}

// Copy directory and everything inside it to destination
// e.g., if the source directory is "./source" and the destination directory is "./destination", it
// copies everything to "./destination"
func CopyDir(source string, destination string) error {
	error := os.MkdirAll(destination, 0755)
	if error != nil {
		return fmt.Errorf("CopyDir: %w", error)
	}

	entries, error := os.ReadDir(source)
	if error != nil {
		return fmt.Errorf("CopyDir: %w", error)
	}

	for _, entry := range entries {
		entryPath := filepath.Join(source, entry.Name())
		entryDestinationPath := filepath.Join(destination, entry.Name())

		if entry.IsDir() {
			error := os.MkdirAll(destination, 0755)
			if error != nil {
				return fmt.Errorf("CopyDir: %w", error)
			}
			CopyDir(entryPath, entryDestinationPath)
		} else {
			input, error := os.ReadFile(entryPath)
			if error != nil {
				return fmt.Errorf("CopyDir: %w", error)
			}

			error = os.WriteFile(entryDestinationPath, input, 0644)
			if error != nil {
				return fmt.Errorf("CopyDir: %w", error)
			}
		}
	}

	return nil
}

// Copy directory and everything inside it to destination using `cp`
// e.g., if the source directory is "./source" and the destination directory is "./destination", it
// copies everything to "./destination"
func CmdCopyDir(sourceDir string, destinationDir string) error {
	error := exec.Command(
		"cp",
		"-r",
		sourceDir,
		destinationDir,
	).Run()
	if error != nil {
		return fmt.Errorf("CopyDirctoryRecursively: %w", error)
	}

	return nil
}

func RemoveDirContents(dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		fullPath := filepath.Join(dirPath, name)

		err = os.RemoveAll(fullPath)
		if err != nil {
			return err
		}
	}

	return nil
}
