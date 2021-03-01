package util

import (
	"os"
	"path/filepath"
)

func GetAbsPath(path string, autoCreate bool) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if _, absence := os.Stat(abs); absence != nil && autoCreate {
		if err := os.MkdirAll(abs, os.ModeDir); err != nil {
			return "", err
		}
	}
	return abs, nil
}
