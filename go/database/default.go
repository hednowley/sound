package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dao"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Default provides access to the default application database.
type Default struct {
	pool *pgxpool.Pool
}

// NewDefault constructs a new default database.
func NewDefault(config *config.Config) (*Default, error) {

	db, err := sql.Open("postgres", config.Db)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+config.MigrationsPath,
		"postgres", driver)
	if err != nil {
		return nil, err
	}

	m.Steps(2)

	pool, err := pgxpool.Connect(context.Background(), config.Db)
	if err != nil {
		return nil, err
	}

	return &Default{pool: pool}, nil
}

func (db *Default) GetConn() (*pgxpool.Conn, error) {
	return db.pool.Acquire(context.Background())
}

// PutArtistByName returns the ID of the artist record with the same name,
// or creates a new one and returns its ID if there is no such artist.
func (db *Default) putArtistByName(conn *pgxpool.Conn, name string) (uint, error) {

	var id uint

	err := conn.QueryRow(context.Background(),
		`
		SELECT 
			id
		FROM
			artists
		WHERE
			lower(name) = lower($1)
		LIMIT
			1
	`, name).Scan(&id)

	if err == nil {
		return id, nil
	}

	err = conn.QueryRow(context.Background(),
		`
		INSERT INTO
			artists (name, starred)
		VALUES
			($1, FALSE)
		RETURNING
			id
	`, name).Scan(&id)

	return id, err
}

// PutAlbumByAttributes returns the ID of the album record with the same name and artist,
// or creates a new one and returns its ID if there is no such album.
func (db *Default) PutAlbumByAttributes(conn *pgxpool.Conn, name string, artist string, disambiguator string) (uint, error) {

	artistID, err := db.putArtistByName(conn, artist)
	if err != nil {
		return 0, err
	}

	var albumID uint

	err = conn.QueryRow(context.Background(),
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

	err = conn.QueryRow(context.Background(),
		`
		INSERT INTO
			albums (artist_id, name, disambiguator, created, starred)
		VALUES
			($1, $2, $3, $4, FALSE)
		RETURNING
			id
	`, artistID, name, disambiguator, time.Now()).Scan(&albumID)

	return albumID, err
}

// PutGenreByName returns the ID of the genre record with the same name,
// or creates a new one and returns its ID if there is no such genre.
func (db *Default) PutGenreByName(conn *pgxpool.Conn, name string) (uint, error) {

	// TODO: Think about whether this needs a transaction
	var genreID uint

	err := conn.QueryRow(context.Background(),
		`
		SELECT
			id
		FROM
			genres
		WHERE
			lower(name) = lower($1)
		LIMIT
			1
	`, name).Scan(&genreID)

	if err == nil {
		return genreID, nil
	}

	if err != pgx.ErrNoRows {
		return 0, err
	}

	err = conn.QueryRow(context.Background(),
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
	conn *pgxpool.Conn,
	artist string,
	albumID uint,
	path string,
	title string,
	track int,
	disc int,
	genreID uint,
	year int,
	artPath *string,
	size int64,
	bitrate int,
	duration int,
	token string,
	providerID string,
) (uint, error) {

	var songID uint

	err := conn.QueryRow(context.Background(),
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
			created,
			starred
		)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, FALSE)
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
		artPath,
		size,
		bitrate,
		duration,
		token,
		providerID,
		time.Now(),
	).Scan(&songID)

	return songID, err
}

