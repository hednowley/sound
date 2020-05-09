package config_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hednowley/sound/config"
)

func TestConfig(t *testing.T) {
	// Get test data path (found relative to this file)
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	dataPath := filepath.Join(dir, "..", "testdata", "config.yaml")

	c, err := config.NewConfig(dataPath)()

	if err != nil {
		t.Error()
	}

	if len(c.ArtSizes) != 2 {
		t.Error()
	}

	if c.Port != 3684 ||
		c.Secret != "changeme" ||
		c.ArtPath != "~/temp/art" ||
		c.Db != "host=localhost port=5432 user=postgres password=sound dbname=sound sslmode=disable" ||
		c.LogConfig != "log-config.xml" ||
		c.AccessControlAllowOrigin != "*" ||
		c.WebsocketTicketExpiry != 30 {
		t.Error()
	}

	if len(c.FileSystemProviders) != 2 {
		t.Error()
	}

	if c.FileSystemProviders[0].Path != "~/temp/my music" ||
		c.FileSystemProviders[0].Name != "Gertrude's bangers" ||
		len(c.FileSystemProviders[0].Extensions) != 2 ||
		c.FileSystemProviders[0].Extensions[0] != "mp3" ||
		c.FileSystemProviders[0].Extensions[1] != "flac" {
		t.Error()
	}

	if len(c.BeetsProviders) != 1 ||
		c.BeetsProviders[0].Database != "~/beets/beetslib.blb" ||
		c.BeetsProviders[0].Name != "My beets music" {
		t.Error()
	}
}
