package services

import (
	"os"
	"time"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/hasher"
)

type Scanner struct {
	InProgress bool
	fileCount  int64
	path       string
	extensions []string
	db         *dao.Database
	id         string
}

func (scanner *Scanner) FileCount() int64 {
	return scanner.fileCount
}

func NewScanner(config *config.Config, database *dao.Database) *Scanner {

	return &Scanner{
		path:       config.Path,
		extensions: config.Extensions,
		db:         database,
		id:         hasher.GetHashFromInt(time.Now().Unix()),
	}
}

func (scanner *Scanner) Start(update bool, delete bool) {

	if !scanner.InProgress {

		scanner.fileCount = 0
		scanner.InProgress = true

		go func() {

			IterateFiles(scanner.path, scanner.extensions, func(path string, info *os.FileInfo) {

				if update {
					data, err := GetMusicData(path)
					if err != nil {
						seelog.Errorf("Could not get music data: %v", err)
						return
					}
					scanner.db.PutSong(&data, scanner.id)
				} else {
					s := scanner.db.GetSongFromPath(path)
					if s == nil {
						data, err := GetMusicData(path)
						if err != nil {
							seelog.Errorf("Could not get music data: %v", err)
							return
						}
						scanner.db.PutSong(&data, scanner.id)
					} else {
						// Update song's scan ID
					}
				}

				scanner.fileCount = scanner.fileCount + 1
			})

			if delete {
				scanner.db.DeleteMissingSongs(scanner.id)
			}

			scanner.InProgress = false
		}()
	}

	return
}