func (db *Default) UpdateSong(
	conn *pgxpool.Conn,
	songID uint,
	artist string,
	albumID uint,
	path string,
	title string,
	track int,
	disc int,
	genreID uint,
	year int,
	art *string,
	size int64,
	bitrate int,
	duration int,
	token string,
	providerID string,
	starred bool,
) error {

	_, err := conn.Exec(context.Background(),
		`
		UPDATE 
			songs
		SET
			artist = $1,
			album_id = $2,
			path = $3,
			title = $4,
			track = $5,
			disc = $6,
			genre_id = $7,
			year = $8,
			art = $9,
			size = $10,
			bitrate = $11,
			duration = $12,
			token = $13,
			provider_id = $14,
			created = $15,
			starred = $16
		WHERE
			id = $17
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
		starred,
		songID,
	)

	return err
}

func (db *Default) InsertArt(conn *pgxpool.Conn, path string, hash string) (*dao.Art, error) {

	var art dao.Art

	err := conn.QueryRow(context.Background(),
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

func (db *Default) GetArtFromHash(conn *pgxpool.Conn, hash string) *dao.Art {
	var a dao.Art

	err := conn.QueryRow(context.Background(),
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
func (db *Default) GetSongIdFromToken(conn *pgxpool.Conn, token string, providerID string) (*uint, error) {
	var songId uint

	err := conn.QueryRow(context.Background(),
		`
		SELECT 
			id
		FROM
			songs
		WHERE
			provider_id = $1  AND token = $2
`, providerID, token,
	).Scan(&songId)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return &songId, err
}

func (db *Default) ReplacePlaylistEntries(conn *pgxpool.Conn, playlistID uint, songIDs []uint) error {

	tx, err := conn.Begin(context.Background())
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

	err = tx.Commit(context.Background())

	return err
}

// Getters

func (db *Default) GetPlaylist(conn *pgxpool.Conn, playlistID uint, username string) (*dao.Playlist, error) {
	var p dao.Playlist

	err := conn.QueryRow(
		context.Background(),
		`
		SELECT 
			playlists.id,
			playlists.name,
			playlists.owner,
			playlists.comment,
			timezone('Etc/UTC', playlists.created),
			timezone('Etc/UTC', playlists.changed),
			playlists.public,
			COUNT(playlist_entries.id),
			COALESCE(SUM(songs.duration), 0)
		FROM
			playlists
		LEFT JOIN
			playlist_entries
		ON
			playlist_entries.playlist_id = playlists.id
		LEFT JOIN
			songs
		ON
			songs.id = playlist_entries.song_id
		WHERE
			playlists.id = $1
			AND (playlists.public = TRUE OR playlists.owner = $2)
		GROUP BY
			playlists.id
		`,
		playlistID,
		username,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Owner,
		&p.Comment,
		&p.Created,
		&p.Changed,
		&p.Public,
		&p.EntryCount,
		&p.Duration)

	if err == pgx.ErrNoRows {
		return nil, &dao.ErrNotFound{}
	}

	return &p, err
}

func (db *Default) GetPlaylistSongIds(conn *pgxpool.Conn, playlistID uint) ([]uint, error) {

	rows, err := conn.Query(context.Background(),
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
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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

func (db *Default) GetPlaylistSongs(conn *pgxpool.Conn, playlistID uint, requestor string) ([]dao.Song, error) {

	rows, err := conn.Query(context.Background(),
		`
		SELECT 
			songs.id,
			songs.artist,
			songs.album_id,
			songs.path,
			songs.title,
			songs.track,
			songs.disc,
			songs.year,
			COALESCE(songs.art, ''),
			timezone('Etc/UTC', songs.created),
			songs.size,
			songs.bitrate,
			songs.duration,
			songs.token,
			songs.provider_id,
			songs.starred,
			albums.name,
			albums.artist_id,
			COALESCE(genres.name, '')
		FROM
			playlist_entries
		LEFT JOIN 
			songs
		ON
			songs.id = playlist_entries.song_id
		LEFT JOIN 
			albums
		ON
			albums.id = songs.album_id
		LEFT JOIN 
			genres
		ON
			genres.id = songs.genre_id
		LEFT JOIN 
			playlists
		ON
			playlists.id = playlist_entries.playlist_id
		WHERE
			playlist_entries.playlist_id = $1
			AND (playlists.public = TRUE OR playlists.owner = $2)
		ORDER BY
			playlist_entries.index ASC
	`, playlistID, requestor)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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

func (db *Default) GetAlbum(conn *pgxpool.Conn, albumID uint) (*dao.Album, error) {

	var a dao.Album

	err := conn.QueryRow(
		context.Background(),
		`
		SELECT 
			albums.id,
			albums.artist_id,
			albums.name,
			timezone('Etc/UTC', albums.created),
			albums.disambiguator,
			albums.starred,
			artists.name,
			COUNT(songs.id),
			COALESCE(SUM(songs.duration), 0),
			array_remove(array_agg(DISTINCT songs.art), NULL),
			array_remove(array_agg(DISTINCT genres.name), NULL),
			array_remove(array_remove(array_agg(DISTINCT songs.year), 0), NULL)
		FROM
			albums
		LEFT JOIN 
			artists
		ON
			artists.id = albums.artist_id
		LEFT JOIN 
			songs
		ON
			songs.album_id = albums.id
		LEFT JOIN 
			genres
		ON
			genres.id = songs.genre_id
		WHERE
			albums.id = $1
		GROUP BY
			albums.id, artists.id`,
		albumID,
	).Scan(
		&a.ID,
		&a.ArtistID,
		&a.Name,
		&a.Created,
		&a.Disambiguator,
		&a.Starred,
		&a.ArtistName,
		&a.SongCount,
		&a.Duration,
		&a.Arts,
		&a.Genres,
		&a.Years,
	)

	if err == pgx.ErrNoRows {
		return nil, &dao.ErrNotFound{}
	}

	return &a, err
}

func (db *Default) GetAlbumSongs(conn *pgxpool.Conn, albumID uint) ([]dao.Song, error) {

	rows, err := conn.Query(context.Background(),
		`
		SELECT 
			songs.id,
			songs.artist,
			songs.album_id,
			songs.path,
			songs.title,
			songs.track,
			songs.disc,
			songs.year,
			COALESCE(songs.art, ''),
			timezone('Etc/UTC', songs.created),
			songs.size,
			songs.bitrate,
			songs.duration,
			songs.token,
			songs.provider_id,
			songs.starred,
			albums.name,
			albums.artist_id,
			COALESCE(genres.name, '')
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
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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
			return nil, err
		}
		songs = append(songs, s)
	}

	return songs, nil
}

