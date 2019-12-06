module Views.Root exposing (view)

import Array
import Audio.AudioMsg exposing (AudioMsg(..))
import Css exposing (..)
import Html.Styled exposing (..)
import Html.Styled.Attributes exposing (class, css)
import Model exposing (Model)
import Msg exposing (Msg(..))
import Routing exposing (Route(..))
import Views.Album
import Views.Albums
import Views.Artist
import Views.Artists
import Views.Home
import Views.Player
import Views.Playlist
import Views.PlaylistItem
import Views.Playlists
import Views.Sidebar


containerStyle : Style
containerStyle =
    Css.batch [ backgroundColor (rgb 240 240 240) ]


view : Model -> Html Msg
view model =
    div
        [ class "app__wrap"
        , css [ containerStyle ]
        ]
        [ div [ class "app__header" ] []
        , div [ class "app__side" ] [ Views.Sidebar.view model ]
        , div [ class "app__main" ]
            [ case model.route of
                Nothing ->
                    Views.Home.view model

                Just (Artist id) ->
                    Views.Artist.view id model

                Just (Album id) ->
                    Views.Album.view id model

                Just (Playlist id) ->
                    Views.Playlist.view id model

                Just Artists ->
                    Views.Artists.view model

                Just Albums ->
                    Views.Albums.view model

                Just Playlists ->
                    Views.Playlists.view model
            ]
        , div [ class "app__playlist" ]
            [ div [] [ text "Playlist" ]
            , div [ class "playlist__items" ]
                (Array.toList <|
                    Array.indexedMap
                        (Views.PlaylistItem.view model)
                        model.player.playlist
                )
            ]
        , div [ class "app__player" ] [ Views.Player.view model ]
        ]
