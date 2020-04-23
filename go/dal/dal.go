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
	"github.com/jackc/pgx/v4/pgxpool"
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

func (dal *DAL) PutSong(conn *pgxpool.Conn, fileInfo *entities.FileInfo, token string, providerID string, songId *uint) error {

	genreID, err := dal.Db.PutGenreByName(conn, fileInfo.Genre)
	if err != nil {
		return err
	}

	var art *dao.Art
	if fileInfo.CoverArt != nil {
		art, err = dal.PutArt(conn, fileInfo.CoverArt)
		if err != nil {
			return err
		}
	}

	albumID, err := dal.Db.PutAlbumByAttributes(conn,
		fileInfo.Album, fileInfo.AlbumArtist, fileInfo.Disambiguator)
	if err != nil {
		return err
	}

	var artPath *string
	if art != nil {
		artPath = &art.Path
	}

	if songId == nil {
		_, err = dal.Db.InsertSong(
			conn,
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

	return dal.Db.UpdateSong(
		conn,
		*songId,
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
		false,
	)
}

func (dal *DAL) PutArt(conn *pgxpool.Conn, art *entities.CoverArtData) (*dao.Art, error) {
	hash := hasher.GetHash(art.Raw)
	a := dal.Db.GetArtFromHash(conn, hash)
	if a != nil {
		// Artwork already exists
		return a, nil
	}

	filePath := fmt.Sprintf("%v.%v", uuid.New().String(), art.Extension)

	err := ioutil.WriteFile(path.Join(dal.artDir, filePath), art.Raw, 0644)
	if err != nil {
		seelog.Errorf("Error saving artwork %v", filePath)
		return nil, err
	}

	// Save the record with the new path
	return dal.Db.InsertArt(conn, filePath, hash)
}

func (dal *DAL) PutPlaylist(conn *pgxpool.Conn, id uint, name string, songIDs []uint, public bool) (uint, error) {

	now := time.Now()
	var playlistID uint

	if id == 0 {
		inserted, err := dal.Db.InsertPlaylist(conn, name, "", public)
		if err != nil {
			return 0, err
		}
		playlistID = inserted
	} else {
		playlist, err := dal.Db.GetPlaylist(conn, id)
		if err != nil {
			return 0, err
		}

		playlist.Changed = &now

		var nameUpdate string
		if name == "" {
			nameUpdate = playlist.Name
		} else {
			nameUpdate = name
		}

		playlist, err = dal.Db.UpdatePlaylist(conn, playlist.ID, nameUpdate, playlist.Comment)
		if err != nil {
			return 0, err
		}

		playlistID = playlist.ID
	}

	err := dal.Db.ReplacePlaylistEntries(conn, playlistID, songIDs)
	if err != nil {
		return 0, err
	}

	return playlistID, nil
}

func (dal *DAL) UpdatePlaylist(
	conn *pgxpool.Conn,
	playlistID uint,
	name string,
	comment string,
	public *bool,
	addedSongs []uint,
	removedSongs []uint,
) error {

	p, err := dal.Db.GetPlaylist(conn, playlistID)
	if err != nil {
		return err
	}
	if p == nil {
		return &dao.ErrNotFound{}
	}

	songIDs, err := dal.Db.GetPlaylistSongIds(conn, playlistID)
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

	dal.Db.ReplacePlaylistEntries(conn, playlistID, songIDs)
	dal.Db.UpdatePlaylist(conn, playlistID, nameUpdate, commentUpdate)
	return nil
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
