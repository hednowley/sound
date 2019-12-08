package provider

import (
	"sync"
	"time"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/ws"
	"github.com/hednowley/sound/ws/dto"
)

// DAL (data access layer) allows high-level manipulation of application data.
type Scanner struct {
	providers []Provider
	dal       *dal.DAL
	hub       ws.IHub
}

// NewDAL constructs a new DAL.
func NewScanner(providers []Provider, dal *dal.DAL, hub ws.IHub) *Scanner {
	return &Scanner{
		providers: providers,
		dal:       dal,
		hub:       hub,
	}
}

func (scanner *Scanner) GetScanFileCount() int64 {
	count := int64(0)
	for _, p := range scanner.providers {
		count += p.FileCount()
	}
	return count
}

func (scanner *Scanner) GetScanStatus() bool {
	scanning := false
	for _, p := range scanner.providers {
		scanning = scanning || p.IsScanning()
	}
	return scanning
}

// StartAllScans asks all providers to start scanning in parallel.
func (scanner *Scanner) StartAllScans(update bool, delete bool) {
	seelog.Info("Starting all scans.")
	var wg sync.WaitGroup
	for _, p := range scanner.providers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			scanner.startScan(p, update, delete)
		}()
	}
	wg.Wait()
}

func (scanner *Scanner) startScan(provider Provider, update bool, delete bool) {
	providerID := provider.ID()

	if provider.IsScanning() {
		seelog.Infof("Skipped '%v' scan as one is already in progress.", providerID)
		return
	}

	seelog.Infof("Started '%v' scan.", providerID)
	synch := NewSynchroniser(scanner.dal, 10)

	scanner.hub.Notify(dto.NewScanStatusNotification(provider.IsScanning(), provider.FileCount()))

	tokens := []string{}

	err := provider.Iterate(func(token string) {
		scanner.hub.Notify(dto.NewScanStatusNotification(provider.IsScanning(), provider.FileCount()))

		if delete {
			tokens = append(tokens, token)
		}

		s := scanner.dal.GetSongFromToken(token, providerID)
		if s == nil || update {
			data, err2 := provider.GetInfo(token)
			if err2 != nil {
				seelog.Errorf("Cannot read music info for '%v': %v", token, err2)
				return
			}

			if s == nil {
				seelog.Infof("Adding token '%v'", token)
				now := time.Now()
				s = &dao.Song{
					Created:    &now,
					ProviderID: providerID,
					Token:      token,
				}
			} else {
				seelog.Infof("Updating token '%v'", token)
				synch.Notify(s.AlbumID) // Notify of potential change to old album
			}

			s = scanner.dal.PutSong(s, data)

			// Notify of change to new album
			synch.Notify(s.AlbumID)

		} else {
			seelog.Infof("Skipping token '%v'", token)
		}
	})
	if err != nil {
		seelog.Errorf("Error during '%v' scan: %v", providerID, err)
	}

	// Make any remaining updates
	synch.Flush()

	if delete {
		seelog.Info("Deleting unscanned items")
		scanner.dal.DeleteMissing(tokens, providerID)

		// Find unscanned songs from same provider

		// Find albums etc with no songs

	}

	seelog.Infof("Finished '%v' scan.", providerID)

	scanner.hub.Notify(dto.NewScanStatusNotification(provider.IsScanning(), provider.FileCount()))
}
