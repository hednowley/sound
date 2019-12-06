module Playlist.Types exposing (PlaylistId(..), getRawPlaylistId)


type PlaylistId
    = PlaylistId Int


getRawPlaylistId : PlaylistId -> Int
getRawPlaylistId playlistId =
    let
        (PlaylistId raw) =
            playlistId
    in
    raw
