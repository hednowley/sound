module Audio.Select exposing (getSongState)

import Audio.State exposing (State(..))
import Dict
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Routing exposing (Route(..))
import Song.Types exposing (SongId(..), getRawSongId)


getSongState : Model -> SongId -> Maybe State
getSongState model songId =
    Dict.get (getRawSongId songId) model.audio.songs
