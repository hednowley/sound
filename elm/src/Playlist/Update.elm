module Playlist.Update exposing (playPlaylist)

import Audio.AudioMsg exposing (AudioMsg(..))
import Entities.Playlist exposing (Playlist)
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg(..))
import Player.Actions exposing (replacePlaylist)
import Playlist.Fetch exposing (fetchPlaylist)
import Playlist.Types exposing (PlaylistId)
import Types exposing (Update)


playLoadedPlaylist : Playlist -> Update Model Msg
playLoadedPlaylist playlist =
    replacePlaylist playlist.songs


playPlaylist : PlaylistId -> Update Model Msg
playPlaylist playlistId =
    fetchPlaylist
        (Just playLoadedPlaylist)
        playlistId
