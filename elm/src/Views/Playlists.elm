module Views.Playlists exposing (view)

import Dict
import Entities.PlaylistSummary exposing (PlaylistSummaries)
import Html.Styled exposing (Html, a, div, text)
import Html.Styled.Attributes exposing (class, href)
import Model exposing (Model)
import Msg exposing (Msg(..))
import String exposing (fromInt)


view : Model -> Html Msg
view model =
    div [ class "home__wrap" ]
        [ viewPlaylists model.playlists
        ]


viewPlaylists : PlaylistSummaries -> Html msg
viewPlaylists playlists =
    div [ class "home__artists" ]
        (List.map
            (\playlist -> div [ class "home__artist" ] [ a [ href <| "/playlist/" ++ fromInt playlist.id ] [ text playlist.name ] ])
            (Dict.values playlists)
        )
