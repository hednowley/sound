package database

import (
	"context"
	"fmt"
	"strings"
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

// PutArtistByName returns the ID of the artist record with the same name,
// or creates a new one and returns its ID if there is no such artist.
func (db *Default) putArtistByName(name string) (uint, error) {

	var id uint

	err := db.conn.QueryRow(context.Background(),
		`
		SELECT 
			id
		FROM
			artists
		WHERE
			name ILIKE $1
		LIMIT
			1
	`, name).Scan(&id)

	if err == nil {
		return id, nil
	}

	err = db.conn.QueryRow(context.Background(),
		`
		INSERT INTO
			artists (name)
		VALUES
			($1)
		RETURNING
			id
	`, name).Scan(&id)

	return id, err
}

// PutAlbumByAttributes returns the ID of the album record with the same name and artist,
// or creates a new one and returns its ID if there is no such album.
func (db *Default) PutAlbumByAttributes(name string, artist string, disambiguator string) (uint, error) {

	artistID, err := db.putArtistByName(artist)
	if err != nil {
		return 0, err
	}

	var albumID uint

	err = db.conn.QueryRow(context.Background(),
		`
		SELECT
			id
		FROM
			albums
		WHERE
			name ILIKE $1 AND
			artist_id = $2 AND
			disambiguator = $3
		LIMIT
			1
	`, name, artistID, disambiguator).Scan(&albumID)

	if err == nil {
		return albumID, nil
	}

	err = db.conn.QueryRow(context.Background(),
		`
		INSERT INTO
			albums (artist_id, name, disambiguator, created)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			id
	`, artistID, name, disambiguator, time.Now()).Scan(&albumID)

	return albumID, err
}

// PutGenreByName returns the ID of the genre record with the same name,
// or creates a new one and returns its ID if there is no such genre.
func (db *Default) PutGenreByName(name string) (uint, error) {

	var genreID uint

	err := db.conn.QueryRow(context.Background(),
		`
		SELECT
			id
		FROM
			genres
		WHERE
			name ILIKE $1
		LIMIT
			1
	`, name).Scan(&genreID)

	if err == nil {
		return genreID, nil
	}

	err = db.conn.QueryRow(context.Background(),
		`
		INSERT INTO
			genres (name)
		VALUES
			($1)
		RETURNING
			id
	`, name).Scan(&genreID)

	return genreID, err
}

func (db *Default) InsertSong(
	artist string,
	albumID uint,
	path string,
	title string,
	track int,
	disc int,
	genreID uint,
	year int,
	art string,
	size int64,
	bitrate int,
	duration int,
	token string,
	providerID string,
) (uint, error) {

	var songID uint

	err := db.conn.QueryRow(context.Background(),
		`
		INSERT INTO
			songs
		(
			artist,
			album_id,
			path,
			title,
			track,
			disc,
			genre_id,
			year,
			art,
			size,
			bitrate,
			duration,
			token,
			provider_id,
			created
		)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING
			id
	`,
		artist,
		albumID,
		path,
		title,
		track,
		disc,
		genreID,
		year,
		art,
		size,
		bitrate,
		duration,
		token,
		providerID,
		time.Now(),
	).Scan(&songID)

	return songID, err
}

func (db *Default) InsertArt(path string, hash string) (*dao.Art, error) {

	var art dao.Art

	err := db.conn.QueryRow(context.Background(),
		`
		INSERT INTO
			arts (path, hash)
		VALUES
			($1, $2)
		RETURNING
			id, path, hash
	`, path, hash).Scan(&art.ID, &art.Path, &art.Hash)

	return &art, err
}

func (db *Default) GetArtFromHash(hash string) *dao.Art {
	var a dao.Art

	err := db.conn.QueryRow(context.Background(),
		`
		SELECT 
			id,
			hash,
			path
		FROM
			arts
		WHERE
			hash = $1
	`, hash,
	).Scan(
		&a.ID,
		&a.Hash,
		&a.Path,
	)

	if err != nil {
		return nil
	}

	return &a
}

// GetSongIdFromToken returns a pointer to the song with the given path and provider,
// or nil if one doesn't exist. Joined entities are not loaded.
func (db *Default) GetSongIdFromToken(token string, providerID string) *uint {
	var songId uint

	err := db.conn.QueryRow(context.Background(),
		`
		SELECT 
			id
		FROM
			songs
		WHERE
			provider_id = $1  AND token = $2
`, providerID, token,
	).Scan(&songId)

	if err != nil {
		return nil
	}

	return &songId
}

func (db *Default) ReplacePlaylistEntries(playlistID uint, songIDs []uint) error {

	tx, err := db.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `
		DELETE FROM
			playlist_entries
		WHERE
			playlist_id = $1
	`, playlistID)

	if err != nil {
		return err
	}

	rows := make([][]interface{}, len(songIDs))
	for index, songID := range songIDs {
		rows[index] = []interface{}{playlistID, songID, index}
	}

	_, err = tx.CopyFrom(
		context.Background(),
		[]string{"playlist_entries"},
		[]string{"playlist_id", "song_id", "index"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return err
	}

	tx.Commit(context.Background())

	return nil
}

// Getters

func (db *Default) GetPlaylist(playlistID uint) (*dao.Playlist, error) {
	var p dao.Playlist

	err := db.conn.QueryRow(context.Background(),
		`
		SELECT 
			playlists.id,
			playlists.name,
			playlists.comment,
			playlists.created,
			playlists.changed,
			COUNT(playlist_entries.id)
		FROM
			playlists
		LEFT JOIN
			playlist_entries
		ON
			playlist_entries.playlist_id = playlists.id
		GROUP BY
			playlists.id
		WHERE
			id = $1
	`, playlistID).Scan(
		&p.ID,
		&p.Name,
		&p.Comment,
		&p.Created,
		&p.Changed,
		&p.EntryCount)

	return &p, err
}

func (db *Default) GetPlaylistSongIds(playlistID uint) ([]uint, error) {

	rows, _ := db.conn.Query(context.Background(),
		`
		SELECT 
			song_id
		FROM
			playlists_entries
		WHERE
			playlist_id = $1
		ORDER BY
			playlist_entries.index ASC
	`, playlistID)

	songIDs := []uint{}
	for rows.Next() {
		var songID uint
		err := rows.Scan(&songID)
		if err != nil {
		}
		songIDs = append(songIDs, songID)
	}

	return songIDs, nil
}

func (db *Default) GetPlaylistSongs(playlistID uint) ([]dao.Song, error) {

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
			playlist_entries
		LEFT JOIN 
			songs
		ON
			song.id = playlist_entries.song_id
		LEFT JOIN 
			albums
		ON
			albums.id = songs.album_id
		LEFT JOIN 
			genres
		ON
			genres.id = songs.genre_id
		WHERE
			playlist_entries.playlist_id = $1 
		ORDER BY
			playlist_entries.index ASC
	`, playlistID)

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

	return songs, nil
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
		ORDER BY
			songs.track ASC
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

func (db *Default) GetArt(artId uint) *dao.Art {
	var a dao.Art

	err := db.conn.QueryRow(context.Background(),
		`
		SELECT 
			id,
			hash,
			path
		FROM
			arts
		WHERE
			id = $1
	`, artId,
	).Scan(
		&a.ID,
		&a.Hash,
		&a.Path,
	)

	if err != nil {
		return nil
	}

	return &a
}

func (db *Default) GetSongsByGenre(genreName string, offset uint, limit uint) []dao.Song {

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
			genres
		ON
			genres.id = songs.genre_id
		LEFT JOIN
			albums
		ON
			albums.id = songs.album_id
		WHERE
			genres.name = $1
		ORDER BY
			songs.title ASC
		OFFSET
			$2
		LIMIT
			$3
	`,
		genreName, offset, limit)

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

func (db *Default) GetSong(songId uint) *dao.Song {
	var s dao.Song

	err := db.conn.QueryRow(context.Background(),
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
			songs.id = $1 
`, songId,
	).Scan(
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
		return nil
	}

	return &s
}

func (db *Default) GetSongPath(songId uint) *string {
	var path string

	err := db.conn.QueryRow(context.Background(),
		`
		SELECT 
			path
		FROM
			songs
		WHERE
			id = $1 
`, songId).Scan(&path)

	if err != nil {
		return nil
	}

	return &path
}

// Collection getters

func (db *Default) GetGenres() []dao.Genre {

	rows, _ := db.conn.Query(context.Background(),
		`
		SELECT 
			genres.id,
			genres.name,
			COUNT(songs.id),
			COUNT(albums.id)
		FROM
			genres
		LEFT JOIN
			songs
		ON
			songs.genre_id = genres.id
		LEFT JOIN
			albums
		ON
			albums.genre_id = genres.id
		GROUP BY
			genres.id, songs.id, albums.id
		ORDER BY
			genres.name ASC`)

	genres := []dao.Genre{}
	for rows.Next() {
		var g dao.Genre
		err := rows.Scan(
			&g.ID,
			&g.Name,
			&g.SongCount,
			&g.AlbumCount,
		)
		if err != nil {
		}
		genres = append(genres, g)
	}

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
		order = ""
	case dao.Recent:
		order = ""
	case dao.Starred:
		order = ""
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
		ORDER BY
			year ASC
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

func (db *Default) GetArtists() []dao.Artist {

	rows, _ := db.conn.Query(context.Background(),
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
		GROUP BY
			albums.id, artists.id, genres.id
	`)

	artists := []dao.Artist{}
	for rows.Next() {
		var a dao.Artist
		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Art,
			&a.Starred,
			&a.AlbumCount,
			&a.Duration,
		)
		if err != nil {
		}
		artists = append(artists, a)
	}

	return artists
}

