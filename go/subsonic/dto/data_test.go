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

func GenerateGenres(count int) []dao.Genre {
	g := make([]dao.Genre, count)
	for i := 1; i < count+1; i++ {
		g[i-1] = *GenerateGenre(i)
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

func GenerateArts(count int) []dao.Art {
	arts := make([]dao.Art, count)
	for i := 1; i < count+1; i++ {
		arts[i-1] = *GenerateArt(i)
	}
	return arts
}

func GenerateSong(i int, genre *dao.Genre, album *dao.Album, art *dao.Art) *dao.Song {
	t := time.Date(2000+i, time.August, 15, 0, 0, 0, 0, time.UTC)

	var artPath string
	if art != nil {
		artPath = art.Path
	}

	var genreName string
	if genre != nil {
		genreName = genre.Name
	}

	s := dao.Song{
		ID:            uint(i),
		Title:         fmt.Sprintf("title%v", i),
		Artist:        album.ArtistName,
		Track:         i,
		Disc:          i,
		AlbumID:       album.ID,
		Art:           artPath,
		Path:          fmt.Sprintf("D:\\music\\%v.mp3", i),
		Size:          int64(i * 100000),
		Year:          1900 + i,
		Created:       &t,
		GenreName:     genreName,
		AlbumArtistID: album.ArtistID,
		AlbumName:     album.Name,
	}

	return &s
}

func GenerateSongs(count int, genre *dao.Genre, album *dao.Album, art *dao.Art) []dao.Song {
	songs := make([]dao.Song, count)
	for i := 0; i < count; i++ {
		songs[i] = *GenerateSong(i+1, genre, album, art)
	}
	return songs
}

func GenerateAlbum(albumId int, genre *dao.Genre, artist *dao.Artist, art *dao.Art) *dao.Album {
	t := time.Date(2000+albumId, time.August, 15, 0, 0, 0, 0, time.UTC)

	a := dao.Album{
		ID:         uint(albumId),
		Name:       fmt.Sprintf("album%v", albumId),
		ArtistID:   artist.ID,
		Created:    &t,
		ArtistName: artist.Name,
		Duration:   albumId,
	}

	artist.AlbumCount++

	return &a
}

func GenerateAlbums(count int, genre *dao.Genre, artist *dao.Artist, art *dao.Art) []dao.Album {
	albums := make([]dao.Album, count)
	for i := 1; i < count+1; i++ {
		albums[i-1] = *GenerateAlbum(i, genre, artist, art)
	}
	return albums
}

func GenerateArtist(i int, art *dao.Art) *dao.Artist {

	a := dao.Artist{
		ID:       uint(i),
		Name:     fmt.Sprintf("artist%v", i),
		Duration: i,
		Arts:     []string{art.Path},
	}

	return &a
}

func GenerateArtists(count int, art *dao.Art) []dao.Artist {
	artists := make([]dao.Artist, count)
	for i := 1; i < count+1; i++ {
		artists[i-1] = *GenerateArtist(i, art)
	}
	return artists
}

func GeneratePlaylist(i int, public bool) *dao.Playlist {

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

	return &p
}

func GeneratePlaylists(public bool, count int) []*dao.Playlist {

	playlists := make([]*dao.Playlist, count)
	for i := 1; i < count+1; i++ {
		playlists[i-1] = GeneratePlaylist(i, public)
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
