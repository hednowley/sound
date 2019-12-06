package provider

import (
	"os"

	"github.com/hednowley/sound/entities"
	"github.com/hednowley/sound/services"
)

// FsProvider scans music files inside a file system directory.
type FsProvider struct {
	id         string
	fileCount  int64
	isScanning bool
	path       string
	extensions []string
}

// NewFsProvider constructs a new provider.
func NewFsProvider(id string, path string, extensions []string) (*FsProvider, error) {
	return &FsProvider{
		id:         id,
		fileCount:  0,
		isScanning: false,
		path:       path,
		extensions: extensions,
	}, nil
}

// FileCount returns the number of files the provider has processed in the current scan.
func (p *FsProvider) FileCount() int64 {
	return p.fileCount
}

// IsScanning returns whether this provider is in the middle of a scan.
func (p *FsProvider) IsScanning() bool {
	return p.isScanning
}

// ID of this provider.
func (p *FsProvider) ID() string {
	return p.id
}

// Iterate through all files in the collection, calling the provided callback synchronously on each.
func (p *FsProvider) Iterate(callback func(token string)) error {
	p.isScanning = true
	p.fileCount = 0

	err := services.IterateFiles(p.path, p.extensions, func(path string, info *os.FileInfo) {
		// Use the path as a unique token
		callback(path)
	})

	if err != nil {
		p.isScanning = false
		return err
	}

	p.isScanning = false
	return nil
}

// GetInfo returns information about the file at the given path.
func (p *FsProvider) GetInfo(token string) (*entities.FileInfo, error) {
	return services.GetMusicData(token)
}
