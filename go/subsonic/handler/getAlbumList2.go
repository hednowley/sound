package handler

import (
	"net/url"
	"strings"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
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
func NewGetAlbumList2Handler(dal *dal.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		listType := parseListType(params.Get("type"))
		if listType == nil {
			return api.NewErrorReponse(dto.Generic, "Unknown type.")
		}

		size := util.ParseUint(params.Get("size"), 10)
		if size > 500 {
			return api.NewErrorReponse(dto.Generic, "Invalid size.")
		}

		offset := util.ParseUint(params.Get("offset"), 0)

		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}
		defer conn.Release()

		albums, err := dal.Db.GetAlbums(conn, *listType, size, offset)
		if err != nil {
			return api.NewErrorReponse(0, "Error")
		}
		return api.NewSuccessfulReponse(dto.NewAlbumList2(albums))
	}
}
