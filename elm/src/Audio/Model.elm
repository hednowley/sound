module Audio.Model exposing (Model, emptyModel)

import Audio.State exposing (State)
import Dict exposing (Dict)


emptyModel : Model
emptyModel =
    { songs = Dict.empty
    }


type alias Model =
    { songs : Dict Int State
    }
