package services_test

import (
	"testing"

	"github.com/hednowley/sound/services"
)

func TestTagReader(t *testing.T) {
	d, err := services.GetMusicData("../testdata/music/1.mp3")
	if err != nil {
		t.Error(err.Error())
	}

	if d.Title != "Front Street (Instrumental)" {
		t.Error()
	}

	if d.Album != "A Day Wit The Homiez (CD, 2002, RonnieCash.com)" {
		t.Error()
	}

	if d.AlbumArtist != "1st Down" {
		t.Error()
	}

	if d.Artist != "1st Down" {
		t.Error()
	}

	if d.Year != 1995 {
		t.Error()
	}

	if d.Track != 4 {
		t.Error()
	}

	if d.Genre != "Hip-Hop" {
		t.Error()
	}

	if d.Size != 8973172 {
		t.Error()
	}

	if d.Disc != 0 {
		t.Error()
	}

	if d.Extension != "mp3" {
		t.Error()
	}
}