func (db *Default) GetPlaylists() []dao.Playlist {

	rows, _ := db.conn.Query(context.Background(),
		`
		SELECT 
			playlists.id,
			playlists.name,
			playlists.comment,
			playlists.created,
			playlists.changed,
			COUNT(playlist_entries.id)
		FROM
			playlists
		LEFT JOIN
			playlist_entries
		ON
			playlist_entries.playlist_id = playlists.id
		GROUP BY
			playlists.id
`)

	playlists := []dao.Playlist{}
	for rows.Next() {
		var p dao.Playlist
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Comment,
			&p.Created,
			&p.Changed,
			&p.EntryCount,
		)
		if err != nil {
		}
		playlists = append(playlists, p)
	}

	return playlists
}

// Putters

func (db *Default) InsertPlaylist(name string, comment string) (uint, error) {

	var playlistID uint

	err := db.conn.QueryRow(context.Background(),
		`
		INSERT INTO
			playlists (name, comment, created, changed)
		VALUES
			($1, $2, $3, $3)
		RETURNING
			id
	`, name, comment, time.Now()).Scan(&playlistID)

	return playlistID, err
}

func (db *Default) UpdatePlaylist(playlistID uint, name string, comment string) (*dao.Playlist, error) {

	var playlist dao.Playlist

	err := db.conn.QueryRow(context.Background(),
		`
		UPDATE
			playlists
		SET
			name = $2,
			comment = $3,
		WHERE
			id = $1
		RETURNING
			id,
			name,
			comment,
			created,
			changed
	`, playlistID, name, comment, time.Now()).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.Comment,
		&playlist.Created,
		&playlist.Changed,
	)

	return &playlist, err
}

