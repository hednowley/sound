module Playlist.Select exposing (getPlaylist, getPlaylistSongs)

import Dict
import Entities.Playlist exposing (Playlist)
import Entities.SongSummary exposing (SongSummary)
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Playlist.Types exposing (PlaylistId, getRawPlaylistId)
import Song.Select exposing (getSong)


getPlaylist : PlaylistId -> Model -> Loadable Playlist
getPlaylist id model =
    Dict.get (getRawPlaylistId id) model.nexus.playlists |> Maybe.withDefault Absent


getPlaylistSongs : Playlist -> Model -> List (Maybe SongSummary)
getPlaylistSongs playlist model =
    List.map (getSong model) playlist.songs
