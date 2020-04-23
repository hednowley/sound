package dto

import (
	"testing"

	"github.com/hednowley/sound/database"
)

func TestGenres(t *testing.T) {

	// genres := GenerateGenres(4)
	// art := GenerateArt(1)
	// artist := GenerateArtist(1, art)
	// albums1 := GenerateAlbums(5, &genres[0], artist, art)
	// GenerateAlbums(3, &genres[1], artist, art)
	// GenerateAlbums(1, &genres[2], artist, art)

	// GenerateSongs(1, &genres[0], &albums1[0], art)
	// GenerateSongs(2, &genres[1], &albums1[0], art)
	// GenerateSongs(3, &genres[2], &albums1[0], art)

	db := database.NewMock()
	conn, _ := db.GetConn()
	defer conn.Release()

	genres, err := db.GetGenres(conn)
	if err != nil {
		t.Error(err)
	}

	DTO := NewGenres(genres)

	xml := `
	<genres>
		<genre songCount="3" albumCount="3">genre_1</genre>
		<genre songCount="1" albumCount="1">genre_2</genre>
		<genre songCount="0" albumCount="0">genre_4</genre>
	</genres>
	`

	json := `
	{
		"genre":[
			{
				"songCount":3,
				"albumCount":3,
				"value":"genre_1"
			},
			{
				"songCount":1,
				"albumCount":1,
				"value":"genre_2"
			},
			{
				"songCount":0,
				"albumCount":0,
				"value":"genre_4"
			}
		]
	}
	`

	err = CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
