module Entities.AlbumSummary exposing (AlbumSummaries, AlbumSummary)

import Album.Types exposing (AlbumId)
import Dict exposing (Dict)


type alias AlbumSummary =
    { id : AlbumId
    , name : String
    , duration : Int
    , year : Maybe Int
    , artId : Maybe String
    }


type alias AlbumSummaries =
    Dict Int AlbumSummary
