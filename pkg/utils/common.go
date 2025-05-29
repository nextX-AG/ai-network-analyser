package utils

import (
	"os"
	"path/filepath"
)

// GetExecutablePath returns the path of the current executable
func GetExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// EnsureDirectory ensures that a directory exists, creating it if necessary
func EnsureDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}
