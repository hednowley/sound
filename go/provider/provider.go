package provider

import (
	"fmt"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/entities"
)

// Provider gives access to a collection of music files.
type Provider interface {
	// Iterate through all files in the collection, calling the provided
	// callback synchronously on each file's unique token.
	Iterate(func(token string) error) error

	// Returns information about the file with the given token.
	GetInfo(token string) (*entities.FileInfo, error)

	// The ID of this provider.
	ID() string

	// Whether this provider is in the middle of a scan.
	IsScanning() bool

	// How many files the provider has processed in the current scan.
	FileCount() int64
}

func addProvider(providers []Provider, provider Provider) ([]Provider, error) {
	for _, p := range providers {
		if p.ID() == provider.ID() {
			return providers, fmt.Errorf("Duplicate provider name (%v)", provider.ID())
		}
	}
	providers = append(providers, provider)
	return providers, nil
}

func NewProviders(config *config.Config) []Provider {
	ps := []Provider{}

	for _, fs := range config.FileSystemProviders {

		p, err := NewFsProvider(fs.Name, fs.Path, fs.Extensions)
		if err != nil {
			seelog.Errorf("Could not create file system provider: %v", err)
			continue
		}

		ps, err = addProvider(ps, p)
		if err != nil {
			seelog.Errorf("Could not create file system provider: %v", err)
		}
	}

	for _, b := range config.BeetsProviders {
		p, err := NewBeetsProvider(b.Name, b.Database)
		if err != nil {
			seelog.Errorf("Could not create beets provider: %v", err)
			continue
		}

		ps, err = addProvider(ps, p)
		if err != nil {
			seelog.Errorf("Could not create beets provider: %v", err)
		}
	}

	return ps
}
