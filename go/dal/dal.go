package dal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/cihub/seelog"
	"github.com/google/uuid"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/database"
	"github.com/hednowley/sound/entities"
	"github.com/hednowley/sound/hasher"
)

// DAL (data access layer) allows querying and writing application data.
type DAL struct {
	Db     *database.Default
	artDir string
}

// NewDAL constructs a new DAL.
func NewDAL(config *config.Config, database *database.Default) *DAL {
	return &DAL{
		Db:     database,
		artDir: config.ArtPath,
	}
}

// PutSong updates the stored song with the given path and provider ID and returns its ID.
// If there is no such song then a new one is created and its ID is returned.
// The associated album, artist, artwork and genre are created too if necessary.
func (dal *DAL) PutSong(fileInfo *entities.FileInfo, token string, providerID string) error {

	genreID, err := dal.Db.PutGenreByName(fileInfo.Genre)
	if err != nil {
		return err
	}

	art, err := dal.PutArt(fileInfo.CoverArt)
	if err != nil {
		return err
	}

	albumID, err := dal.Db.PutAlbumByAttributes(
		fileInfo.Album, fileInfo.AlbumArtist, fileInfo.Disambiguator)
	if err != nil {
		return err
	}

	var artPath string
	if art != nil {
		artPath = art.Path
	}

	_, err = dal.Db.InsertSong(
		fileInfo.Artist,
		albumID,
		fileInfo.Path,
		fileInfo.Title,
		fileInfo.Track,
		fileInfo.Disc,
		genreID,
		fileInfo.Year,
		artPath,
		fileInfo.Size,
		fileInfo.Bitrate,
		fileInfo.Duration,
		token,
		providerID,
	)

	return err
}

func (dal *DAL) PutArt(art *entities.CoverArtData) (*dao.Art, error) {
	hash := hasher.GetHash(art.Raw)
	a := dal.Db.GetArtFromHash(hash)
	if a != nil {
		// Artwork already exists
		return a, nil
	}

	filePath := fmt.Sprintf("%v.%v", uuid.New().String(), art.Extension)

	err := ioutil.WriteFile(path.Join(dal.artDir, "art", filePath), art.Raw, 0644)
	if err != nil {
		seelog.Errorf("Error saving artwork %v", filePath)
		return nil, err
	}

	// Save the record with the new path
	return dal.Db.InsertArt(filePath, hash)
}

func (dal *DAL) PutPlaylist(id uint, name string, songIDs []uint) (uint, error) {

	now := time.Now()
	var playlistID uint

	if id == 0 {
		inserted, err := dal.Db.InsertPlaylist(name, "")
		if err != nil {
			return 0, err
		}
		playlistID = inserted
	} else {
		playlist, err := dal.Db.GetPlaylist(id)
		if err != nil {
			return 0, err
		}
		if playlist == nil {
			return 0, &dao.ErrNotFound{}
		}

		playlist.Changed = &now

		var nameUpdate string
		if name == "" {
			nameUpdate = playlist.Name
		} else {
			nameUpdate = name
		}

		playlist, err = dal.Db.UpdatePlaylist(playlist.ID, nameUpdate, playlist.Comment)
		if err != nil {
			return 0, err
		}

		playlistID = playlist.ID
	}

	err := dal.Db.ReplacePlaylistEntries(playlistID, songIDs)
	if err != nil {
		return 0, err
	}

	return playlistID, nil
}

func (dal *DAL) GetSong(id uint) (*dao.Song, error) {
	s := dal.Db.GetSong(id)
	if s == nil {
		return nil, &dao.ErrNotFound{}
	}
	return s, nil
}

func (dal *DAL) GetAlbum(id uint) (*dao.Album, error) {
	s := dal.Db.GetAlbum(id)
	if s == nil {
		return nil, &dao.ErrNotFound{}
	}
	return s, nil
}

func (dal *DAL) GetArtist(id uint) (*dao.Artist, error) {
	a := dal.Db.GetArtist(id)
	if a == nil {
		return nil, &dao.ErrNotFound{}
	}
	return a, nil
}

func (dal *DAL) UpdatePlaylist(
	playlistID uint,
	name string,
	comment string,
	public *bool,
	addedSongs []uint,
	removedSongs []uint,
) error {

	p, err := dal.Db.GetPlaylist(playlistID)
	if err != nil {
		return err
	}
	if p == nil {
		return &dao.ErrNotFound{}
	}

	songIDs, err := dal.Db.GetPlaylistSongIds(playlistID)
	if err != nil {
		return err
	}

	var nameUpdate string
	if len(name) != 0 {
		nameUpdate = name
	} else {
		nameUpdate = p.Name
	}

	var commentUpdate string
	if len(comment) != 0 {
		commentUpdate = comment
	} else {
		commentUpdate = p.Comment
	}

	for _, index := range removedSongs {
		songIDs = append(songIDs[:index], songIDs[index+1:]...)
	}

	for _, songID := range addedSongs {
		songIDs = append(songIDs, songID)
	}

	dal.Db.ReplacePlaylistEntries(playlistID, songIDs)
	dal.Db.UpdatePlaylist(playlistID, nameUpdate, commentUpdate)
	return nil
}

func (dal *DAL) DeletePlaylist(id uint) error {
	return dal.Db.DeletePlaylist(id)
}

func (dal *DAL) GetAlbums(listType dao.AlbumList2Type, size uint, offset uint) []dao.Album {
	return dal.Db.GetAlbums(listType, size, offset)
}

func (d *DAL) DeleteMissing(tokens []string, providerID string) {
	d.Db.DeleteMissing(tokens, providerID)
}

// GetArtPath checks that an artwork file exists for the given ID and returns
// the full path if so.
func (dal *DAL) GetArtPath(id string) (string, error) {
	p := path.Join(dal.artDir, "art", id)
	_, err := os.Stat(p)
	if err != nil {
		return "", err
	}
	return p, nil
}

func (dal *DAL) StarSong(songID uint, star bool) error {
	return nil
}

func (dal *DAL) StarAlbum(albumID uint, star bool) error {
	return nil
}

func (dal *DAL) StarArtist(artistID uint, star bool) error {
	return nil
}
