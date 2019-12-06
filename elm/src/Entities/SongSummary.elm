module Entities.SongSummary exposing (SongSummary)

import Song.Types exposing (SongId)


type alias SongSummary =
    { id : SongId
    , name : String
    , track : Int
    }
