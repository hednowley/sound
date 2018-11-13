package dto

import (
	"testing"
)

func TestGenres(t *testing.T) {

	genres := GenerateGenres(4)
	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	albums1 := GenerateAlbums(5, genres[0], artist, art)
	GenerateAlbums(3, genres[1], artist, art)
	GenerateAlbums(1, genres[2], artist, art)

	GenerateSongs(1, genres[0], albums1[0], art)
	GenerateSongs(2, genres[1], albums1[0], art)
	GenerateSongs(3, genres[2], albums1[0], art)

	DTO := NewGenres(genres)

	xml := `
	<genres>
		<genre songCount="1" albumCount="5">genre1</genre>
		<genre songCount="2" albumCount="3">genre2</genre>
		<genre songCount="3" albumCount="1">genre3</genre>
		<genre songCount="0" albumCount="0">genre4</genre>
	</genres>
	`

	json := `
	{
		"genre":[
		   {
			  "songCount":1,
			  "albumCount":5,
			  "value":"genre1"
		   },
		   {
			  "songCount":2,
			  "albumCount":3,
			  "value":"genre2"
		   },
		   {
			  "songCount":3,
			  "albumCount":1,
			  "value":"genre3"
		   },
		   {
			  "songCount":0,
			  "albumCount":0,
			  "value":"genre4"
		   }
		]
	 }
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
