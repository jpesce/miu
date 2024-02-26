package file

/*
Utils to deal with files
*/

import (
  "os"
  "path/filepath"
  "strings"
  "os/exec"
  "fmt"
)

// Given a complete path, return the filename without the file's extension.
// E.g. "path/to/file.go" -> "file"
// "file.go" -> "file"
func FileNameWithoutExtension(path string) string {
  fileName := filepath.Base(path)
  return PathWithoutExtension(fileName)
}

// Given a path, return the path without the file's extension.
// E.g. "path/to/file.go" -> "path/to/file"
func PathWithoutExtension(path string) string {
  return strings.TrimSuffix(path, filepath.Ext(path))
}

// Copy file to destination and create any necessary parent directories along the way
func CopyFileToDestination(sourceFile string, destinationFile string) (error) {
  destinationDirectory := filepath.Dir(destinationFile)

  error := os.MkdirAll(destinationDirectory, 0755)
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
// E.g., if the source directory is "./source" and the destination directory is "./destination", it
// copies everything to "./destination"
func CopyDirectoryRecursively(sourceDirectory string, destinationDirectory string) (error) {
  error := os.MkdirAll(destinationDirectory, 0755)
  if error != nil {
    return fmt.Errorf("CopyFileToDestination: %w", error)
  }

  entries, error := os.ReadDir(sourceDirectory)
  if error != nil {
    return fmt.Errorf("CopyDirectoryRecursively: %w", error)
  }

  for _, entry := range entries {
    entryPath := filepath.Join(sourceDirectory, entry.Name())
    entryDestinationPath := filepath.Join(destinationDirectory, entry.Name())

    if entry.IsDir() {
      error := os.MkdirAll(destinationDirectory, 0755)
      if error != nil {
        return fmt.Errorf("CopyFileToDestination: %w", error)
      }
      CopyDirectoryRecursively(entryPath, entryDestinationPath)
    } else {
      input, error := os.ReadFile(entryPath)
      if error != nil {
        return fmt.Errorf("CopyFileToDestination: %w", error)
      }

      error = os.WriteFile(entryDestinationPath, input, 0644)
      if error != nil {
        return fmt.Errorf("CopyFileToDestination: %w", error)
      }
    }
  }

  return nil
}

// Copy directory and everything inside it to destination
// E.g., if the source directory is "./source" and the destination directory is "./destination", it
// copies everything to "./destination"
func CmdCopyDirectoryRecursively(sourceDirectory string, destinationDirectory string) (error) {
  error := exec.Command(
    "cp",
    "-r",
    sourceDirectory,
    destinationDirectory,
  ).Run()
  if error != nil {
    return fmt.Errorf("CopyDirectoryRecursively: %w", error)
  }

  return nil
}
