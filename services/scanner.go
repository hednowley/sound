package services

import (
	"os"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dao"
)

type Scanner struct {
	InProgress bool
	fileCount  int64
	path       string
	extensions []string
	db         *dao.Database
}

func (scanner *Scanner) FileCount() int64 {
	return scanner.fileCount
}

func NewScanner(config *config.Config, database *dao.Database) *Scanner {
	scanner := Scanner{
		path:       config.Path,
		extensions: config.Extensions,
		db:         database,
	}
	return &scanner
}

func (scanner *Scanner) Start() {

	if !scanner.InProgress {

		scanner.fileCount = 0
		scanner.InProgress = true

		go func() {

			IterateFiles(scanner.path, scanner.extensions, func(path string, info *os.FileInfo) {

				data, err := GetMusicData(path)
				if err != nil {
					seelog.Errorf("Could not get music data: %v", err)
					return
				}
				scanner.db.PutSong(&data)
				scanner.fileCount = scanner.fileCount + 1
			})

			scanner.InProgress = false
		}()
	}

	return
}
