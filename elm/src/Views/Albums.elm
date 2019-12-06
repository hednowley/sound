module Views.Albums exposing (view)

import Album.Types exposing (getRawAlbumId)
import Dict
import Entities.AlbumSummary exposing (AlbumSummaries)
import Html.Styled exposing (Html, a, div, text)
import Html.Styled.Attributes exposing (class, href)
import Model exposing (Model)
import Msg exposing (Msg(..))
import String exposing (fromInt)


view : Model -> Html Msg
view model =
    div [ class "home__wrap" ]
        [ viewAlbums model.albums
        ]


viewAlbums : AlbumSummaries -> Html msg
viewAlbums albums =
    div [ class "home__artists" ]
        (List.map
            (\album ->
                div
                    [ class "home__artist" ]
                    [ a [ href <| "/album/" ++ fromInt (getRawAlbumId album.id) ] [ text album.name ] ]
            )
            (Dict.values albums)
        )