func (db *Default) GetArt(conn *pgxpool.Conn, artId uint) *dao.Art {
	var a dao.Art

	err := conn.QueryRow(context.Background(),
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

func (db *Default) GetSongsByGenre(conn *pgxpool.Conn, genreName string, offset uint, limit uint) ([]dao.Song, error) {

	rows, err := conn.Query(context.Background(),
		`
		SELECT 
			songs.id,
			songs.artist,
			songs.album_id,
			songs.path,
			songs.title,
			songs.track,
			songs.disc,
			songs.year,
			COALESCE(songs.art, ''),
			timezone('Etc/UTC', songs.created),
			songs.size,
			songs.bitrate,
			songs.duration,
			songs.token,
			songs.provider_id,
			songs.starred,
			albums.name,
			albums.artist_id,
			COALESCE(genres.name, '')
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
			lower(genres.name) = lower($1)
		ORDER BY
			songs.title ASC
		OFFSET
			$2
		LIMIT
			$3
	`,
		genreName, offset, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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
			return nil, err
		}
		songs = append(songs, s)
	}

	return songs, nil
}

func (db *Default) GetArtist(conn *pgxpool.Conn, artistId uint) (*dao.Artist, error) {
	var a dao.Artist

	err := conn.QueryRow(context.Background(),
		`
		SELECT 
			artists.id,
			artists.name,
			artists.starred,
			COUNT(DISTINCT albums.id),
			COALESCE(SUM(songs.duration), 0),
			array_remove(array_agg(DISTINCT songs.art), NULL)
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
		&a.Starred,
		&a.AlbumCount,
		&a.Duration,
		&a.Arts,
	)

	if err == pgx.ErrNoRows {
		return nil, &dao.ErrNotFound{}
	}

	return &a, err
}

func (db *Default) GetSong(conn *pgxpool.Conn, songId uint) (*dao.Song, error) {
	var s dao.Song

	err := conn.QueryRow(context.Background(),
		`
		SELECT 
			songs.id,
			songs.artist,
			songs.album_id,
			songs.path,
			songs.title,
			songs.track,
			songs.disc,
			songs.year,
			COALESCE(songs.art, ''),
			timezone('Etc/UTC', songs.created),
			songs.size,
			songs.bitrate,
			songs.duration,
			songs.token,
			songs.provider_id,
			songs.starred,
			albums.name,
			albums.artist_id,
			COALESCE(genres.name, '')
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

	if err == pgx.ErrNoRows {
		return nil, &dao.ErrNotFound{}
	}

	return &s, err
}

func (db *Default) GetSongPath(conn *pgxpool.Conn, songId uint) *string {
	var path string

	err := conn.QueryRow(context.Background(),
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

func (db *Default) GetGenres(conn *pgxpool.Conn) ([]dao.Genre, error) {

	rows, err := conn.Query(context.Background(),
		`
		SELECT 
			genres.id,
			genres.name,
			COUNT(DISTINCT songs.id),
			COUNT(DISTINCT albums.id)
		FROM
			genres
		LEFT JOIN
			songs
		ON
			songs.genre_id = genres.id
		LEFT JOIN
			albums
		ON
			albums.id = songs.album_id
		GROUP BY
			genres.id
		ORDER BY
			genres.name ASC`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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
			return nil, err
		}
		genres = append(genres, g)
	}

	return genres, nil
}

func (db *Default) GetAlbums(conn *pgxpool.Conn, listType dao.AlbumList2Type, limit uint, offset uint) ([]dao.Album, error) {

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

	query := fmt.Sprintf(`
		SELECT 
			albums.id,
			albums.artist_id,
			albums.name,
			timezone('Etc/UTC', albums.created),
			albums.disambiguator,
			albums.starred,
			artists.name,
			COUNT(songs.id),
			COALESCE(SUM(songs.duration), 0),
			array_remove(array_agg(DISTINCT songs.art), NULL),
			array_remove(array_agg(DISTINCT genres.name), NULL),
			array_remove(array_remove(array_agg(DISTINCT songs.year), 0), NULL)
		FROM
			albums
		LEFT JOIN 
			artists
		ON
			artists.id = albums.artist_id
		LEFT JOIN 
			songs
		ON
			songs.album_id = albums.id
		LEFT JOIN 
			genres
		ON
			genres.id = songs.genre_id
		GROUP BY
			albums.id, artists.id
		%v
		LIMIT
			$1
		OFFSET
			$2

	`, order)

	rows, err := conn.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	albums := []dao.Album{}
	for rows.Next() {
		var a dao.Album
		rowErr := rows.Scan(
			&a.ID,
			&a.ArtistID,
			&a.Name,
			&a.Created,
			&a.Disambiguator,
			&a.Starred,
			&a.ArtistName,
			&a.SongCount,
			&a.Duration,
			&a.Arts,
			&a.Genres,
			&a.Years,
		)
		if rowErr != nil {
			return nil, rowErr
		}
		albums = append(albums, a)
	}

	return albums, nil
}

func (db *Default) GetAlbumsByArtist(conn *pgxpool.Conn, artistId uint) ([]dao.Album, error) {

	rows, err := conn.Query(context.Background(),
		`
		SELECT 
			albums.id,
			albums.artist_id,
			albums.name,
			timezone('Etc/UTC', albums.created),
			albums.disambiguator,
			albums.starred,
			artists.name,
			COUNT(songs.id),
			COALESCE(SUM(songs.duration), 0),
			array_remove(array_agg(DISTINCT songs.art), NULL),
			array_remove(array_agg(DISTINCT genres.name), NULL),
			array_remove(array_remove(array_agg(DISTINCT songs.year), 0), NULL)
		FROM
			albums
		LEFT JOIN 
			artists
		ON
			artists.id = albums.artist_id
		LEFT JOIN 
			songs
		ON
			songs.album_id = albums.id
		LEFT JOIN 
			genres
		ON
			genres.id = songs.genre_id
		WHERE
			albums.artist_id = $1
		GROUP BY
			albums.id, artists.id
	`, artistId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	albums := []dao.Album{}
	for rows.Next() {
		var a dao.Album
		err := rows.Scan(
			&a.ID,
			&a.ArtistID,
			&a.Name,
			&a.Created,
			&a.Disambiguator,
			&a.Starred,
			&a.ArtistName,
			&a.SongCount,
			&a.Duration,
			&a.Arts,
			&a.Genres,
			&a.Years,
		)
		if err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}

	return albums, nil
}

func (db *Default) GetArtists(conn *pgxpool.Conn) ([]dao.Artist, error) {

	rows, err := conn.Query(context.Background(),
		`
		SELECT 
			artists.id,
			artists.name,
			artists.starred,
			COUNT(DISTINCT albums.id),
			COALESCE(SUM(songs.duration), 0),
			array_remove(array_agg(DISTINCT songs.art), NULL)
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
			artists.id
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	artists := []dao.Artist{}
	for rows.Next() {
		var a dao.Artist
		rowErr := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Starred,
			&a.AlbumCount,
			&a.Duration,
			&a.Arts,
		)
		if rowErr != nil {
			return nil, rowErr
		}
		artists = append(artists, a)
	}

	return artists, nil
}

func (db *Default) GetPlaylists(conn *pgxpool.Conn, username string) ([]dao.Playlist, error) {

	rows, err := conn.Query(context.Background(),
		`
		SELECT 
			playlists.id,
			playlists.name,
			playlists.comment,
			timezone('Etc/UTC', playlists.created),
			timezone('Etc/UTC', playlists.changed),
			playlists.public,
			playlists.owner,
			COUNT(playlist_entries.id)
		FROM
			playlists
		LEFT JOIN
			playlist_entries
		ON
			playlist_entries.playlist_id = playlists.id
		WHERE
			playlists.public = TRUE OR playlists.owner = $1
		GROUP BY
			playlists.id
		`, username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	playlists := []dao.Playlist{}
	for rows.Next() {
		var p dao.Playlist
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Comment,
			&p.Created,
			&p.Changed,
			&p.Public,
			&p.Owner,
			&p.EntryCount,
		)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}

	return playlists, nil
}

// Putters

func (db *Default) InsertPlaylist(conn *pgxpool.Conn, name string, comment string, owner string, public bool) (uint, error) {

	var playlistID uint

	err := conn.QueryRow(context.Background(),
		`
		INSERT INTO
			playlists (name, comment, owner, created, changed, public)
		VALUES
			($1, $2, $3, $4, $4, $5)
		RETURNING
			id
	`, name, comment, owner, time.Now(), public).Scan(&playlistID)

	return playlistID, err
}

func (db *Default) UpdatePlaylist(conn *pgxpool.Conn, playlistID uint, name string, public bool, comment string) (*dao.Playlist, error) {

	var playlist dao.Playlist

	err := conn.QueryRow(context.Background(),
		`
		UPDATE
			playlists
		SET
			name = $2,
			public = $3,
			comment = $4,
			changed = $5
		WHERE
			id = $1
		RETURNING
			id,
			name,
			comment,
			timezone('Etc/UTC', created),
			timezone('Etc/UTC', changed)
	`, playlistID, name, public, comment, time.Now()).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.Comment,
		&playlist.Created,
		&playlist.Changed,
	)

	return &playlist, err
}

// Deleters

func (db *Default) DeletePlaylist(conn *pgxpool.Conn, playlistID uint, owner string) error {

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		`
			DELETE FROM
				playlist_entries
			WHERE id IN (
				SELECT
					playlist_entries.id
				FROM
					playlist_entries
				LEFT JOIN
					playlists
				ON
					playlists.id = playlist_entries.playlist_id
				WHERE
					playlist_entries.playlist_id = $1
					AND playlists.owner = $2
			)
		`,
		playlistID,
		owner,
	)
	if err != nil {
		return err
	}

	ct, err := tx.Exec(
		context.Background(),
		`
			DELETE FROM
				playlists
			WHERE
				id = $1 AND owner = $2
		`,
		playlistID,
		owner,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return &dao.ErrNotFound{}
	}

	err = tx.Commit(context.Background())

	return err

}

func (db *Default) DeleteMissing(conn *pgxpool.Conn, tokens []string, providerID string) error {

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	values := []interface{}{providerID}
	params := []string{}
	for _, token := range tokens {
		values = append(values, token)
		params = append(params, fmt.Sprintf("$%d", len(values)))
	}

	query := fmt.Sprintf(`
		DELETE FROM
			playlist_entries
		USING
			playlist_entries AS pe
		INNER JOIN
			songs
		ON
			songs.id = pe.song_id
		WHERE
			songs.provider_id = $1 AND
			songs.token NOT IN (%v)
	`, strings.Join(params, ","))

	_, err = tx.Exec(
		context.Background(),
		query,
		values...,
	)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
		DELETE FROM
			songs
		WHERE
			provider_id = $1 AND
			token NOT IN (%v)
	`, strings.Join(params, ","))

	_, err = tx.Exec(
		context.Background(),
		query,
		values...,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		`
		DELETE FROM
			albums
		WHERE id in (
			SELECT
				albums.id
			FROM
				albums
			LEFT JOIN
				songs
			ON
				songs.album_id = albums.id
			GROUP BY
				albums.id
			HAVING
				COUNT(songs.id) = 0
		)
		`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		`
		DELETE FROM
			artists
		WHERE id IN (
			SELECT
				artists.id
			FROM
				artists
			LEFT JOIN
				albums
			ON
				albums.artist_id = artists.id
			GROUP BY
				artists.id
			HAVING
				COUNT(albums.id) = 0
		)
		`)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (db *Default) SearchAlbums(conn *pgxpool.Conn, query string, count uint, offset uint) ([]dao.Album, error) {

	rows, err := conn.Query(context.Background(),
		`
		SELECT 
			albums.id,
			albums.artist_id,
			albums.name,
			timezone('Etc/UTC', albums.created),
			albums.disambiguator,
			albums.starred,
			artists.name,
			COUNT(songs.id),
			COALESCE(SUM(songs.duration), 0),
			array_remove(array_agg(DISTINCT songs.art), NULL),
			array_remove(array_agg(DISTINCT genres.name), NULL),
			array_remove(array_remove(array_agg(DISTINCT songs.year), 0), NULL)
		FROM
			albums
		LEFT JOIN 
			artists
		ON
			artists.id = albums.artist_id
		LEFT JOIN 
			songs
		ON
			songs.album_id = albums.id
		LEFT JOIN 
			genres
		ON
			genres.id = songs.genre_id
		WHERE
			albums.name ILIKE $1
		GROUP BY
			albums.id, artists.id
		ORDER BY
			albums.name ASC
		LIMIT
			$2
		OFFSET
			$3
	
	`, "%"+query+"%", count, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	albums := []dao.Album{}
	for rows.Next() {
		var a dao.Album
		err := rows.Scan(
			&a.ID,
			&a.ArtistID,
			&a.Name,
			&a.Created,
			&a.Disambiguator,
			&a.Starred,
			&a.ArtistName,
			&a.SongCount,
			&a.Duration,
			&a.Arts,
			&a.Genres,
			&a.Years,
		)
		if err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}

	return albums, nil
}

func (db *Default) SearchArtists(conn *pgxpool.Conn, query string, limit uint, offset uint) ([]dao.Artist, error) {

	rows, err := conn.Query(context.Background(),
		`
		SELECT 
			artists.id,
			artists.name,
			artists.starred,
			COUNT(DISTINCT albums.id),
			COALESCE(SUM(songs.duration), 0),
			array_remove(array_agg(DISTINCT songs.art), NULL)
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
			artists.name ILIKE $1
		GROUP BY
			artists.id
		ORDER BY
			artists.name ASC
		LIMIT
			$2
		OFFSET
			$3
		
	`, "%"+query+"%", limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	artists := []dao.Artist{}
	for rows.Next() {
		var a dao.Artist
		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Starred,
			&a.AlbumCount,
			&a.Duration,
			&a.Arts,
		)
		if err != nil {
			return nil, err
		}
		artists = append(artists, a)
	}

	return artists, nil
}

func (db *Default) SearchSongs(conn *pgxpool.Conn, query string, count uint, offset uint) ([]dao.Song, error) {
	rows, err := conn.Query(context.Background(),
		`
		SELECT 
			songs.id,
			songs.artist,
			songs.album_id,
			songs.path,
			songs.title,
			songs.track,
			songs.disc,
			songs.year,
			COALESCE(songs.art, ''),
			timezone('Etc/UTC', songs.created),
			songs.size,
			songs.bitrate,
			songs.duration,
			songs.token,
			songs.provider_id,
			songs.starred,
			albums.name,
			albums.artist_id,
			COALESCE(genres.name, '')
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
			songs.title ILIKE $1
		ORDER BY
			songs.title ASC
		LIMIT
			$2
		OFFSET
			$3
	`,
		"%"+query+"%", count, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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
			return nil, err
		}
		songs = append(songs, s)
	}

	return songs, nil
}

func (db *Default) GetRandomSongs(conn *pgxpool.Conn, size uint, from uint, to uint, genre string) ([]dao.Song, error) {

	values := []interface{}{size}
	wheres := []string{}

	if from != 0 {
		values = append(values, from)
		wheres = append(wheres, fmt.Sprintf("Year >= $%v", len(values)))
	}

	if to != 0 {
		values = append(values, to)
		wheres = append(wheres, fmt.Sprintf("Year <= $%v", len(values)))
	}

	if len(genre) > 0 {
		values = append(values, genre)
		wheres = append(wheres, fmt.Sprintf("lower(genres.name) = lower($%v)", len(values)))
	}

	var where string
	if len(wheres) > 0 {
		where = "WHERE " + strings.Join(wheres, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT 
			songs.id,
			songs.artist,
			songs.album_id,
			songs.path,
			songs.title,
			songs.track,
			songs.disc,
			songs.year,
			COALESCE(songs.art, ''),
			timezone('Etc/UTC', songs.created),
			songs.size,
			songs.bitrate,
			songs.duration,
			songs.token,
			songs.provider_id,
			songs.starred,
			albums.name,
			albums.artist_id,
			COALESCE(genres.name, '')
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
		%v
		ORDER BY
			RANDOM()
		LIMIT
			$1`,
		where)

	rows, err := conn.Query(
		context.Background(),
		query,
		values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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
			return nil, err
		}
		songs = append(songs, s)
	}

	return songs, nil
}
