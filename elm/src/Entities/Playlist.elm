module Entities.Playlist exposing (Playlist)

import Song.Types exposing (SongId)


type alias Playlist =
    { id : Int
    , name : String
    , songs : List SongId
    }