// Deleters

func (db *Default) DeletePlaylist(id uint) error {
	p, err := db.GetPlaylist(id)
	if err != nil {
		return err
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

func (db *Default) SearchAlbums(query string, count uint, offset uint) []dao.Album {

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
			albums.name ILIKE %$1%
		GROUP BY
			albums.id, artists.id, genres.id
		ORDER BY
			albums.name ASC
		LIMIT
			$2
		OFFSET
			$3
	
	`, query, count, offset)

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

func (db *Default) SearchArtists(query string, count uint, offset uint) []dao.Artist {

	rows, _ := db.conn.Query(context.Background(),
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
		GROUP BY
			albums.id, artists.id, genres.id
		WHERE
			artists.name ILIKE %$1%
		ORDER BY
			artists.name ASC
		LIMIT
			$2
		OFFSET
			$3
		
	`, query, count, offset)

	artists := []dao.Artist{}
	for rows.Next() {
		var a dao.Artist
		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Art,
			&a.Starred,
			&a.AlbumCount,
			&a.Duration,
		)
		if err != nil {
		}
		artists = append(artists, a)
	}

	return artists
}

func (db *Default) SearchSongs(query string, count uint, offset uint) []dao.Song {
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
			genres
		ON
			genres.id = songs.genre_id
		LEFT JOIN
			albums
		ON
			albums.id = songs.album_id
		WHERE
			songs.title ILIKE %$1%
		ORDER BY
			songs.title ASC
		LIMIT
			$2
		OFFSET
			$3
	`,
		query, count, offset)

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

func (db *Default) GetRandomSongs(size uint, from uint, to uint, genre string) []dao.Song {

	values := []interface{}{size}
	wheres := []string{}
	next := len(values) + 1

	if from != 0 {
		wheres = append(wheres, fmt.Sprintf("Year >= $%v", next))
		values = append(values, from)
		next++
	}

	if to != 0 {
		wheres = append(wheres, fmt.Sprintf("Year <= $%v", next))
		values = append(values, to)
		next++
	}

	if len(genre) > 0 {
		wheres = append(wheres, fmt.Sprintf("genres.name ILIKE $%v", next))
		values = append(values, genre)
	}

	where := strings.Join(wheres, " AND ")

	rows, _ := db.conn.Query(context.Background(),
		fmt.Sprintf(`
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
			genres
		ON
			genres.id = songs.genre_id
		LEFT JOIN
			albums
		ON
			albums.id = songs.album_id
		WHERE
			%v
		ORDER BY
			RANDOM()
		LIMIT
			$1
	`, where), values...)

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

func (db *Default) StarSong(songID uint, star bool) error {
	return db.db.Model(dao.Song{ID: songID}).Updates(dao.Song{Starred: star}).Error
}

func (db *Default) StarAlbum(albumID uint, star bool) error {
	return db.db.Model(dao.Album{ID: albumID}).Updates(dao.Album{Starred: star}).Error
}

func (db *Default) StarArtist(artistID uint, star bool) error {
	return db.db.Model(dao.Artist{ID: artistID}).Updates(dao.Artist{Starred: star}).Error
}
