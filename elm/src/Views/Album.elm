module Views.Album exposing (view)

import Album.Select exposing (getAlbum, getAlbumArt, getAlbumSongs)
import Album.Types exposing (AlbumId)
import Audio.Msg exposing (AudioMsg(..))
import Html.Styled exposing (Html, button, div, img, text)
import Html.Styled.Attributes exposing (class, src)
import Html.Styled.Events exposing (onClick)
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg(..))
import Player.Msg exposing (PlayerMsg(..))
import Views.Song


view : AlbumId -> Model -> Html Msg
view id model =
    case getAlbum id model of
        Absent ->
            div [] [ text "No album" ]

        Loading _ ->
            div [] [ text "Loading album" ]

        Loaded album ->
            div []
                [ div []
                    [ div [] [ text album.name ]
                    , img [ class "album__art", src <| getAlbumArt album.artId ] []
                    , button [ onClick <| PlayerMsg (PlayAlbum id) ] [ text "Play album" ]
                    ]
                , div [] <|
                    List.map (Views.Song.view model) (getAlbumSongs album model)
                ]
