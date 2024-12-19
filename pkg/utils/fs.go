package utils

import (
	"os"
)

func IsDir(dir string) bool {
	f, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return f.IsDir()
}

func IsFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}

	return f.Mode().IsRegular()
}
