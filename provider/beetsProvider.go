package provider

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/entities"
	"github.com/hednowley/sound/hasher"

	// SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

// BeetsProvider scans music files indexed inside a beets database.
type BeetsProvider struct {
	beets      *sql.DB
	id         string
	fileCount  int64
	isScanning bool
	scanID     string
	lastAlbum  int
	lastArt    *entities.CoverArtData
}

// NewBeetsProvider constructs a new provider.
func NewBeetsProvider(id string, beetsPath string) (*BeetsProvider, error) {
	beets, err := sql.Open("sqlite3", beetsPath)
	if err != nil {
		return nil, err
	}

	return &BeetsProvider{
		id:         id,
		beets:      beets,
		fileCount:  0,
		isScanning: false,
	}, nil
}

// FileCount returns the number of files the provider has processed in the current scan.
func (p *BeetsProvider) FileCount() int64 {
	return p.fileCount
}

// IsScanning returns whether this provider is in the middle of a scan.
func (p *BeetsProvider) IsScanning() bool {
	return p.isScanning
}

// ID of this provider.
func (p *BeetsProvider) ID() string {
	return p.id
}

func (p *BeetsProvider) ScanID() string {
	return p.scanID
}

// Iterate through all files in the collection, calling the provided callback synchronously on each.
func (p *BeetsProvider) Iterate(callback func(path string)) error {
	p.isScanning = true
	p.fileCount = 0
	p.scanID = hasher.GetHashFromInt(time.Now().Unix())

	// This should probably look at id's not paths.
	// Would need a new "external_id" field in our DB.
	rows, err := p.beets.Query("SELECT path FROM items")
	if err != nil {
		return err
	}

	var pathBlob []byte

	for rows.Next() {
		err = rows.Scan(&pathBlob)
		if err != nil {
			return err
		}

		p.fileCount = p.fileCount + 1
		callback(string(pathBlob))
	}

	p.isScanning = false
	return nil
}

// GetInfo returns information about the file at the given path.
func (p *BeetsProvider) GetInfo(path string) (*entities.FileInfo, error) {

	row := p.beets.QueryRow("SELECT title, artist, album, album_id, albumartist, genre, track, disc, year, bitrate, length FROM items WHERE PATH = ?", []byte(path))

	var title string
	var artist string
	var album string
	var albumID int
	var albumArtist string
	var genre string
	var track int
	var disc int
	var year int
	var bitrate int
	var duration float64

	err := row.Scan(&title, &artist, &album, &albumID, &albumArtist, &genre, &track, &disc, &year, &bitrate, &duration)

	if err == nil {

		art, err := p.getArt(albumID)
		if err != nil {
			seelog.Errorf("Error getting art for '%v': %v", title, err)
			art = nil
		}

		return &entities.FileInfo{
			Path:        path,
			Title:       title,
			Album:       album,
			Artist:      artist,
			AlbumArtist: albumArtist,
			Genre:       genre,
			Track:       track,
			Disc:        disc,
			Year:        year,
			Bitrate:     bitrate / 1000,
			Duration:    int(duration),
			CoverArt:    art,
		}, nil
	}
	if err == sql.ErrNoRows {
		return nil, errors.New("song not found in database")
	}

	return nil, err
}

func (p *BeetsProvider) getArt(album int) (*entities.CoverArtData, error) {

	if album == p.lastAlbum {
		return p.lastArt, nil
	}

	p.lastAlbum = album
	var pathBlob []byte
	var art *entities.CoverArtData

	row := p.beets.QueryRow("SELECT artpath FROM albums WHERE id = ?", album)
	err := row.Scan(&pathBlob)

	if err == nil {
		if len(pathBlob) == 0 {
			goto end
		}

		path := string(pathBlob)
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			goto end
		}

		art = &entities.CoverArtData{
			Raw:       bytes,
			Extension: strings.TrimPrefix(filepath.Ext(path), "."),
		}
	}
	if err == sql.ErrNoRows {
		err = errors.New("album not found in database")
	}

end:
	// Cache the art (this will speed up requests if they arrive grouped by album)
	p.lastArt = art
	return art, err
}
