module Song.Select exposing (getSong)

import AudioState exposing (State(..))
import Dict
import Entities.SongSummary exposing (SongSummary)
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Routing exposing (Route(..))
import Song.Types exposing (SongId(..))


getSong : Model -> SongId -> Maybe SongSummary
getSong model songId =
    let
        (SongId id) =
            songId
    in
    Dict.get id model.songs
