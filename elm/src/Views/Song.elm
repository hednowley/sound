module Views.Song exposing (view)

import Audio.Msg exposing (AudioMsg(..))
import Entities.SongSummary exposing (SongSummary)
import Html.Styled exposing (Html, button, div, text)
import Html.Styled.Attributes exposing (class)
import Html.Styled.Events exposing (onClick)
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg(..))
import Player.Msg exposing (PlayerMsg(..))
import String exposing (fromInt)


view : Model -> Maybe SongSummary -> Html Msg
view model maybeSong =
    case maybeSong of
        Just song ->
            div [ class "album__song" ]
                [ button [ onClick <| PlayerMsg (Play song.id) ] [ text "Play" ]
                , button [ onClick <| PlayerMsg (Queue song.id) ] [ text "Queue" ]
                , div [] [ text <| fromInt song.track ]
                , div [] [ text song.name ]
                ]

        Nothing ->
            div [ class "album__song" ] []
