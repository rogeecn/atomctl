package utils

import (
	"os"
	"strings"
)

func IsDir(dir string) bool {
	fs, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return fs.IsDir()

}

func IsFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true
}

func IsTplFile(path string) bool {
	return strings.HasSuffix(path, ".tpl")
}

func TplToGo(path string) string {
	return strings.Replace(path, ".tpl", ".go", 1)
}
