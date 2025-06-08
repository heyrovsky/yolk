package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func AbsPath(folder string) (string, error) {
	absPath, err := filepath.Abs(folder)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	return absPath, nil
}

// CreateFileWithDirs ensures the parent directory exists, then creates the file
func CreateFileWithDirs(filePath string) (*os.File, error) {
	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", filePath, err)
	}

	return file, nil
}

// DeleteFileIfExists deletes the file if it exists, otherwise does nothing.
func DeleteFileIfExists(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("failed to delete file %s: %w", filePath, err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check file %s: %w", filePath, err)
	}
	return nil
}
