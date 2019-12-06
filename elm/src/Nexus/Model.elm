module Nexus.Model exposing (Model, empty)

import Dict exposing (Dict)
import Entities.Album exposing (Album)
import Entities.Artist exposing (Artist)
import Entities.Playlist exposing (Playlist)
import Loadable exposing (Loadable(..))


type alias Model =
    { playlists : Dict Int (Loadable Playlist)
    , artists : Dict Int (Loadable Artist)
    , albums : Dict Int (Loadable Album)
    }


empty : Model
empty =
    { playlists = Dict.empty
    , artists = Dict.empty
    , albums = Dict.empty
    }
