package provider

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/entities"

	// SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

// BeetsProvider scans music files indexed inside a beets database.
type BeetsProvider struct {
	beets      *sql.DB
	id         string
	fileCount  int64
	isScanning bool
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

// Iterate through all files in the collection, calling the provided callback synchronously on each.
func (p *BeetsProvider) Iterate(callback func(token string) error) error {
	p.isScanning = true
	p.fileCount = 0

	// Ordered so most recently added songs are scanned first.
	rows, err := p.beets.Query("SELECT id FROM items ORDER BY id DESC")
	if err != nil {
		p.isScanning = false
		return err
	}

	var id int
	var callbackErr error

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			p.isScanning = false
			return err
		}

		p.fileCount = p.fileCount + 1

		// Use the id as a unique token
		callbackErr = callback(strconv.Itoa(id))
		if callbackErr != nil {
			break
		}
	}

	p.isScanning = false
	return callbackErr
}

// GetInfo returns information about the file with the given token.
func (p *BeetsProvider) GetInfo(token string) (*entities.FileInfo, error) {

	id, err := strconv.Atoi(token)
	if err != nil {
		return nil, fmt.Errorf("Unknown token '%v'", token)
	}
	row := p.beets.QueryRow("SELECT path, title, artist, album, album_id, albumartist, genre, track, disc, original_year, year, bitrate, length FROM items WHERE ID = ?", id)

	var path string
	var title string
	var artist string
	var album string
	var albumID int
	var albumArtist string
	var genre string
	var track int
	var disc int
	var originalYear int
	var year int
	var bitrate int
	var duration float64

	err = row.Scan(&path, &title, &artist, &album, &albumID, &albumArtist, &genre, &track, &disc, &originalYear, &year, &bitrate, &duration)

	if err == nil {

		art, err := p.getArt(albumID)
		if err != nil {
			seelog.Errorf("Error getting art for '%v': %v", title, err)
			art = nil
		}

		// Prefer original year
		if originalYear != 0 {
			year = originalYear
		}

		// Prefer original year
		if originalYear != 0 {
			year = originalYear
		}

		return &entities.FileInfo{
			Path:          path,
			Title:         title,
			Album:         album,
			Artist:        artist,
			AlbumArtist:   albumArtist,
			Genre:         genre,
			Track:         track,
			Disc:          disc,
			Year:          year,
			Bitrate:       bitrate / 1000,
			Duration:      int(duration),
			CoverArt:      art,
			Disambiguator: strconv.Itoa(albumID),
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
