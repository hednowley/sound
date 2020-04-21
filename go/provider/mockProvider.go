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
}

func NewMockProvider(id string, files []*entities.FileInfo) *MockProvider {
	return &MockProvider{
		files:      files,
		id:         id,
		isScanning: false,
		count:      0,
	}
}

func (p *MockProvider) Iterate(callback func(path string) error) error {
	p.isScanning = true
	p.count = 0

	var callbackErr error
	for _, f := range p.files {
		callbackErr = callback(f.Path)
		p.count = p.count + 1

		if callbackErr != nil {
			break
		}
	}
	p.isScanning = false
	return callbackErr
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
