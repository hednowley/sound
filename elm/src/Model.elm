module Model exposing (Model, SocketModelWrap(..), getSocketModel, setSocketModel)

import Audio.Model
import AudioState
import Browser.Navigation exposing (Key)
import Config exposing (Config)
import Dict exposing (Dict)
import Entities.AlbumSummary exposing (AlbumSummaries)
import Entities.ArtistSummary exposing (ArtistSummaries)
import Entities.PlaylistSummary exposing (PlaylistSummaries)
import Entities.SongSummary exposing (SongSummary)
import Loadable exposing (Loadable(..))
import Nexus.Model
import Player.Model
import Routing exposing (Route)
import Socket.Model
import Url exposing (Url)


type alias Model =
    { key : Key
    , url : Url
    , username : String
    , password : String
    , message : String
    , token : Loadable String
    , isScanning : Bool
    , scanCount : Int
    , scanShouldUpdate : Bool
    , scanShouldDelete : Bool
    , playlists : PlaylistSummaries
    , nexus : Nexus.Model.Model
    , artists : ArtistSummaries
    , albums : AlbumSummaries
    , songs : Dict Int SongSummary
    , config : Config
    , route : Maybe Route
    , socket : SocketModelWrap
    , audio : Audio.Model.Model
    , player : Player.Model.Model
    }


{-| Type to avoid type recursion
-}
type SocketModelWrap
    = SocketModelWrap (Socket.Model.Model Model)


getSocketModel : Model -> Socket.Model.Model Model
getSocketModel model =
    let
        (SocketModelWrap s) =
            model.socket
    in
    s


setSocketModel : Model -> Socket.Model.Model Model -> Model
setSocketModel model socket =
    { model
        | socket = SocketModelWrap socket
    }
