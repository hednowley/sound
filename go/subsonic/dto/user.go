package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/config"
)

type User struct {
	XMLName             xml.Name `xml:"user" json:"-"`
	Username            string   `xml:"username,attr" json:"username"`
	Email               string   `xml:"email,attr" json:"email"`
	ScrobblingEnabled   bool     `xml:"scrobblingEnabled,attr" json:"scrobblingEnabled"`
	AdminRole           bool     `xml:"adminRole,attr" json:"adminRole"`
	SettingsRole        bool     `xml:"settingsRole,attr" json:"settingsRole"`
	DownloadRole        bool     `xml:"downloadRole,attr" json:"downloadRole"`
	UploadRole          bool     `xml:"uploadRole,attr" json:"uploadRole"`
	PlaylistRole        bool     `xml:"playlistRole,attr" json:"playlistRole"`
	CoverArtRole        bool     `xml:"coverArtRole,attr" json:"coverArtRole"`
	CommentRole         bool     `xml:"commentRole,attr" json:"commentRole"`
	PodcastRole         bool     `xml:"podcastRole,attr" json:"podcastRole"`
	StreamRole          bool     `xml:"streamRole,attr" json:"streamRole"`
	JukeboxRole         bool     `xml:"jukeboxRole,attr" json:"jukeboxRole"`
	ShareRole           bool     `xml:"shareRole,attr" json:"shareRole"`
	VideoConversionRole bool     `xml:"videoConversionRole,attr" json:"videoConversionRole"`
	Folder              []uint   `xml:"folder,attr" json:"folder"`
}

func NewUser(user config.User) User {
	return User{
		Username:     user.Username,
		Email:        user.Email,
		AdminRole:    true,
		SettingsRole: true,
		DownloadRole: true,
		UploadRole:   true,
		PlaylistRole: true,
		StreamRole:   true,
		Folder:       []uint{0},
	}
}
