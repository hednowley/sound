package database

import (
	"fmt"
	"time"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dao"
	"github.com/jinzhu/gorm"

	// Postgres driver for GORM
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Default provides access to the default application database.
type Default struct {
	db *gorm.DB
}

// NewDefault constructs a new default database.
func NewDefault(config *config.Config) (*Default, error) {
	db, err := gorm.Open("postgres", config.Db)
	db = db.Set("gorm:association_autoupdate", false).
		Set("gorm:association_autocreate", false)
	database := Default{db: db}

	// db.LogMode(true)

	db.AutoMigrate(dao.Song{})
	db.AutoMigrate(dao.Artist{})
	db.AutoMigrate(dao.Album{})
	db.AutoMigrate(dao.Art{})
	db.AutoMigrate(dao.Genre{})
	db.AutoMigrate(dao.Playlist{})
	db.AutoMigrate(dao.PlaylistEntry{})

	return &database, err
}

// Close the database.
func (db *Default) Close() {
	db.Close()
}

// PutArtist updates the artist record with the same ID,
// or creates one if no such record exists.
func (db *Default) PutArtist(artist *dao.Artist) {
	db.db.Save(&artist)
}

// PutArtistByName returns the ID of the artist record with the same name,
// or creates a new one and returns its ID if there is no such artist.
func (db *Default) putArtistByName(name string) *dao.Artist {
	var artist dao.Artist
	db.db.FirstOrCreate(&artist, dao.Artist{Name: name})
	return &artist
}

// PutAlbum updates the album record with the same ID,
// or creates a new one and returns its ID if there is no such album.
func (db *Default) PutAlbum(album *dao.Album) {
	db.db.Save(&album)
}

// PutAlbumByAttributes returns the ID of the album record with the same name and artist,
// or creates a new one and returns its ID if there is no such album.
func (db *Default) PutAlbumByAttributes(name string, artist string, disambiguator string) *dao.Album {

	album := dao.Album{
		Name:          name,
		ArtistID:      db.putArtistByName(artist).ID,
		Disambiguator: disambiguator,
	}

	if db.db.Where(&album).First(&album).RecordNotFound() {
		now := time.Now()
		album.Created = &now
		db.db.Create(&album)
	}

	return &album
}

// PutGenreByName returns the ID of the genre record with the same name,
// or creates a new one and returns its ID if there is no such genre.
func (db *Default) PutGenreByName(name string) *dao.Genre {
	var genre dao.Genre
	db.db.FirstOrCreate(&genre, dao.Genre{Name: name})
	return &genre
}

// PutSong updates the song record with the same ID
// or creates a new one if there is no such song.
func (db *Default) PutSong(song *dao.Song) {
	db.db.Save(&song)
}

// PutArt updates the art record with the same ID,
// or creates a new one and returns its ID if there is no such art.
func (db *Default) PutArt(art *dao.Art) {
	db.db.Save(&art)
}

func (db *Default) GetArtFromHash(hash string) *dao.Art {
	var a dao.Art
	if db.db.Where(dao.Art{Hash: hash}).First(&a).RecordNotFound() {
		return nil
	}
	return &a
}

// GetSongFromToken returns a pointer to the song with the given path and provider,
// or nil if one doesn't exist. Joined entities are not loaded.
func (db *Default) GetSongFromToken(token string, providerID string) *dao.Song {
	var f dao.Song
	if db.db.
		Where(&dao.Song{
			Token:      token,
			ProviderID: providerID,
		}).First(&f).RecordNotFound() {
		return nil
	}
	return &f
}

func (db *Default) ReplacePlaylistEntries(playlist *dao.Playlist, entries []*dao.PlaylistEntry) {
	db.db.Delete(&dao.PlaylistEntry{}, "playlist_id = ?", playlist.ID)
	for _, e := range entries {
		e.PlaylistID = playlist.ID
		db.db.Save(e)
	}
}

// Getters

func (db *Default) GetPlaylist(id uint, entries bool, songs bool, albums bool, artists bool) *dao.Playlist {
	var p dao.Playlist
	d := db.db
	if entries {
		d = d.Preload("Entries")
	}
	if songs {
		d = d.Preload("Entries.Song")
	}
	if albums {
		d = d.Preload("Entries.Song.Album")
	}
	if artists {
		d = d.Preload("Entries.Song.Album.Artist")
	}
	if d.Where(dao.Playlist{ID: id}).
		First(&p).
		RecordNotFound() {
		return nil
	}
	return &p
}

func (db *Default) GetAlbum(id uint, genre bool, artist bool, songs bool) *dao.Album {
	var a dao.Album
	d := db.db
	if genre {
		d = d.Preload("Genre")
	}
	if artist {
		d = d.Preload("Artist")
	}
	if songs {
		d = d.Preload("Songs")
	}
	if d.Where(dao.Album{ID: id}).
		First(&a).
		RecordNotFound() {
		return nil
	}
	for i := range a.Songs {
		a.Songs[i].Album = &a
	}
	return &a
}

func (db *Default) GetArt(id uint) *dao.Art {
	var f dao.Art
	if db.db.
		Where(dao.Art{ID: id}).
		First(&f).
		RecordNotFound() {
		return nil
	}
	return &f
}

func (db *Default) GetGenre(name string) *dao.Genre {
	var g dao.Genre
	if db.db.
		Preload("Songs").
		Preload("Songs.Genre").
		Preload("Songs.Album").
		Preload("Songs.Album.Artist").
		Where(dao.Genre{Name: name}).
		First(&g).
		RecordNotFound() {
		return nil
	}
	return &g
}

func (db *Default) GetArtist(id uint, albums bool, songs bool) *dao.Artist {
	var a dao.Artist
	d := db.db
	if albums {
		d = d.Preload("Albums")
	}
	if songs {
		d = d.Preload("Albums.Songs")
	}
	if d.Where(&dao.Artist{ID: id}).
		First(&a).
		RecordNotFound() {
		return nil
	}
	return &a
}

func (db *Default) GetSong(id uint, genre bool, album bool, artist bool) *dao.Song {
	var s dao.Song
	d := db.db
	if genre {
		d = d.Preload("Genre")
	}
	if album {
		d = d.Preload("Album")
	}
	if artist {
		d = d.Preload("Album.Artist")
	}
	if d.Where(dao.Song{ID: id}).
		First(&s).
		RecordNotFound() {
		return nil
	}
	return &s
}

// Collection getters

func (db *Default) GetGenres() []*dao.Genre {
	var genres []*dao.Genre
	db.db.
		Preload("Songs").
		Preload("Albums").
		Find(&genres)
	return genres
}

func (db *Default) GetAlbums(listType dao.AlbumList2Type, size uint, offset uint) []*dao.Album {

	var albums []*dao.Album
	query := db.db.Preload("Songs").Preload("Artist")

	switch listType {
	case dao.Random:
		query = query.Order(gorm.Expr("random()"))
	case dao.Newest:
		query = query.Order("created desc")
	case dao.Frequent:
		query = query
	case dao.Recent:
		query = query
	case dao.Starred:
		query = query
	case dao.AlphabeticalByName:
		query = query.Order("name asc")
	case dao.AlphabeticalByArtist:
		query = query.Joins("JOIN artists ON albums.artist_id = artists.id").Order("UPPER(artists.name) asc")
	case dao.ByYear:
		query = query
	case dao.ByGenre:
		query = query
	}

	query.Limit(size).Offset(offset).Find(&albums)
	return albums
}

func (db *Default) GetArtists() []*dao.Artist {
	var artists []*dao.Artist
	db.db.Preload("Albums").Find(&artists)
	return artists
}

func (db *Default) GetPlaylists() []*dao.Playlist {
	var playlists []*dao.Playlist
	db.db.Preload("Entries").Find(&playlists)
	return playlists
}

// Putters

func (db *Default) AddPlaylist(playlist *dao.Playlist) uint {
	db.db.Create(&playlist)
	return playlist.ID
}

func (db *Default) UpdatePlaylist(playlist *dao.Playlist) {
	db.db.Save(&playlist)
}

// Deleters

func (db *Default) DeletePlaylist(id uint) error {
	p := db.GetPlaylist(id, false, false, false, false)
	if p == nil {
		return &dao.ErrNotFound{}
	}

	db.db.Delete(&p)
	return nil
}

func (db *Default) Empty() {
	db.db.Delete(dao.PlaylistEntry{})
	db.db.Delete(dao.Playlist{})
	db.db.Delete(dao.Song{})
	db.db.Delete(dao.Album{})
	db.db.Delete(dao.Artist{})
	db.db.Delete(dao.Genre{})
	db.db.Delete(dao.Art{})
}

func (db *Default) DeleteMissing(tokens []string, providerID string) {
	db.db.Where("provider_id = ?", providerID).Not("token IN (?)", tokens).Delete(dao.Song{})
	db.db.Exec("DELETE FROM albums WHERE id in (SELECT albums.id FROM albums LEFT JOIN songs ON songs.album_id = albums.id GROUP BY albums.id HAVING count(songs.id) = 0)")
	db.db.Exec("DELETE FROM artists WHERE id in (SELECT artists.id FROM artists LEFT JOIN albums ON albums.artist_id = artists.id GROUP BY artists.id HAVING count(albums.id) = 0)")
	//db.db.Joins("JOIN albums ON albums.artist_id = artists.id").Group("artists.id").Having("count(albums.id) = 1").Delete(dao.Artist{})
}

func (db *Default) SearchAlbums(query string, count uint, offset uint) []*dao.Album {
	var albums []*dao.Album
	db.db.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", query)).Limit(count).Offset(offset).Find(&albums)
	return albums
}

func (db *Default) SearchArtists(query string, count uint, offset uint) []*dao.Artist {
	var artists []*dao.Artist
	db.db.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", query)).Limit(count).Offset(offset).Find(&artists)
	return artists
}

func (db *Default) SearchSongs(query string, count uint, offset uint) []*dao.Song {
	var songs []*dao.Song
	db.db.Preload("Genre").
		Preload("Album").
		Preload("Album.Artist").
		Where("title ILIKE ?", fmt.Sprintf("%%%s%%", query)).
		Limit(count).
		Offset(offset).
		Find(&songs)
	return songs
}
