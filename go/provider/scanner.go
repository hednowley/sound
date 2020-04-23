package provider

import (
	"sync"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/dto"
)

// DAL (data access layer) allows high-level manipulation of application data.
type Scanner struct {
	providers []Provider
	dal       *dal.DAL
	hub       socket.IHub
}

// NewDAL constructs a new DAL.
func NewScanner(providers []Provider, dal *dal.DAL, hub socket.IHub) *Scanner {
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

	scanner.hub.Notify(dto.NewScanStatusNotification(provider.IsScanning(), provider.FileCount()))

	defer scanner.hub.Notify(dto.NewScanStatusNotification(provider.IsScanning(), provider.FileCount()))

	tokens := []string{}

	conn, err := scanner.dal.Db.GetConn()
	if err != nil {
		seelog.Errorf("Can't start '%v' scan: %v", providerID, err)
		return
	}
	defer conn.Release()

	err = provider.Iterate(func(token string) error {
		scanner.hub.Notify(dto.NewScanStatusNotification(provider.IsScanning(), provider.FileCount()))

		if delete {
			tokens = append(tokens, token)
		}

		existing, getSongErr := scanner.dal.Db.GetSongIdFromToken(conn, token, providerID)
		if getSongErr != nil {
			return getSongErr
		}

		if existing == nil || update {
			data, err2 := provider.GetInfo(token)
			if err2 != nil {
				seelog.Errorf("Cannot read music info for '%v': %v", token, err2)
				return nil
			}

			if existing == nil {
				seelog.Infof("Adding token '%v'", token)
			} else {
				seelog.Infof("Updating token '%v'", token)
			}

			putErr := scanner.dal.PutSong(conn, data, token, providerID, existing)
			if putErr != nil {
				return putErr
			}

		} else {
			seelog.Infof("Skipping token '%v'", token)
		}

		return nil
	})
	if err != nil {
		seelog.Errorf("Error during '%v' scan: %v", providerID, err)
	}

	if delete {
		seelog.Info("Deleting unscanned items")
		err = scanner.dal.Db.DeleteMissing(conn, tokens, providerID)

		if err != nil {
			seelog.Errorf("Error during '%v' scan: %v", providerID, err)
		}
	}

	seelog.Infof("Finished '%v' scan.", providerID)

}
