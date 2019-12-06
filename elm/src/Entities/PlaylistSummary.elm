module Entities.PlaylistSummary exposing (PlaylistSummaries, PlaylistSummary)

import Dict exposing (Dict)


type alias PlaylistSummary =
    { id : Int
    , name : String
    }


type alias PlaylistSummaries =
    Dict Int PlaylistSummary
