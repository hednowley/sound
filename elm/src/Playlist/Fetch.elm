module Playlist.Fetch exposing (fetchPlaylist)

import Entities.Playlist exposing (Playlist)
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg)
import Nexus.Fetch exposing (fetch)
import Playlist.Types exposing (PlaylistId, getRawPlaylistId)
import Socket.DTO.Playlist exposing (convert, decode)
import Socket.DTO.SongSummary exposing (convertMany)
import Song.Types exposing (SongId(..), getRawSongId)
import Types exposing (Update)
import Util exposing (insertMany)


fetchPlaylist : Maybe (Playlist -> Update Model Msg) -> PlaylistId -> Update Model Msg
fetchPlaylist maybeCallback =
    fetch
        getRawPlaylistId
        "getPlaylist"
        decode
        saveSongs
        convert
        { get = .playlists
        , set = \repo -> \m -> { m | playlists = repo }
        }
        maybeCallback


saveSongs : Socket.DTO.Playlist.Playlist -> Model -> Model
saveSongs playlist model =
    let
        songs =
            convertMany playlist.songs
    in
    { model
        | songs =
            insertMany
                (.id >> getRawSongId)
                identity
                songs
                model.songs
    }
