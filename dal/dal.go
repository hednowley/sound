package dal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"time"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/database"
	"github.com/hednowley/sound/entities"
	"github.com/hednowley/sound/hasher"
	"github.com/hednowley/sound/interfaces"
)

// DAL (data access layer) allows high-level manipulation of application data.
type DAL struct {
	db     *database.Default
	artDir string
	resize bool
}

// NewDAL constructs a new DAL.
func NewDAL(config *config.Config, database *database.Default) interfaces.DAL {
	return &DAL{
		db:     database,
		artDir: config.ArtPath,
		resize: config.ResizeArt,
	}
}

// PutSong updates the stored song with the given path and provider ID and returns its ID.
// If there is no such song then a new one is created and its ID is returned.
// The associated album, artist, artwork and genre are created too if necessary.
func (dal *DAL) PutSong(song *dao.Song, data *entities.FileInfo) *dao.Song {

	genre := dal.db.PutGenreByName(data.Genre)
	art := dal.PutArt(data.CoverArt)
	album := dal.db.PutAlbumByAttributes(data.Album, data.AlbumArtist, data.Disambiguator)

	song.Path = data.Path
	song.AlbumID = album.ID
	song.Title = data.Title
	song.Track = data.Track
	song.Disc = data.Disc
	song.GenreID = genre.ID
	song.Year = data.Year
	song.Extension = data.Extension
	song.Size = data.Size
	song.Bitrate = data.Bitrate
	song.Duration = data.Duration

	if art != nil {
		song.Art = art.Path
	}

	dal.db.PutSong(song)
	return song
}

func (dal *DAL) PutArt(art *entities.CoverArtData) *dao.Art {

	if art == nil {
		return nil
	}

	hash := hasher.GetHash(art.Raw)
	a := dal.db.GetArtFromHash(hash)
	if a != nil {
		// Artwork already exists
		return a
	}

	// Save the hash to get a fresh ID
	a = &dao.Art{
		Hash: hash,
	}
	dal.db.PutArt(a)

	a = &dao.Art{
		Path: fmt.Sprintf("%v.%v", a.ID, art.Extension),
		Hash: hash,
	}

	err := ioutil.WriteFile(path.Join(dal.artDir, "art", a.Path), art.Raw, 0644)
	if err != nil {
		seelog.Errorf("Error saving artwork %v", a.Path)
	}

	// Save the record with the new path
	dal.db.PutArt(a)
	return a
}

// SynchroniseAlbum updates the given album's aggregate properties, e.g. duration.
func (dal *DAL) SynchroniseAlbum(id uint) (*dao.Album, error) {

	seelog.Infof("Synchronising album %v", id)

	a, err := dal.GetAlbum(id, false, false, true)
	if err != nil {
		return nil, err
	}

	artSet := false
	genreSet := false
	yearSet := false
	duration := 0

	for _, song := range a.Songs {

		duration = duration + song.Duration

		if !artSet && song.Art != "" {
			a.Art = song.Art
			artSet = true
		}
		if !genreSet && song.GenreID != 0 {
			a.GenreID = song.GenreID
			genreSet = true
		}
		if !yearSet && song.Year != 0 {
			a.Year = song.Year
			yearSet = true
		}
	}

	a.Duration = duration
	dal.db.PutAlbum(a)
	return a, nil
}

// SynchroniseAlbum updates the given artist's aggregate properties, e.g. album count.
func (dal *DAL) SynchroniseArtist(id uint) error {

	seelog.Infof("Synchronising artist %v", id)

	a := dal.db.GetArtist(id, true, false)
	if a == nil {
		return &dao.ErrNotFound{}
	}

	artSet := false
	duration := 0

	for _, album := range a.Albums {
		if !artSet && album.Art != "" {
			a.Art = album.Art
			artSet = true
		}
		duration = duration + album.Duration
	}

	a.Duration = duration
	dal.db.PutArtist(a)
	return nil
}

func (dal *DAL) deleteSong(song *dao.Song) {

	for _, p := range dal.db.GetPlaylists() {

		// Indexes of entries to delete
		var d []int

		for i, e := range p.Entries {
			if e.SongID == song.ID {
				d = append(d, i)
			}
		}

		if len(d) > 0 {

		}

	}

	// Delete songs

	// Delete albums

	// Delete artists

	// Delete genres

	// Delete from playlists

	// Delete art
}

