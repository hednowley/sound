package provider

import (
	"github.com/cihub/seelog"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/entities"
)

// Provider gives access to a collection of music files.
type Provider interface {
	// Iterate through all files in the collection, calling the provided
	// callback synchronously on each file's unique token.
	Iterate(func(token string)) error

	// Returns information about the file with the given token.
	GetInfo(token string) (*entities.FileInfo, error)

	// The ID of this provider.
	ID() string

	// Whether this provider is in the middle of a scan.
	IsScanning() bool

	// How many files the provider has processed in the current scan.
	FileCount() int64
}

func NewProviders(config *config.Config) []Provider {
	ps := []Provider{}

	for _, fs := range config.FileSystemProviders {
		p, err := NewFsProvider("fs", fs.Path, fs.Extensions)
		if err != nil {
			seelog.Errorf("Could not create file system provider: %v", err)
		} else {
			ps = append(ps, p)
		}
	}

	for _, b := range config.BeetsProviders {
		p, err := NewBeetsProvider("beets", b.Database)
		if err != nil {
			seelog.Errorf("Could not create beets provider: %v", err)
		} else {
			ps = append(ps, p)
		}
	}

	return ps
}
