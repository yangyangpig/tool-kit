package util

import (
	"os"
	"path/filepath"
	"strings"
)

func GetProtoName(path string) (name string) {
	if strings.Contains(path, ".proto") {
		fileName := filepath.Base(path)
		return strings.Split(fileName, ".")[0]
	}
	return
}

func MkDirIfNotExist(dir string) error {
	if len(dir) == 0 {
		return nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}


