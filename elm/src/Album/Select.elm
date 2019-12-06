module Album.Select exposing (getAlbum, getAlbumArt, getAlbumSongs)

import Album.Types exposing (AlbumId, getRawAlbumId)
import Dict
import Entities.Album exposing (Album)
import Entities.SongSummary exposing (SongSummary)
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Song.Select exposing (getSong)


getAlbum : AlbumId -> Model -> Loadable Album
getAlbum albumId model =
    Dict.get (getRawAlbumId albumId) model.nexus.albums |> Maybe.withDefault Absent


getAlbumSongs : Album -> Model -> List (Maybe SongSummary)
getAlbumSongs album model =
    List.map (getSong model) album.songs



-- |> List.sortBy .track


getAlbumArt : Maybe String -> String
getAlbumArt art =
    case art of
        Nothing ->
            ""

        Just id ->
            "/api/art?size=120&id=" ++ id
