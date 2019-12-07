package provider

import (
	"github.com/cihub/seelog"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/util"
)

// Synchroniser efficiently keeps artist and album information in sync
// with song information.
type Synchroniser struct {
	dal             *dal.DAL
	albumQueue      []uint
	flushThreshhold int
}

// NewSynchroniser returns a new synchroniser. flushthreshhold is how many
// albums need to be changed before they and their artists are synced. A high
// value will be more efficient if albums arrive roughly grouped by artist.
func NewSynchroniser(dal *dal.DAL, flushThreshhold int) *Synchroniser {
	return &Synchroniser{
		dal:             dal,
		flushThreshhold: flushThreshhold,
	}
}

// Notify tells the synchroniser that the songs of the album with
// the given ID may have changed.
func (s *Synchroniser) Notify(album uint) {

	// Don't add duplicates to the queue
	if util.Contains(s.albumQueue, album) {
		return
	}

	if s.flushThreshhold != 0 && len(s.albumQueue) > s.flushThreshhold {
		s.Flush()
	}

	s.albumQueue = append(s.albumQueue, album)
}

// Flush synchronises the pending albums and artists.
func (s *Synchroniser) Flush() error {

	var artists []uint

	seelog.Info("Synchroniser flushing...")

	for _, album := range s.albumQueue {
		a, err := s.dal.SynchroniseAlbum(album)
		if err != nil {
			return err
		}

		// Don't add duplicates to the queue
		if !util.Contains(artists, a.ArtistID) {
			artists = append(artists, a.ArtistID)
		}
	}

	for _, artist := range artists {
		s.dal.SynchroniseArtist(artist)
	}

	s.albumQueue = nil
	seelog.Info("Synchroniser flushed.")
	return nil
}
