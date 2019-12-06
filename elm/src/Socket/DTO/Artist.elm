module Socket.DTO.Artist exposing (Artist, convert, decode)

import Artist.Types exposing (ArtistId(..))
import Entities.Artist
import Json.Decode exposing (Decoder, field, int, list, map3, string)
import Socket.DTO.AlbumSummary exposing (AlbumSummary)


type alias Artist =
    { id : Int
    , name : String
    , albums : List AlbumSummary
    }


decode : Decoder Artist
decode =
    map3 Artist
        (field "id" int)
        (field "name" string)
        (field "albums"
            (list Socket.DTO.AlbumSummary.decode)
        )


convert : Artist -> Entities.Artist.Artist
convert artist =
    { id = ArtistId artist.id
    , name = artist.name
    , albums = List.map Socket.DTO.AlbumSummary.convert artist.albums
    }
