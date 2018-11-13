package handler

import (
	"net/url"
	"strings"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func parseListType(param string) *dao.AlbumList2Type {
	param = strings.ToLower(param)
	var t dao.AlbumList2Type
	if param == "random" {
		t = dao.Random
		return &t
	} else if param == "newest" {
		t = dao.Newest
		return &t
	} else if param == "frequent" {
		t = dao.Frequent
		return &t
	} else if param == "recent" {
		t = dao.Recent
		return &t
	} else if param == "starred" {
		t = dao.Starred
		return &t
	} else if param == "alphabeticalbyname" {
		t = dao.AlphabeticalByName
		return &t
	} else if param == "alphabeticalbyartist" {
		t = dao.AlphabeticalByArtist
		return &t
	} else if param == "byyear" {
		t = dao.ByYear
		return &t
	} else if param == "bygenre" {
		t = dao.ByGenre
		return &t
	}

	return nil
}

// NewGetAlbumList2Handler is a handler for getting information about a sample albums.
func NewGetAlbumList2Handler(database *dal.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		listType := parseListType(params.Get("type"))
		if listType == nil {
			return api.NewErrorReponse(dto.Generic, "Unknown type.")
		}

		size := api.ParseUint(params.Get("size"), 10)
		if size > 500 {
			return api.NewErrorReponse(dto.Generic, "Invalid size.")
		}

		offset := api.ParseUint(params.Get("offset"), 0)

		albums := database.GetAlbums(*listType, size, offset)
		return api.NewSuccessfulReponse(dto.NewAlbumList2(albums))
	}
}