func (dal *DAL) PutPlaylist(id uint, name string, songIDs []uint) (uint, error) {

	now := time.Now()
	var p *dao.Playlist

	if id == 0 {
		p = &dao.Playlist{
			Name:    name,
			Created: &now,
			Changed: &now,
		}
		dal.db.AddPlaylist(p)
	} else {
		p = dal.db.GetPlaylist(id, false, false, false, false)
		if p == nil {
			return 0, &dao.ErrNotFound{}
		}

		p.Changed = &now

		if name != "" {
			p.Name = name
		}

		dal.db.UpdatePlaylist(p)
	}

	entries := []*dao.PlaylistEntry{}
	i := 0
	for _, songID := range songIDs {

		// We can't trust the song actually exists!
		_, err := dal.GetSong(songID, false, false, false)
		if err == nil {
			e := &dao.PlaylistEntry{
				PlaylistID: p.ID,
				Index:      i,
				SongID:     songID,
			}
			entries = append(entries, e)
			i++
		}
	}

	dal.db.ReplacePlaylistEntries(p, entries)
	return p.ID, nil
}

func (dal *DAL) GetSong(id uint, genre bool, album bool, artist bool) (*dao.Song, error) {
	s := dal.db.GetSong(id, genre, album, artist)
	if s == nil {
		return nil, &dao.ErrNotFound{}
	}
	return s, nil
}

func (dal *DAL) GetSongFromToken(token string, providerID string) *dao.Song {
	return dal.db.GetSongFromToken(token, providerID)
}

func (dal *DAL) GetAlbum(id uint, genre bool, artist bool, songs bool) (*dao.Album, error) {
	s := dal.db.GetAlbum(id, genre, artist, songs)
	if s == nil {
		return nil, &dao.ErrNotFound{}
	}
	return s, nil
}

func (dal *DAL) GetArtist(id uint) (*dao.Artist, error) {
	a := dal.db.GetArtist(id, true, true)
	if a == nil {
		return nil, &dao.ErrNotFound{}
	}
	return a, nil
}

func (dal *DAL) GetGenre(name string) (*dao.Genre, error) {
	g := dal.db.GetGenre(name)
	if g == nil {
		return nil, &dao.ErrNotFound{}
	}
	return g, nil
}

func (dal *DAL) GetPlaylist(id uint) (*dao.Playlist, error) {
	p := dal.db.GetPlaylist(id, true, true, true, true)
	if p == nil {
		return nil, &dao.ErrNotFound{}
	}

	// Sort entries
	s := NewPlaylistSorter(p)
	sort.Sort(s)

	return p, nil
}

func (dal *DAL) UpdatePlaylist(id uint, name string, comment string, public *bool, addedSongs []uint, removedSongs []uint) error {

	p := dal.db.GetPlaylist(id, true, false, false, false)
	if p == nil {
		return &dao.ErrNotFound{}
	}

	if len(name) != 0 {
		p.Name = name
	}

	if len(comment) != 0 {
		p.Comment = comment
	}

	if public != nil {
		p.Public = *public
	}

	e := p.Entries

	for _, index := range removedSongs {
		e = append(e[:index], e[index+1:]...)
	}

	for _, song := range addedSongs {
		entry := dao.PlaylistEntry{
			PlaylistID: p.ID,
			SongID:     song,
		}
		e = append(e, &entry)
	}

	// Reassign IDs
	for i := range e {
		e[i].Index = i
	}

	dal.db.ReplacePlaylistEntries(p, e)
	dal.db.UpdatePlaylist(p)
	return nil
}

func (dal *DAL) DeletePlaylist(id uint) error {
	return dal.db.DeletePlaylist(id)
}

func (dal *DAL) GetAlbums(listType dao.AlbumList2Type, size uint, offset uint) []*dao.Album {
	return dal.db.GetAlbums(listType, size, offset)
}

func (dal *DAL) GetArtists() []*dao.Artist {
	return dal.db.GetArtists()
}

func (dal *DAL) GetGenres() []*dao.Genre {
	return dal.db.GetGenres()
}

// GetPlaylists returns all playlists.
func (dal *DAL) GetPlaylists() []*dao.Playlist {
	playlists := dal.db.GetPlaylists()
	for _, p := range playlists {
		s := NewPlaylistSorter(p)
		sort.Sort(s)
	}
	return playlists
}

// Empty deletes all data.
func (d *DAL) Empty() {
	d.db.Empty()
}

func (d *DAL) DeleteMissing(tokens []string, providerID string) {
	d.db.DeleteMissing(tokens, providerID)
}

func (d *DAL) SearchArtists(query string, count uint, offset uint) []*dao.Artist {
	return d.db.SearchArtists(query, count, offset)
}

func (d *DAL) SearchSongs(query string, count uint, offset uint) []*dao.Song {
	return d.db.SearchSongs(query, count, offset)
}

func (d *DAL) SearchAlbums(query string, count uint, offset uint) []*dao.Album {
	return d.db.SearchAlbums(query, count, offset)
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
