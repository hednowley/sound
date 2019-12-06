module Updaters exposing
    ( logOut
    , onUrlChange
    )

import Album.Fetch exposing (fetchAlbum)
import Artist.Fetch exposing (fetchArtist)
import Artist.Types exposing (ArtistId(..))
import Audio.Select exposing (..)
import AudioState exposing (State(..))
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg)
import Playlist.Fetch exposing (fetchPlaylist)
import Ports
import Routing exposing (Route(..))
import Socket.Core exposing (sendMessage)
import Socket.MessageId exposing (MessageId(..))
import Socket.Methods.GetAlbums exposing (getAlbums)
import Socket.Methods.GetArtists exposing (getArtists)
import Socket.Methods.GetPlaylists exposing (getPlaylists)
import Types exposing (Update)
import Url exposing (Url)


logOut : Update Model Msg
logOut model =
    ( { model | username = "", token = Absent }, Ports.websocketClose () )


onUrlChange : Url -> Update Model Msg
onUrlChange url model =
    let
        m =
            { model | route = Routing.parseUrl url }
    in
    case m.route of
        Just (Artist id) ->
            fetchArtist Nothing id m

        Just (Album id) ->
            fetchAlbum Nothing id m

        Just (Playlist id) ->
            fetchPlaylist Nothing id m

        Just Playlists ->
            sendMessage getPlaylists False m

        Just Albums ->
            sendMessage getAlbums False m

        Just Artists ->
            sendMessage getArtists False m

        Nothing ->
            ( m, Cmd.none )
