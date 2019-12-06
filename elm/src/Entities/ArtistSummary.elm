module Entities.ArtistSummary exposing (ArtistSummaries, ArtistSummary)

import Artist.Types exposing (ArtistId)
import Dict exposing (Dict)


type alias ArtistSummary =
    { id : ArtistId
    , name : String
    }


type alias ArtistSummaries =
    Dict Int ArtistSummary
