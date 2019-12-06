module Socket.DTO.AlbumSummary exposing (AlbumSummary, convert, decode)

import Album.Types exposing (AlbumId(..))
import Entities.AlbumSummary
import Json.Decode exposing (Decoder, field, int, map5, maybe, string)


type alias AlbumSummary =
    { id : Int
    , name : String
    , duration : Int
    , year : Maybe Int
    , artId : Maybe String
    }


decode : Decoder AlbumSummary
decode =
    map5 AlbumSummary
        (field "id" int)
        (field "name" string)
        (field "duration" int)
        (maybe <| field "year" int)
        (maybe <| field "coverArt" string)


convert : AlbumSummary -> Entities.AlbumSummary.AlbumSummary
convert album =
    { id = AlbumId album.id
    , name = album.name
    , duration = album.duration
    , year = album.year
    , artId = album.artId
    }
