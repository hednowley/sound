module Socket.DTO.ArtistSummary exposing (ArtistSummary, convert, decode)

import Artist.Types exposing (ArtistId(..))
import Entities.ArtistSummary
import Json.Decode exposing (Decoder, field, int, map2, string)


type alias ArtistSummary =
    { id : Int
    , name : String
    }


decode : Decoder ArtistSummary
decode =
    map2 ArtistSummary
        (field "id" int)
        (field "name" string)


convert : ArtistSummary -> Entities.ArtistSummary.ArtistSummary
convert album =
    { id = ArtistId album.id
    , name = album.name
    }
