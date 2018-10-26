package provider

import (
	"errors"

	"github.com/hednowley/sound/entities"
)

type MockProvider struct {
	files      []*entities.FileInfo
	id         string
	isScanning bool
	count      int
	scanID     string
}

func NewMockProvider(id string, files []*entities.FileInfo) *MockProvider {
	return &MockProvider{
		files:      files,
		id:         id,
		isScanning: false,
		count:      0,
	}
}

func (p *MockProvider) SetScanID(id string) {
	p.scanID = id
}

func (p *MockProvider) ScanID() string {
	return p.scanID
}

func (p *MockProvider) Iterate(callback func(path string)) error {
	p.isScanning = true
	p.count = 0
	for _, f := range p.files {
		callback(f.Path)
		p.count = p.count + 1
	}
	p.isScanning = false
	return nil
}

func (p *MockProvider) GetInfo(path string) (*entities.FileInfo, error) {
	for _, f := range p.files {
		if f.Path == path {
			return f, nil
		}
	}
	return nil, errors.New("Bad path")
}

func (p *MockProvider) ID() string {
	return p.id
}

func (p *MockProvider) IsScanning() bool {
	return p.isScanning
}

func (p *MockProvider) FileCount() int64 {
	return int64(p.count)
}
