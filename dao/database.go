package dao

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/entities"
	"github.com/hednowley/sound/hasher"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	db     *gorm.DB
	artDir string
}

func NewDatabase(config *config.Config) (*Database, error) {
	db, err := gorm.Open("postgres", config.Db)
	database := Database{
		db:     db,
		artDir: config.ArtPath,
	}

	// db.LogMode(true)

	db.AutoMigrate(Song{})
	db.AutoMigrate(Artist{})
	db.AutoMigrate(Album{})
	db.AutoMigrate(Art{})
	db.AutoMigrate(Genre{})
	db.AutoMigrate(Playlist{})
	db.AutoMigrate(PlaylistEntry{})

	return &database, err
}

func (db *Database) Close() {
	db.Close()
}

func (db *Database) putArt(art *entities.CoverArtData) uint {

	if art == nil {
		return 0
	}

	hash := hasher.GetHash(art.Raw)

	a := Art{Hash: hash}
	if db.db.Where(&a).First(&a).RecordNotFound() {
		db.db.Create(&a)
		a.Path = path.Join(db.artDir, "art", fmt.Sprintf("%v.%v", a.ID, art.Extension))
		db.db.Save(&a)
		err := ioutil.WriteFile(a.Path, art.Raw, 0644)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return a.ID
}

func (db *Database) PutSong(data *entities.FileData) uint {

	genreID := db.putGenre(data.Genre)
	artID := db.putArt(data.CoverArt)
	albumID := db.putAlbum(data)

	song := Song{
		Path: data.Path,
	}

	result := db.db.Where(&song).First(&song)

	song.AlbumID = albumID
	song.Title = data.Title
	song.Track = data.Track
	song.Disc = data.Disc
	song.GenreID = genreID
	song.Year = data.Year
	song.ArtID = artID
	song.Extension = data.Extension
	song.Size = data.Size

	if result.RecordNotFound() {
		now := time.Now()
		song.Created = &now
		db.db.Create(&song)
	} else {
		db.db.Save(&song)
	}

	db.generateAlbumData(albumID)

	return song.ID
}

func (db *Database) putGenre(name string) uint {
	var genre Genre
	db.db.FirstOrCreate(&genre, Genre{Name: name})
	return genre.ID
}

func (db *Database) putArtist(name string) uint {

	var artist Artist
	db.db.FirstOrCreate(&artist, Artist{Name: name})
	return artist.ID
}

func (db *Database) putAlbum(data *entities.FileData) uint {

	album := Album{
		Name:     data.Album,
		ArtistID: db.putArtist(data.AlbumArtist),
	}

	if db.db.Where(&album).First(&album).RecordNotFound() {
		now := time.Now()
		album.Created = &now
		db.db.Create(&album)
	}

	return album.ID
}

func (db *Database) generateAlbumData(id uint) error {
	var a Album
	if db.db.
		Preload("Songs").
		Preload("Songs.Genre").
		Where(&Album{
			ID: id,
		}).First(&a).RecordNotFound() {
		return errors.New("No such album")
	}

	artSet := false
	genreSet := false
	yearSet := false

	for _, song := range a.Songs {

		if artSet && genreSet && yearSet {
			break
		}

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

	db.db.Save(&a)

	db.generateArtistData(a.ArtistID)

	return nil
}

func (db *Database) generateArtistData(id uint) error {
	var a Artist
	if db.db.Preload("Albums").Where(&Artist{
		ID: id,
	}).First(&a).RecordNotFound() {
		return errors.New("No such artist")
	}

	for _, album := range a.Albums {
		if album.ArtID != 0 {
			a.ArtID = album.ArtID
			break
		}
	}

	db.db.Save(&a)
	return nil
}

type AlbumList2Type int

const (
	Random               AlbumList2Type = 0
	Newest               AlbumList2Type = 1
	Frequent             AlbumList2Type = 2
	Recent               AlbumList2Type = 3
	Starred              AlbumList2Type = 4
	AlphabeticalByName   AlbumList2Type = 5
	AlphabeticalByArtist AlbumList2Type = 6
	ByYear               AlbumList2Type = 7
	ByGenre              AlbumList2Type = 8
)

func (db *Database) GetAlbums(listType AlbumList2Type, size uint, offset uint) []*Album {

	var albums []*Album
	query := db.db.Preload("Songs").Preload("Artist")

	switch listType {
	case Random:
		query = query.Order(gorm.Expr("random()"))
	case Newest:
		query = query.Order("created desc")
	case Frequent:
		query = query
	case Recent:
		query = query
	case Starred:
		query = query
	case AlphabeticalByName:
		query = query.Order("name asc")
	case AlphabeticalByArtist:
		query = query.Joins("JOIN artists ON albums.artist_id = artists.id").Order("UPPER(artists.name) asc")
	case ByYear:
		query = query
	case ByGenre:
		query = query
	}

	query.Limit(size).Offset(offset).Find(&albums)
	return albums
}

func (db *Database) GetArtists() []*Artist {
	var artists []*Artist
	db.db.Preload("Albums").Find(&artists)
	return artists
}

func (db *Database) GetArt(id uint) (*Art, error) {
	var f Art
	if db.db.Where(&Art{
		ID: id,
	}).First(&f).RecordNotFound() {
		return &f, errors.New("No such art")
	}
	return &f, nil
}

func (db *Database) GetSong(id uint) (*Song, error) {
	var f Song
	if db.db.
		Preload("Genre").
		Preload("Album").
		Preload("Album.Artist").
		Preload("Art").
		Where(&Song{
			ID: id,
		}).First(&f).RecordNotFound() {
		return &f, errors.New("No such song")
	}
	return &f, nil
}

func (db *Database) GetArtist(id uint) (*Artist, error) {
	var a Artist
	if db.db.Preload("Albums").Preload("Albums.Songs").Where(&Artist{
		ID: id,
	}).First(&a).RecordNotFound() {
		return &a, &ErrNotFound{}
	}
	return &a, nil
}

func (db *Database) GetAlbum(id uint) (*Album, error) {
	var a Album
	if db.db.
		Preload("Genre").
		Preload("Artist").
		Preload("Songs").
		Where(&Album{
			ID: id,
		}).First(&a).RecordNotFound() {
		return &a, &ErrNotFound{}
	}
	for i := range a.Songs {
		a.Songs[i].Album = &a
	}
	return &a, nil
}

func (db *Database) GetGenres() []*Genre {
	var genres []*Genre
	db.db.
		Preload("Songs").
		Preload("Albums").
		Find(&genres)
	return genres
}

func (db *Database) GetGenre(name string) (*Genre, error) {
	var g Genre
	if db.db.
		Preload("Songs").
		Preload("Songs.Genre").
		Preload("Songs.Album").
		Preload("Songs.Album.Artist").
		Where(&Genre{
			Name: name,
		}).First(&g).RecordNotFound() {
		return nil, &ErrNotFound{}
	}
	return &g, nil
}

func (db *Database) PutPlaylist(id uint, name string, songIDs []uint) (uint, error) {

	now := time.Now()
	var playlist Playlist

	if id == 0 {
		playlist.Name = name
		playlist.Created = &now
		playlist.Changed = &now
		db.db.Create(&playlist)
	} else {
		if db.db.Where(Playlist{ID: id}).
			First(&playlist).
			RecordNotFound() {
			return 0, &ErrNotFound{}
		}

		if name != "" {
			playlist.Name = name
		}
		playlist.Changed = &now
		db.db.Save(&playlist)

	}

	entries := []PlaylistEntry{}
	i := 0
	for _, songID := range songIDs {

		// We can't trust the song actually exists!
		_, err := db.GetSong(songID)
		if err == nil {
			p := PlaylistEntry{
				PlaylistID: playlist.ID,
				Index:      i,
				SongID:     songID,
			}
			entries = append(entries, p)
		}
		i++
	}

	db.db.Model(&playlist).Association("Entries").Replace(&entries)
	db.db.Delete(PlaylistEntry{}, "playlist_id IS NULL")
	return playlist.ID, nil
}

func (db *Database) GetPlaylists() []*Playlist {
	var playlists []*Playlist
	db.db.Preload("Entries").Find(&playlists)
	return playlists
}

func (db *Database) DeletePlaylist(id uint) error {
	var p Playlist
	if db.db.Where(&Playlist{
		ID: id,
	}).First(&p).RecordNotFound() {
		return &ErrNotFound{}
	}

	db.db.Delete(&p)
	return nil
}

func (db *Database) GetPlaylist(id uint) (*Playlist, error) {
	var a Playlist
	if db.db.
		Preload("Entries").
		Preload("Entries.Song").
		Preload("Entries.Song.Album").
		Preload("Entries.Song.Album.Artist").
		Where(&Playlist{
			ID: id,
		}).First(&a).RecordNotFound() {
		return &a, errors.New("No such playlist")
	}
	return &a, nil
}

func (db *Database) UpdatePlaylist(id uint, name string, comment string, public *bool, addedSongs []uint, removedSongs []uint) error {
	var p Playlist
	if db.db.Where(&Playlist{
		ID: id,
	}).First(&p).RecordNotFound() {
		return errors.New("No such playlist")
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

	for _, index := range removedSongs {
		p.Entries = append(p.Entries[:index], p.Entries[index+1:]...)
	}

	for _, song := range addedSongs {
		entry := PlaylistEntry{
			PlaylistID: p.ID,
			SongID:     song,
		}
		p.Entries = append(p.Entries, &entry)
	}

	// Reassign IDs
	for i := range p.Entries {
		p.Entries[i].Index = i
	}

	db.db.Save(&p)
	return nil
}
