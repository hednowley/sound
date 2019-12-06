module Views.MiniAlbum exposing (view)

import Album.Select exposing (getAlbumArt)
import Album.Types exposing (getRawAlbumId)
import Audio.AudioMsg exposing (AudioMsg(..))
import Entities.AlbumSummary exposing (AlbumSummary)
import Html.Styled exposing (Html, a, button, div, img, text)
import Html.Styled.Attributes exposing (class, href, src)
import Html.Styled.Events exposing (onClick)
import Loadable exposing (Loadable(..))
import Msg exposing (Msg(..))
import Player.Msg exposing (PlayerMsg(..))
import String exposing (fromInt)


view : AlbumSummary -> Html Msg
view album =
    div [ class "home__artist" ]
        [ div []
            [ img [ class "artist__album--art", src <| getAlbumArt album.artId ] []
            , a [ href <| "/album/" ++ fromInt (getRawAlbumId album.id) ] [ text album.name ]
            , button [ onClick <| PlayerMsg (PlayAlbum album.id) ] [ text "Play" ]
            ]
        ]
