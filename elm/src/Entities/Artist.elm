module Entities.Artist exposing (Artist)

import Artist.Types exposing (ArtistId)
import Entities.AlbumSummary exposing (AlbumSummary)


type alias Artist =
    { id : ArtistId
    , name : String
    , albums : List AlbumSummary
    }
