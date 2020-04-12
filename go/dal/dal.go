package dal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/cihub/seelog"
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
func (dal *DAL) PutSong(song *dao.Song, data *entities.FileInfo) *dao.Song {

	genre := dal.Db.PutGenreByName(data.Genre)
	art := dal.PutArt(data.CoverArt)
	album := dal.Db.PutAlbumByAttributes(data.Album, data.AlbumArtist, data.Disambiguator)

	song.Path = data.Path
	song.Artist = data.Artist
	song.AlbumID = album.ID
	song.Title = data.Title
	song.Track = data.Track
	song.Disc = data.Disc
	song.GenreID = genre.ID
	song.Year = data.Year
	song.Size = data.Size
	song.Bitrate = data.Bitrate
	song.Duration = data.Duration

	song.AlbumName = album.Name
	song.AlbumArtistID = album.ArtistID
	song.GenreName = genre.Name

	if art != nil {
		song.Art = art.Path
	}

	dal.Db.PutSong(song)
	return song
}

func (dal *DAL) PutArt(art *entities.CoverArtData) *dao.Art {

	if art == nil {
		return nil
	}

	hash := hasher.GetHash(art.Raw)
	a := dal.Db.GetArtFromHash(hash)
	if a != nil {
		// Artwork already exists
		return a
	}

	// Save the hash to get a fresh ID
	a = &dao.Art{
		Hash: hash,
	}
	dal.Db.PutArt(a)

	a.Path = fmt.Sprintf("%v.%v", a.ID, art.Extension)

	err := ioutil.WriteFile(path.Join(dal.artDir, "art", a.Path), art.Raw, 0644)
	if err != nil {
		seelog.Errorf("Error saving artwork %v", a.Path)
	}

	// Save the record with the new path
	dal.Db.PutArt(a)
	return a
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
		dal.Db.AddPlaylist(p)
	} else {
		p = dal.Db.GetPlaylist(id, false, false)
		if p == nil {
			return 0, &dao.ErrNotFound{}
		}

		p.Changed = &now

		if name != "" {
			p.Name = name
		}

		dal.Db.UpdatePlaylist(p)
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

	dal.Db.ReplacePlaylistEntries(p, entries)
	return p.ID, nil
}

func (dal *DAL) GetSong(id uint, genre bool, album bool, artist bool) (*dao.Song, error) {
	s := dal.Db.GetSong(id, genre, album, artist)
	if s == nil {
		return nil, &dao.ErrNotFound{}
	}
	return s, nil
}

func (dal *DAL) GetSongFromToken(token string, providerID string) *dao.Song {
	return dal.Db.GetSongFromToken(token, providerID)
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

func (dal *DAL) GetGenre(name string) (*dao.Genre, error) {
	g := dal.Db.GetGenre(name)
	if g == nil {
		return nil, &dao.ErrNotFound{}
	}
	return g, nil
}

func (dal *DAL) GetPlaylist(id uint) (*dao.Playlist, error) {
	p := dal.Db.GetPlaylist(id, true, true)
	if p == nil {
		return nil, &dao.ErrNotFound{}
	}
	return p, nil
}

func (dal *DAL) UpdatePlaylist(id uint, name string, comment string, public *bool, addedSongs []uint, removedSongs []uint) error {

	p := dal.Db.GetPlaylist(id, true, false)
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

	dal.Db.ReplacePlaylistEntries(p, e)
	dal.Db.UpdatePlaylist(p)
	return nil
}

func (dal *DAL) DeletePlaylist(id uint) error {
	return dal.Db.DeletePlaylist(id)
}

func (dal *DAL) GetAlbums(listType dao.AlbumList2Type, size uint, offset uint) []dao.Album {
	return dal.Db.GetAlbums(listType, size, offset)
}

func (dal *DAL) GetArtists(includeAlbums bool) []*dao.Artist {
	return dal.Db.GetArtists(includeAlbums)
}

func (dal *DAL) GetGenres() []*dao.Genre {
	return dal.Db.GetGenres()
}

// GetPlaylists returns all playlists.
func (dal *DAL) GetPlaylists() []*dao.Playlist {
	return dal.Db.GetPlaylists()
}

// Empty deletes all data.
func (d *DAL) Empty() {
	d.Db.Empty()
}

func (d *DAL) DeleteMissing(tokens []string, providerID string) {
	d.Db.DeleteMissing(tokens, providerID)
}

func (d *DAL) SearchArtists(query string, count uint, offset uint) []*dao.Artist {
	return d.Db.SearchArtists(query, count, offset)
}

func (d *DAL) SearchSongs(query string, count uint, offset uint) []*dao.Song {
	return d.Db.SearchSongs(query, count, offset)
}

func (d *DAL) SearchAlbums(query string, count uint, offset uint) []*dao.Album {
	return d.Db.SearchAlbums(query, count, offset)
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

func (dal *DAL) GetRandomSongs(size uint, from uint, to uint, genre string) []*dao.Song {
	return dal.Db.GetRandomSongs(size, from, to, genre)
}

func (dal *DAL) StarSong(songID uint, star bool) error {
	return dal.Db.StarSong(songID, star)
}

func (dal *DAL) StarAlbum(albumID uint, star bool) error {
	return dal.Db.StarAlbum(albumID, star)
}

func (dal *DAL) StarArtist(artistID uint, star bool) error {
	return dal.Db.StarArtist(artistID, star)
}
