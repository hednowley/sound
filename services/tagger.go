package services

import (
	"os"
	"path"
	"strings"

	"github.com/dhowden/tag"
	"github.com/hednowley/sound/entities"
)

func GetMusicData(filePath string) (data entities.FileData, err error) {
	file, _ := os.Open(filePath)
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return
	}

	m, _ := tag.ReadFrom(file)

	track, _ := m.Track()
	disc, _ := m.Disc()

	pic := m.Picture()
	var art *entities.CoverArtData
	if pic != nil {
		art = &entities.CoverArtData{
			Extension: pic.Ext,
			Raw:       pic.Data,
		}
	}

	var albumArtist string
	if len(m.AlbumArtist()) == 0 {
		albumArtist = m.Artist()
	} else {
		albumArtist = m.AlbumArtist()
	}

	data = entities.FileData{
		Path:        filePath,
		Artist:      m.Artist(),
		Album:       m.Album(),
		AlbumArtist: albumArtist,
		Title:       m.Title(),
		Genre:       m.Genre(),
		Year:        m.Year(),
		Track:       track,
		Disc:        disc,
		CoverArt:    art,
		Size:        info.Size(),
		Extension:   strings.TrimPrefix(path.Ext(filePath), "."),
	}
	return
}
