package dal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cihub/seelog"
	"github.com/google/uuid"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/database"
	"github.com/hednowley/sound/entities"
	"github.com/hednowley/sound/hasher"
	"github.com/hednowley/sound/services"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DAL (data access layer) allows querying and writing application data.
type DAL struct {
	Db        *database.Default
	artDir    string
	artSizes  []uint
	resizeArt bool
}

// NewDAL constructs a new DAL.
func NewDAL(config *config.Config, database *database.Default) *DAL {
	return &DAL{
		Db:        database,
		artDir:    config.ArtPath,
		artSizes:  config.ArtSizes,
		resizeArt: config.ResizeArt,
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

// GetArtPath gets the full path of the artwork file with the given filename
// scaled to the given size (where 0 means unscaled).
// This file may or may not exist!
func (dal *DAL) GetArtPath(filename string, size uint) string {

	if size != 0 {
		ext := filepath.Ext(filename)
		return filepath.Join(
			dal.artDir,
			fmt.Sprintf("%v_%v.jpg", strings.TrimSuffix(filename, ext), size))
	}

	return path.Join(dal.artDir, filename)
}

// GetArt returns the path of the art file.
// Depending on the configuration, this file may be scaled to the requested
// dimension.
func (dal *DAL) GetArt(filename string, size uint) *string {

	scaled := dal.GetArtPath(filename, size)
	_, err := os.Stat(scaled)
	if err == nil {
		// Artwork exists at exactly the requested size 🎉
		return &scaled
	}

	if size == 0 {
		// We couldn't find the original artwork even though it was asked for 🤷‍♂️
		return nil
	}

	if dal.resizeArt {
		p := dal.GetArtPath(filename, size)

		// Resize on demand
		err = services.Resize(dal.GetArtPath(filename, 0), p, size)
		if err == nil {
			return &p
		}

		seelog.Errorf("Error resizing %v", p)
		return nil
	}

	closest := uint(0)
	for _, s := range dal.artSizes {
		if s >= size && (closest == 0 || closest > s) {
			closest = s
		}
	}

	p := dal.GetArtPath(filename, closest)
	_, err = os.Stat(p)
	if err == nil {
		return &p
	}

	// Make a scaled copy of the image as this size is inside our presets
	err = services.Resize(dal.GetArtPath(filename, 0), p, closest)
	if err == nil {
		return &p
	}

	seelog.Errorf("Error resizing %v", p)
	return nil

}

func (dal *DAL) PutArt(conn *pgxpool.Conn, art *entities.CoverArtData) (*dao.Art, error) {
	hash := hasher.GetHash(art.Raw)
	a := dal.Db.GetArtFromHash(conn, hash)

	var filename string
	var artPath string

	if a == nil {
		// Save art to disk
		filename = fmt.Sprintf("%v.%v", uuid.New().String(), art.Extension)
		artPath = dal.GetArtPath(filename, 0)

		err := ioutil.WriteFile(artPath, art.Raw, 0644)
		if err != nil {
			seelog.Errorf("Error saving artwork %v", artPath)
			return nil, err
		}

		// Save the record with the new path
		a, err = dal.Db.InsertArt(conn, filename, hash)
		if err != nil {
			return a, err
		}
	} else {
		filename = a.Path
		artPath = dal.GetArtPath(a.Path, 0)
	}

	// Make resized versions of the art
	for _, size := range dal.artSizes {
		resized := dal.GetArtPath(filename, size)

		// Check if resized file already exists on disk
		_, err := os.Stat(resized)
		if err == nil {
			continue
		}

		// Ignore resizing errors
		err = services.Resize(artPath, resized, size)
		if err != nil {
			seelog.Errorf("Error resizing %v", resized)
		}
	}

	return a, nil
}

func (dal *DAL) UpdatePlaylist(
	conn *pgxpool.Conn,
	playlistID uint,
	name string,
	comment string,
	public *bool,
	addedSongs []uint,
	removedSongs []uint,
	requestor string,
) error {

	p, err := dal.Db.GetPlaylist(conn, playlistID, requestor)
	if err != nil {
		return err
	}
	if p == nil {
		return &dao.ErrNotFound{}
	}

	songIDs, err := dal.Db.GetPlaylistSongIds(conn, playlistID, requestor)
	if err != nil {
		return err
	}

	var nameUpdate string
	if len(name) != 0 {
		nameUpdate = name
	} else {
		nameUpdate = p.Name
	}

	var publicUpdate bool
	if public != nil {
		publicUpdate = *public
	} else {
		publicUpdate = p.Public
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

	err = dal.Db.ReplacePlaylistEntries(conn, playlistID, songIDs)
	if err != nil {
		return err
	}

	_, err = dal.Db.UpdatePlaylist(conn, playlistID, nameUpdate, publicUpdate, commentUpdate)
	if err != nil {
		return err
	}

	return nil
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
