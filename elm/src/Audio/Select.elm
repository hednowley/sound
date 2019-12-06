module Audio.Select exposing (getSongState)

import AudioState exposing (State(..))
import Dict
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Routing exposing (Route(..))
import Song.Types exposing (SongId(..), getRawSongId)


getSongState : SongId -> Model -> Maybe State
getSongState songId model =
    Dict.get (getRawSongId songId) model.songCache
