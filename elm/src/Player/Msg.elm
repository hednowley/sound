module Player.Msg exposing (PlayerMsg(..))

import Album.Types exposing (AlbumId)
import Array exposing (Array)
import Player.Repeat exposing (Repeat)
import Playlist.Types exposing (PlaylistId)
import Song.Types exposing (SongId)


type PlayerMsg
    = Play SongId
    | PlayItem Int
    | Pause
    | Resume
    | Queue SongId
    | PlayAlbum AlbumId
    | PlayPlaylist PlaylistId
    | Next
    | Prev
    | SetShuffle Bool
    | Shuffled (Array SongId)
    | SetRepeat Repeat
