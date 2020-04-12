package database

import (
	"context"
	"fmt"
	"time"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dao"
	"github.com/jinzhu/gorm"

	// Postgres driver for GORM
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jackc/pgx/v4"
)

// Default provides access to the default application database.
type Default struct {
	db   *gorm.DB
	conn *pgx.Conn
}

// NewDefault constructs a new default database.
func NewDefault(config *config.Config) (*Default, error) {
	db, err := gorm.Open("postgres", config.Db)
	db = db.Set("gorm:association_autoupdate", false).
		Set("gorm:association_autocreate", false)

	conn, err := pgx.Connect(context.Background(), config.Db)
	if err != nil {
		return nil, err
	}

	database := Default{db: db, conn: conn}

	//db.LogMode(true)

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
	db.db.Where("name ILIKE ?", name).Limit(1).Find(&artist)
	if artist.ID != 0 {
		return &artist
	}

	artist.Name = name
	db.db.Create(&artist)
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

func (db *Default) GetPlaylist(id uint, entries bool, songs bool) *dao.Playlist {
	var p dao.Playlist
	d := db.db
	if entries {
		d = d.Preload("Entries")
	}
	if songs {
		d = d.Preload("Entries.Song")
	}
	if d.Where(dao.Playlist{ID: id}).
		First(&p).
		RecordNotFound() {
		return nil
	}
	return &p
}

func (db *Default) GetAlbum(albumID uint) *dao.Album {

	var a dao.Album

	err := db.conn.QueryRow(context.Background(),
		`
		SELECT 
			albums.id,
			albums.artist_id,
			albums.name,
			albums.created,
			albums.art,
			albums.genre_id,
			albums.year,
			albums.disambiguator,
			albums.starred,
			artists.name,
			genres.name,
			COUNT(songs.id),
			SUM(songs.duration)

		FROM
			albums
		LEFT JOIN 
			artists
		ON
			artists.id = albums.artist_id
		LEFT JOIN 
			genres
		ON
			genres.id = albums.genre_id
		LEFT JOIN 
			songs
		ON
			songs.album_id = albums.id
		WHERE
			albums.id = $1
		GROUP BY
			albums.id, artists.id, genres.id
	`, albumID,
	).Scan(
		&a.ID,
		&a.ArtistID,
		&a.Name,
		&a.Created,
		&a.Art,
		&a.GenreID,
		&a.Year,
		&a.Disambiguator,
		&a.Starred,
		&a.ArtistName,
		&a.GenreName,
		&a.SongCount,
		&a.Duration,
	)

	if err != nil {
		return nil
	}

	return &a
}

func (db *Default) GetAlbumSongs(albumID uint) []dao.Song {

	rows, _ := db.conn.Query(context.Background(),
		`
		SELECT 
			songs.id,
			songs.artist,
			songs.album_id,
			songs.path,
			songs.title,
			songs.track,
			songs.disc,
			songs.genre_id,
			songs.year,
			songs.art,
			songs.created,
			songs.size,
			songs.bitrate,
			songs.duration,
			songs.token,
			songs.provider_id,
			songs.starred,
			albums.name,
			albums.artist_id,
			genres.name
		FROM
			songs
		LEFT JOIN 
			albums
		ON
			albums.id = songs.album_id
		LEFT JOIN 
			genres
		ON
			genres.id = songs.genre_id
		WHERE
			album_id = $1 
	`,
		albumID)

	songs := []dao.Song{}
	for rows.Next() {
		var s dao.Song
		err := rows.Scan(
			&s.ID,
			&s.Artist,
			&s.AlbumID,
			&s.Path,
			&s.Title,
			&s.Track,
			&s.Disc,
			&s.GenreID,
			&s.Year,
			&s.Art,
			&s.Created,
			&s.Size,
			&s.Bitrate,
			&s.Duration,
			&s.Token,
			&s.ProviderID,
			&s.Starred,
			&s.AlbumName,
			&s.AlbumArtistID,
			&s.GenreName,
		)
		if err != nil {
		}
		songs = append(songs, s)
	}

	return songs
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

func (db *Default) GetArtist(artistId uint) *dao.Artist {
	var a dao.Artist

	err := db.conn.QueryRow(context.Background(),
		`
		SELECT 
			artists.id,
			artists.name,
			artists.art,
			artists.starred,
			COUNT(albums.id),
			SUM(songs.duration)

		FROM
			artists
		LEFT JOIN 
			albums
		ON
			albums.artist_id = artists.id
		LEFT JOIN 
			songs
		ON
			songs.album_id = albums.id
		WHERE
			artists.id = $1
		GROUP BY
			artists.id
`, artistId,
	).Scan(
		&a.ID,
		&a.Name,
		&a.Art,
		&a.Starred,
		&a.AlbumCount,
		&a.Duration,
	)

	if err != nil {
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

func (db *Default) GetAlbums(listType dao.AlbumList2Type, size uint, offset uint) []dao.Album {

	var order string

	switch listType {
	case dao.Random:
		order = "ORDER BY RANDOM()"
	case dao.Newest:
		order = "ORDER BY albums.created DESC"
	case dao.Frequent:

	case dao.Recent:

	case dao.Starred:

	case dao.AlphabeticalByName:
		order = "ORDER BY albums.name ASC"
	case dao.AlphabeticalByArtist:
		order = "ORDER BY UPPER(artists.name) ASC"
	case dao.ByYear:
		order = "ORDER BY albums.year ASC"
	case dao.ByGenre:
	}

	rows, _ := db.conn.Query(context.Background(),
		fmt.Sprintf(`
		SELECT 
			albums.id,
			albums.artist_id,
			albums.name,
			albums.created,
			albums.art,
			albums.genre_id,
			albums.year,
			albums.disambiguator,
			albums.starred,
			artists.name,
			genres.name,
			COUNT(songs.id),
			SUM(songs.duration)
			FROM
			albums
		LEFT JOIN 
			artists
		ON
			artists.id = albums.artist_id
		LEFT JOIN 
			genres
		ON
			genres.id = albums.genre_id
		LEFT JOIN 
			songs
		ON
			songs.album_id = albums.id
		GROUP BY
			albums.id, artists.id, genres.id
		%v
		LIMIT
			$1
		OFFSET
			$2
	
	`, order), size, offset)

	albums := []dao.Album{}
	for rows.Next() {
		var a dao.Album
		err := rows.Scan(
			&a.ID,
			&a.ArtistID,
			&a.Name,
			&a.Created,
			&a.Art,
			&a.GenreID,
			&a.Year,
			&a.Disambiguator,
			&a.Starred,
			&a.ArtistName,
			&a.GenreName,
			&a.SongCount,
			&a.Duration,
		)
		if err != nil {
		}
		albums = append(albums, a)
	}

	return albums
}

func (db *Default) GetAlbumsByArtist(artistId uint) []dao.Album {

	rows, _ := db.conn.Query(context.Background(),
		`
		SELECT 
			albums.id,
			albums.artist_id,
			albums.name,
			albums.created,
			albums.art,
			albums.genre_id,
			albums.year,
			albums.disambiguator,
			albums.starred,
			artists.name,
			genres.name,
			COUNT(songs.id),
			SUM(songs.duration)
		FROM
			albums
		LEFT JOIN 
			genres
		ON
			genres.id = albums.genre_id
		LEFT JOIN 
			songs
		ON
			songs.album_id = albums.id
		WHERE
			album.artist_id = $1
		GROUP BY
			albums.id, genres.id
	`, artistId)

	albums := []dao.Album{}
	for rows.Next() {
		var a dao.Album
		err := rows.Scan(
			&a.ID,
			&a.ArtistID,
			&a.Name,
			&a.Created,
			&a.Art,
			&a.GenreID,
			&a.Year,
			&a.Disambiguator,
			&a.Starred,
			&a.ArtistName,
			&a.GenreName,
			&a.SongCount,
			&a.Duration,
		)
		if err != nil {
		}
		albums = append(albums, a)
	}

	return albums
}

func (db *Default) GetArtists(includeAlbums bool) []*dao.Artist {
	var artists []*dao.Artist
	query := db.db

	if includeAlbums {
		query = query.Preload("Albums")
	}

	query.Find(&artists)
	return artists
}

func (db *Default) GetPlaylists() []*dao.Playlist {
	var playlists []*dao.Playlist
	db.db.Preload("Entries", func(db *gorm.DB) *gorm.DB {
		return db.Order("playlist_entries.index ASC")
	}).Find(&playlists)
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
	p := db.GetPlaylist(id, false, false)
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

func (db *Default) GetRandomSongs(size uint, from uint, to uint, genre string) []*dao.Song {
	var songs []*dao.Song
	query := db.db.Preload("Genre").
		Preload("Album").
		Preload("Album.Artist")

	if from != 0 {
		query = query.Where("Year >= ?", from)
	}

	if to != 0 {
		query = query.Where("Year <= ?", to)
	}

	if len(genre) > 0 {
		query = query.Joins("JOIN genres ON songs.genre_id = genres.id").Where("genres.name ILIKE ?", genre)
	}

	query.Order(gorm.Expr("random()")).Limit(size).Find(&songs)
	return songs
}

func (db *Default) StarSong(songID uint, star bool) error {
	return db.db.Model(dao.Song{ID: songID}).Updates(dao.Song{Starred: star}).Error
}

func (db *Default) StarAlbum(albumID uint, star bool) error {
	return db.db.Model(dao.Album{ID: albumID}).Updates(dao.Album{Starred: star}).Error
}

func (db *Default) StarArtist(artistID uint, star bool) error {
	return db.db.Model(dao.Artist{ID: artistID}).Updates(dao.Artist{Starred: star}).Error
}
