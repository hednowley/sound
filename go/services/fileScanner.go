package services

import (
	"os"
	"path/filepath"
	"strings"
)

func hasExtension(name string, extensions []string) bool {
	for _, e := range extensions {
		if strings.HasSuffix(name, "."+e) {
			return true
		}
	}
	return false
}

func IterateFiles(root string, extensions []string, action func(path string, info *os.FileInfo)) (err error) {

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && hasExtension(info.Name(), extensions) {
			action(path, &info)
		}
		return nil
	})
	if err != nil {
		return
	}

	return
}
