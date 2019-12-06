package interfaces

import (
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/entities"
)

type DAL interface {
	PutPlaylist(id uint, name string, songIDs []uint) (uint, error)
	GetSong(id uint, genre bool, album bool, artist bool) (*dao.Song, error)
	GetAlbum(id uint, genre bool, artist bool, songs bool) (*dao.Album, error)
	GetArtPath(id string) (string, error)
	GetArtist(id uint) (*dao.Artist, error)
	GetGenre(name string) (*dao.Genre, error)
	GetPlaylist(id uint) (*dao.Playlist, error)
	GetSongFromToken(token string, providerID string) *dao.Song
	UpdatePlaylist(id uint, name string, comment string, public *bool, addedSongs []uint, removedSongs []uint) error
	DeletePlaylist(id uint) error
	GetAlbums(listType dao.AlbumList2Type, size uint, offset uint) []*dao.Album
	GetArtists(includeAlbums bool) []*dao.Artist
	GetGenres() []*dao.Genre
	GetPlaylists() []*dao.Playlist
	SynchroniseAlbum(id uint) (*dao.Album, error)
	SynchroniseArtist(id uint) error
	PutSong(song *dao.Song, data *entities.FileInfo) *dao.Song
	PutArt(art *entities.CoverArtData) *dao.Art
	Empty()
	DeleteMissing(tokens []string, providerID string)

	SearchArtists(query string, count uint, offset uint) []*dao.Artist
	SearchAlbums(query string, count uint, offset uint) []*dao.Album
	SearchSongs(query string, count uint, offset uint) []*dao.Song

	GetRandomSongs(size uint, from uint, to uint, genre string) []*dao.Song

	StarSong(songID uint, star bool) error
	StarAlbum(albumID uint, star bool) error
	StarArtist(artistID uint, star bool) error
}
