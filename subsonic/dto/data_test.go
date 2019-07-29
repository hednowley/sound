package dto

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hednowley/sound/dao"
)

func GenerateGenre(id int) *dao.Genre {
	return &dao.Genre{
		ID:   uint(id),
		Name: fmt.Sprintf("genre%v", id),
	}
}

func GenerateGenres(count int) []*dao.Genre {
	g := make([]*dao.Genre, count)
	for i := 1; i < count+1; i++ {
		g[i-1] = GenerateGenre(i)
	}
	return g
}

func GenerateArt(id int) *dao.Art {
	return &dao.Art{
		ID:   uint(id),
		Hash: fmt.Sprintf("cde22cace9e39eee25f1ece6%v", id),
		Path: fmt.Sprintf("%v.jpg", id),
	}
}

func GenerateArts(count int) []*dao.Art {
	arts := make([]*dao.Art, count)
	for i := 1; i < count+1; i++ {
		arts[i-1] = GenerateArt(i)
	}
	return arts
}

func GenerateSong(i int, genre *dao.Genre, album *dao.Album, art *dao.Art) *dao.Song {
	t := time.Date(2000+i, time.August, 15, 0, 0, 0, 0, time.UTC)

	var artPath string
	if art != nil {
		artPath = art.Path
	}

	var genreID uint
	if genre != nil {
		genreID = genre.ID
	}

	s := dao.Song{
		ID:        uint(i),
		Title:     fmt.Sprintf("title%v", i),
		Artist:    album.Artist.Name,
		Track:     i,
		Disc:      i,
		Extension: fmt.Sprintf("mp%v", i),
		GenreID:   genreID,
		Genre:     genre,
		Album:     album,
		AlbumID:   album.ID,
		Art:       artPath,
		Path:      fmt.Sprintf("D:\\music\\%v.mp3", i),
		Size:      int64(i * 100000),
		Year:      1900 + i,
		Created:   &t,
	}

	album.Songs = append(album.Songs, &s)
	if genre != nil {
		genre.Songs = append(genre.Songs, &s)
	}
	return &s
}

func GenerateSongs(count int, genre *dao.Genre, album *dao.Album, art *dao.Art) []*dao.Song {
	songs := make([]*dao.Song, count)
	for i := 1; i < count+1; i++ {
		songs[i-1] = GenerateSong(i, genre, album, art)
	}
	return songs
}

func GenerateAlbum(i int, genre *dao.Genre, artist *dao.Artist, art *dao.Art) *dao.Album {
	t := time.Date(2000+i, time.August, 15, 0, 0, 0, 0, time.UTC)

	var artPath string
	if art != nil {
		artPath = art.Path
	}

	var genreID uint
	if genre != nil {
		genreID = genre.ID
	}

	a := dao.Album{
		ID:       uint(i),
		Name:     fmt.Sprintf("album%v", i),
		Year:     1900 + i,
		GenreID:  genreID,
		Genre:    genre,
		Artist:   artist,
		ArtistID: artist.ID,
		Art:      artPath,
		Created:  &t,
		Duration: i,
	}

	artist.Albums = append(artist.Albums, &a)
	if genre != nil {
		genre.Albums = append(genre.Albums, &a)
	}
	return &a
}

func GenerateAlbums(count int, genre *dao.Genre, artist *dao.Artist, art *dao.Art) []*dao.Album {
	albums := make([]*dao.Album, count)
	for i := 1; i < count+1; i++ {
		albums[i-1] = GenerateAlbum(i, genre, artist, art)
	}
	return albums
}

func GenerateArtist(i int, art *dao.Art) *dao.Artist {

	var artPath string
	if art != nil {
		artPath = art.Path
	}

	a := dao.Artist{
		ID:       uint(i),
		Name:     fmt.Sprintf("artist%v", i),
		Art:      artPath,
		Duration: i,
	}

	return &a
}

func GenerateArtists(count int, art *dao.Art) []*dao.Artist {
	artists := make([]*dao.Artist, count)
	for i := 1; i < count+1; i++ {
		artists[i-1] = GenerateArtist(i, art)
	}
	return artists
}

func GeneratePlaylist(i int, public bool, songs ...*dao.Song) *dao.Playlist {

	t1 := time.Date(2000+i, time.August, 10, 0, 30, 0, 9, time.UTC)
	t2 := time.Date(2001+i, time.January, 15, 5, 1, 0, 1, time.UTC)

	p := dao.Playlist{
		ID:      uint(i),
		Name:    fmt.Sprintf("playlist%v", i),
		Comment: fmt.Sprintf("comment%v", i),
		Public:  public,
		Changed: &t1,
		Created: &t2,
	}

	entries := make([]*dao.PlaylistEntry, len(songs))

	for i, s := range songs {
		entries[i] = &dao.PlaylistEntry{
			ID:         uint(i + 1),
			Index:      i,
			PlaylistID: p.ID,
			SongID:     s.ID,
			Song:       s,
		}
	}

	p.Entries = entries
	return &p
}

func GeneratePlaylists(public bool, songCollections ...[]*dao.Song) []*dao.Playlist {

	count := len(songCollections)
	playlists := make([]*dao.Playlist, count)
	for i := 1; i < count+1; i++ {
		playlists[i-1] = GeneratePlaylist(i, public, songCollections[i-1]...)
	}

	return playlists
}

func CompareSerialisations(s1 string, s2 string) bool {

	re := regexp.MustCompile(`\r?\n\s+`)
	s1 = re.ReplaceAllString(s1, "")
	s2 = re.ReplaceAllString(s2, "")

	if s1 == s2 {
		return true
	}

	fmt.Println("")
	fmt.Println("--START--" + s1 + "--END--")
	fmt.Println("")
	fmt.Println("--START--" + s2 + "--END--")
	fmt.Println("")
	return false
}

func CheckSerialisation(o interface{}, xmlS string, jsonS string) error {

	// https://jsonformatter.curiousconcept.com/
	// https://www.liquid-technologies.com/online-xml-formatter

	m, err := xml.Marshal(o)
	if err != nil {
		return err
	}

	s := string(m)

	if !CompareSerialisations(xmlS, s) {
		return errors.New("Bad XML serialisation")
	}

	m, err = json.Marshal(o)
	if err != nil {
		return err
	}

	s = string(m)
	if !CompareSerialisations(jsonS, s) {
		return errors.New("Bad JSON serialisation")
	}

	return nil
}
