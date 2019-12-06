module Player.Model exposing (Model, emptyModel)

import Array exposing (Array)
import Song.Types exposing (SongId)


emptyModel : Model
emptyModel =
    { shuffle = False
    , repeat = All
    , playlist = Array.empty
    , unshuffledPlaylist = Array.empty
    , playing = Nothing
    }


type alias Model =
    { shuffle : Bool
    , repeat : Repeat
    , playlist : Array SongId
    , unshuffledPlaylist : Array SongId
    , playing : Maybe Int
    }


type Repeat
    = None
    | One
    | All
