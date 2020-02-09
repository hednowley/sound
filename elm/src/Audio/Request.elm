module Audio.Request exposing (LoadRequest, makeLoadRequest)

import Song.Types exposing (SongId, getRawSongId)
import String exposing (fromInt)


type alias LoadRequest =
    { url : String
    , songId : Int
    }


makeLoadRequest : SongId -> LoadRequest
makeLoadRequest songId =
    { url = "/api/stream?id=" ++ fromInt (getRawSongId songId)
    , songId = getRawSongId songId
    }
