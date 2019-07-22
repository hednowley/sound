package interfaces

import (
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/entities"
)

type DAL interface {
	PutPlaylist(id uint, name string, songIDs []uint) (uint, error)
	GetSong(id uint, genre bool, album bool, artist bool, art bool) (*dao.Song, error)
	GetAlbum(id uint, genre bool, artist bool, songs bool) (*dao.Album, error)
	GetArt(id uint) (*dao.Art, error)
	GetArtist(id uint) (*dao.Artist, error)
	GetGenre(name string) (*dao.Genre, error)
	GetPlaylist(id uint) (*dao.Playlist, error)
	GetSongFromToken(token string, providerID string) *dao.Song
	UpdatePlaylist(id uint, name string, comment string, public *bool, addedSongs []uint, removedSongs []uint) error
	DeletePlaylist(id uint) error
	GetAlbums(listType dao.AlbumList2Type, size uint, offset uint) []*dao.Album
	GetArtists() []*dao.Artist
	GetGenres() []*dao.Genre
	GetPlaylists() []*dao.Playlist
	SynchroniseAlbum(id uint) (*dao.Album, error)
	SynchroniseArtist(id uint) error
	PutSong(song *dao.Song, data *entities.FileInfo) *dao.Song
	Empty()
	DeleteMissing(tokens []string, providerID string)
}
