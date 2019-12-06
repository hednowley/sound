module Artist.Fetch exposing (fetchArtist)

import Artist.Types exposing (ArtistId, getRawArtistId)
import Entities.Artist exposing (Artist)
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg)
import Nexus.Fetch exposing (fetch)
import Socket.DTO.Artist exposing (convert, decode)
import Song.Types exposing (SongId(..))
import Types exposing (Update)


noOp : Socket.DTO.Artist.Artist -> Model -> Model
noOp _ model =
    model


fetchArtist : Maybe (Artist -> Update Model Msg) -> ArtistId -> Update Model Msg
fetchArtist maybeCallback =
    fetch
        getRawArtistId
        "getArtist"
        decode
        noOp
        convert
        { get = .artists
        , set = \repo -> \m -> { m | artists = repo }
        }
        maybeCallback
