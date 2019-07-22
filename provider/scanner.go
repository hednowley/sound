package provider

import (
	"sync"
	"time"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/ws/dto"
)

// DAL (data access layer) allows high-level manipulation of application data.
type Scanner struct {
	providers []Provider
	dal       interfaces.DAL
	hub       interfaces.Hub
}

// NewDAL constructs a new DAL.
func NewScanner(providers []Provider, dal interfaces.DAL, hub interfaces.Hub) *Scanner {
	return &Scanner{
		providers: providers,
		dal:       dal,
		hub:       hub,
	}
}

func (dal *Scanner) GetScanFileCount() int64 {
	count := int64(0)
	for _, p := range dal.providers {
		count += p.FileCount()
	}
	return count
}

func (dal *Scanner) GetScanStatus() bool {
	scanning := false
	for _, p := range dal.providers {
		scanning = scanning || p.IsScanning()
	}
	return scanning
}

// StartAllScans asks all providers to start scanning in parallel.
func (dal *Scanner) StartAllScans(update bool, delete bool) {
	seelog.Info("Starting all scans.")
	var wg sync.WaitGroup
	for _, p := range dal.providers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dal.startScan(p, update, delete)
		}()
	}
	wg.Wait()
}

func (dal *Scanner) startScan(provider Provider, update bool, delete bool) {
	providerID := provider.ID()

	if provider.IsScanning() {
		seelog.Infof("Skipped '%v' scan as one is already in progress.", providerID)
		return
	}

	seelog.Infof("Started '%v' scan.", providerID)
	synch := NewSynchroniser(dal.dal, 10)

	dal.hub.Notify(dto.NewScanStatusNotification(provider.IsScanning(), provider.FileCount()))

	tokens := []string{}

	err := provider.Iterate(func(token string) {
		dal.hub.Notify(dto.NewScanStatusNotification(provider.IsScanning(), provider.FileCount()))

		if delete {
			tokens = append(tokens, token)
		}

		s := dal.dal.GetSongFromToken(token, providerID)
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

			s = dal.dal.PutSong(s, data)

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
		dal.dal.DeleteMissing(tokens, providerID)

		// Find unscanned songs from same provider

		// Find albums etc with no songs

	}

	seelog.Infof("Finished '%v' scan.", providerID)

	dal.hub.Notify(dto.NewScanStatusNotification(provider.IsScanning(), provider.FileCount()))
}
