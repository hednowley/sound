module Views.PlaylistItem exposing (view)

import Audio.AudioMsg exposing (AudioMsg(..))
import Html.Styled exposing (Html, button, div, text)
import Html.Styled.Attributes exposing (class)
import Html.Styled.Events exposing (onClick)
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg(..))
import Player.Msg exposing (PlayerMsg(..))
import Song.Select exposing (getSong)
import Song.Types exposing (SongId)


view : Model -> Int -> SongId -> Html Msg
view model index songId =
    case getSong model songId of
        Just song ->
            div [ class "playlist__item" ]
                [ button [ onClick <| PlayerMsg (PlayItem index) ] [ text "Play" ]
                , div [] [ text song.name ]
                ]

        Nothing ->
            div [] [ text "Nothing" ]
