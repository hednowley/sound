package dto

import (
	"testing"
)

func TestSong(t *testing.T) {

	genre := GenerateGenre(1)
	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	album := GenerateAlbum(1, genre, artist, art)
	song := GenerateSong(1, genre, album, art)

	DTO := NewSong(song)

	xml := `
	<song id="1" title="title1" album="album1" artist="artist1" track="1" year="1901" genre="genre1" coverArt="1.jpg" size="100000" contentType="" suffix="mp1" duration="0" bitRate="0" path="D:\music\1.mp3" isVideo="false" playCount="0" discNumber="1" created="2001-08-15T00:00:00Z" albumId="1" artistId="1" type="music"></song>
	`

	json := `
	{
		"id":"1",
		"title":"title1",
		"album":"album1",
		"artist":"artist1",
		"track":1,
		"year":1901,
		"genre":"genre1",
		"coverArt":"1.jpg",
		"size":100000,
		"contentType":"",
		"suffix":"mp1",
		"duration":0,
		"bitRate":0,
		"path":"D:\\music\\1.mp3",
		"isVideo":false,
		"playCount":0,
		"discNumber":1,
		"created":"2001-08-15T00:00:00Z",
		"albumId":"1",
		"artistId":"1",
		"type":"music"
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSongWithoutYear(t *testing.T) {

	genre := GenerateGenre(1)
	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	album := GenerateAlbum(1, genre, artist, art)
	song := GenerateSong(1, genre, album, art)
	song.Year = 0

	DTO := NewSong(song)

	xml := `
	<song id="1" title="title1" album="album1" artist="artist1" track="1" genre="genre1" coverArt="1.jpg" size="100000" contentType="" suffix="mp1" duration="0" bitRate="0" path="D:\music\1.mp3" isVideo="false" playCount="0" discNumber="1" created="2001-08-15T00:00:00Z" albumId="1" artistId="1" type="music"></song>
	`

	json := `
	{
		"id":"1",
		"title":"title1",
		"album":"album1",
		"artist":"artist1",
		"track":1,
		"genre":"genre1",
		"coverArt":"1.jpg",
		"size":100000,
		"contentType":"",
		"suffix":"mp1",
		"duration":0,
		"bitRate":0,
		"path":"D:\\music\\1.mp3",
		"isVideo":false,
		"playCount":0,
		"discNumber":1,
		"created":"2001-08-15T00:00:00Z",
		"albumId":"1",
		"artistId":"1",
		"type":"music"
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSongWithoutArt(t *testing.T) {

	genre := GenerateGenre(1)
	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	album := GenerateAlbum(1, genre, artist, art)
	song := GenerateSong(1, genre, album, nil)

	DTO := NewSong(song)

	xml := `
	<song id="1" title="title1" album="album1" artist="artist1" track="1" year="1901" genre="genre1" size="100000" contentType="" suffix="mp1" duration="0" bitRate="0" path="D:\music\1.mp3" isVideo="false" playCount="0" discNumber="1" created="2001-08-15T00:00:00Z" albumId="1" artistId="1" type="music"></song>
	`

	json := `
	{
		"id":"1",
		"title":"title1",
		"album":"album1",
		"artist":"artist1",
		"track":1,
		"year":1901,
		"genre":"genre1",
		"size":100000,
		"contentType":"",
		"suffix":"mp1",
		"duration":0,
		"bitRate":0,
		"path":"D:\\music\\1.mp3",
		"isVideo":false,
		"playCount":0,
		"discNumber":1,
		"created":"2001-08-15T00:00:00Z",
		"albumId":"1",
		"artistId":"1",
		"type":"music"
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSongWithoutGenre(t *testing.T) {

	genre := GenerateGenre(1)
	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	album := GenerateAlbum(1, genre, artist, art)
	song := GenerateSong(1, nil, album, art)

	DTO := NewSong(song)

	xml := `
	<song id="1" title="title1" album="album1" artist="artist1" track="1" year="1901" coverArt="1.jpg" size="100000" contentType="" suffix="mp1" duration="0" bitRate="0" path="D:\music\1.mp3" isVideo="false" playCount="0" discNumber="1" created="2001-08-15T00:00:00Z" albumId="1" artistId="1" type="music"></song>
	`

	json := `
	{
		"id":"1",
		"title":"title1",
		"album":"album1",
		"artist":"artist1",
		"track":1,
		"year":1901,
		"coverArt":"1.jpg",
		"size":100000,
		"contentType":"",
		"suffix":"mp1",
		"duration":0,
		"bitRate":0,
		"path":"D:\\music\\1.mp3",
		"isVideo":false,
		"playCount":0,
		"discNumber":1,
		"created":"2001-08-15T00:00:00Z",
		"albumId":"1",
		"artistId":"1",
		"type":"music"
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
