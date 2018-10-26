package provider

import (
	"github.com/cihub/seelog"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/entities"
)

// Provider gives access to a collection of music files.
type Provider interface {
	// Iterate through all files in the collection, calling the provided callback synchronously on each.
	Iterate(func(path string)) error

	// Returns information about the file at the given path.
	GetInfo(path string) (*entities.FileInfo, error)

	// The ID of this provider.
	ID() string

	// Whether this provider is in the middle of a scan.
	IsScanning() bool

	// How many files the provider has processed in the current scan.
	FileCount() int64

	// An identifier for the current scan
	ScanID() string
}

func NewProviders(config *config.Config) []Provider {
	ps := []Provider{}

	if len(config.Path) > 0 {
		p, err := NewFsProvider("fs", config.Path, config.Extensions)
		if err != nil {
			seelog.Errorf("Could not create file system provider: %v", err)
		} else {
			ps = append(ps, p)
		}
	}

	if len(config.BeetsDB) > 0 {
		p, err := NewBeetsProvider("beets", config.BeetsDB)
		if err != nil {
			seelog.Errorf("Could not create beets provider: %v", err)
		} else {
			ps = append(ps, p)
		}
	}

	return ps
}
