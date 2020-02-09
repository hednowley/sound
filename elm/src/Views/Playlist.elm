module Views.Playlist exposing (view)

import Audio.Msg exposing (AudioMsg(..))
import Html.Styled exposing (Html, button, div, text)
import Html.Styled.Events exposing (onClick)
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg(..))
import Player.Msg exposing (PlayerMsg(..))
import Playlist.Select exposing (getPlaylist, getPlaylistSongs)
import Playlist.Types exposing (PlaylistId)
import Views.Song


view : PlaylistId -> Model -> Html Msg
view playlistId model =
    case getPlaylist playlistId model of
        Absent ->
            div [] [ text "No playlist" ]

        Loading _ ->
            div [] [ text "Loading playlist" ]

        Loaded playlist ->
            div []
                [ div []
                    [ div [] [ text playlist.name ]
                    , button [ onClick <| PlayerMsg (PlayPlaylist playlistId) ] [ text "Play playlist" ]
                    ]
                , div [] <|
                    List.map (Views.Song.view model) (getPlaylistSongs playlist model)
                ]
