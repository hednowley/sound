module Player.Model exposing (Model, emptyModel)

import Array exposing (Array)
import Player.Repeat exposing (Repeat(..))
import Song.Types exposing (SongId)


emptyModel : Model
emptyModel =
    { shuffle = False
    , repeat = None
    , playlist = Array.empty
    , unshuffledPlaylist = Array.empty
    , playing = Nothing
    , isPaused = True
    }


type alias Model =
    { shuffle : Bool
    , repeat : Repeat
    , playlist : Array SongId
    , unshuffledPlaylist : Array SongId
    , playing : Maybe Int
    , isPaused : Bool
    }
