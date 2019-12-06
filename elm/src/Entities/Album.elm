module Entities.Album exposing (Album)

import Song.Types exposing (SongId)


type alias Album =
    { id : Int
    , artId : Maybe String
    , name : String
    , songs : List SongId
    }
