module Views.Artists exposing (view)

import Artist.Types exposing (getRawArtistId)
import Dict
import Entities.ArtistSummary exposing (ArtistSummaries)
import Html.Styled exposing (Html, a, div, text)
import Html.Styled.Attributes exposing (class, href)
import Model exposing (Model)
import Msg exposing (Msg(..))
import String exposing (fromInt)


view : Model -> Html Msg
view model =
    div [ class "home__wrap" ]
        [ viewAlbums model.artists
        ]


viewAlbums : ArtistSummaries -> Html msg
viewAlbums albums =
    div [ class "home__artists" ]
        (List.map
            (\album ->
                div
                    [ class "home__artist" ]
                    [ a [ href <| "/artist/" ++ fromInt (getRawArtistId album.id) ] [ text album.name ] ]
            )
            (Dict.values albums)
        )
