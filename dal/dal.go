package dal

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"sync"
	"time"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/database"
	"github.com/hednowley/sound/entities"
	"github.com/hednowley/sound/hasher"
	"github.com/hednowley/sound/provider"
)

// DAL (data access layer) allow high-level manipulation of application data.
type DAL struct {
	db        *database.Default
	providers []provider.Provider
	artDir    string
}

// NewDAL constructs a new DAL.
func NewDAL(providers []provider.Provider, config *config.Config, database *database.Default) *DAL {
	return &DAL{
		db:        database,
		artDir:    config.ArtPath,
		providers: providers,
	}
}

// putSong updates the stored song with the given path and provider ID and returns its ID.putSong
// If there is no such song then a new one is created and its ID is returned.
// The associated album, artist, artwork and genre are created too if necessary.
func (dal *DAL) putSong(song *dao.Song, data *entities.FileInfo) *dao.Song {

	genre := dal.db.PutGenreByName(data.Genre)
	art := dal.putArt(data.CoverArt)
	album := dal.db.PutAlbumByNameAndArtist(data.Album, data.AlbumArtist)

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
		song.ArtID = art.ID
	}

	dal.db.PutSong(song)
	return song
}

func (dal *DAL) putArt(art *entities.CoverArtData) *dao.Art {

	if art == nil {
		return nil
	}

	hash := hasher.GetHash(art.Raw)
	a := dal.db.GetArtFromHash(hash)
	if a != nil {
		return a
	}

	a = &dao.Art{}
	dal.db.PutArt(a)

	p := path.Join(dal.artDir, "art", fmt.Sprintf("%v.%v", a.ID, art.Extension))
	a = &dao.Art{
		Path: p,
		Hash: hash,
	}

	err := ioutil.WriteFile(a.Path, art.Raw, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}

	dal.db.PutArt(a)
	return a
}

// Should "stagger" this to run every 50 songs
func (dal *DAL) synchroniseAlbum(id uint) (*dao.Album, error) {

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

		if !artSet && song.ArtID != 0 {
			a.ArtID = song.ArtID
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

func (dal *DAL) synchroniseArtist(id uint) error {

	seelog.Infof("Synchronising artist %v", id)

	a := dal.db.GetArtist(id, true, false)
	if a == nil {
		return &dao.ErrNotFound{}
	}

	artSet := false
	duration := 0

	for _, album := range a.Albums {
		if !artSet && album.ArtID != 0 {
			a.ArtID = album.ArtID
			artSet = true
		}
		duration = duration + album.Duration
	}

	a.Duration = duration
	dal.db.PutArtist(a)
	return nil
}

func (dal *DAL) updateSongScanID(song *dao.Song, scanID string) {
	song.ScanID = scanID
	dal.db.PutSong(song)
}

func (dal *DAL) DeleteMissingSongs(scanID string, scannerID string) {
	// Delete songs

	// Delete albums

	// Delete artists

	// Delete genres

	// Delete from playlists

	// Delete art
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
		_, err := dal.GetSong(songID, false, false, false, false)
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

func (dal *DAL) GetSong(id uint, genre bool, album bool, artist bool, art bool) (*dao.Song, error) {
	s := dal.db.GetSong(id, genre, album, artist, art)
	if s == nil {
		return nil, &dao.ErrNotFound{}
	}
	return s, nil
}

func (dal *DAL) GetAlbum(id uint, genre bool, artist bool, songs bool) (*dao.Album, error) {
	s := dal.db.GetAlbum(id, genre, artist, songs)
	if s == nil {
		return nil, &dao.ErrNotFound{}
	}
	return s, nil
}

func (dal *DAL) GetArt(id uint) (*dao.Art, error) {
	a := dal.db.GetArt(id)
	if a == nil {
		return nil, &dao.ErrNotFound{}
	}
	return a, nil
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

func (dal *DAL) GetScanFileCount() int64 {
	count := int64(0)
	for _, p := range dal.providers {
		count += p.FileCount()
	}
	return count
}

func (dal *DAL) GetScanStatus() bool {
	scanning := false
	for _, p := range dal.providers {
		scanning = scanning || p.IsScanning()
	}
	return scanning
}

// StartAllScans asks all providers to start scanning in parallel.
func (dal *DAL) StartAllScans(update bool, delete bool) {
	seelog.Info("Starting all scans.")
	var wg sync.WaitGroup
	for _, p := range dal.providers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dal.startScan(p, update, delete)
		}()
	}
	wg.Wait()
}

func (dal *DAL) startScan(provider provider.Provider, update bool, delete bool) {
	providerID := provider.ID()

	if provider.IsScanning() {
		seelog.Infof("Skipped '%v' scan as one is already in progress.", providerID)
		return
	}

	seelog.Infof("Started '%v' scan.", providerID)
	scanID := provider.ScanID()
	synch := NewSynchroniser(dal, 10)

	err := provider.Iterate(func(token string) {
		s := dal.db.GetSongFromToken(token, providerID)
		if s == nil || update {
			data, err2 := provider.GetInfo(token)
			if err2 != nil {
				seelog.Errorf("Cannot read music info for '%v': %v", token, err2)
				return
			}

			if s == nil {
				seelog.Infof("Adding token '%v'", token)
				now := time.Now()
				s = &dao.Song{
					Created:    &now,
					ProviderID: providerID,
					Token:      token,
				}
			} else {
				seelog.Infof("Updating token '%v'", token)
				synch.Notify(s.AlbumID) // Notify of potential change to old album
			}

			s.ScanID = scanID
			s = dal.putSong(s, data)

			// Notify of change to new album
			synch.Notify(s.AlbumID)

		} else {
			seelog.Infof("Skipping token '%v'", token)
			dal.updateSongScanID(s, scanID)
		}
	})
	if err != nil {
		seelog.Errorf("Error during '%v' scan: %v", providerID, err)
	}

	// Make any remaining updates
	synch.Flush()

	seelog.Infof("Finished '%v' scan.", providerID)
}

// Empty deletes all data.
func (d *DAL) Empty() {
	d.db.Empty()
}
