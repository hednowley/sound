package services_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hednowley/sound/services"
)

func TestScanner(t *testing.T) {

	errors := []string{}
	files := make(map[string]*bool)

	b1 := false
	b2 := false
	b3 := false

	files[filepath.Join("..", "testdata", "music", "1.mp3")] = &b1
	files[filepath.Join("..", "testdata", "music", "2.mp3")] = &b2
	files[filepath.Join("..", "testdata", "music", "subfolder", "3.mp3")] = &b3

	e := []string{
		"mp3",
		"flac",
	}

	path := filepath.Join("..", "testdata", "music")
	services.IterateFiles(path, e, func(path string, info *os.FileInfo) {
		for k, v := range files {
			if strings.HasSuffix(path, k) {
				if *v {
					errors = append(errors, fmt.Sprintf("Double scan: %v", k))
					return
				}
				*files[k] = true
				return
			}
		}

		errors = append(errors, fmt.Sprintf("Unexpected scan: %v", path))
	})

	if len(errors) > 0 {
		for _, v := range errors {
			t.Error(v)
		}
	}

	for k, v := range files {
		if !*v {
			t.Error(fmt.Sprintf("Unscanned file: %v", k))
		}
	}

}
