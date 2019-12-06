module Socket.DTO.Album exposing (Album, convert, decode)

import Entities.Album
import Json.Decode exposing (Decoder, field, int, list, map6, maybe, string)
import Socket.DTO.SongSummary exposing (SongSummary)
import Song.Types exposing (SongId(..))


type alias Album =
    { id : Int
    , artId : Maybe String
    , name : String
    , duration : Int
    , year : Maybe Int
    , songs : List SongSummary
    }


decode : Decoder Album
decode =
    map6 Album
        (field "id" int)
        (maybe <| field "coverArt" string)
        (field "name" string)
        (field "duration" int)
        (maybe <| field "year" int)
        (field "songs" <| list Socket.DTO.SongSummary.decode)


convert : Album -> Entities.Album.Album
convert album =
    { id = album.id
    , artId = album.artId
    , name = album.name
    , songs = List.map (.id >> SongId) album.songs
    }
