module Album.Fetch exposing (fetchAlbum)

import Album.Types exposing (AlbumId, getRawAlbumId)
import Entities.Album exposing (Album)
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg)
import Nexus.Fetch exposing (fetch)
import Socket.DTO.Album exposing (convert, decode)
import Socket.DTO.SongSummary exposing (convertMany)
import Song.Types exposing (SongId(..), getRawSongId)
import Types exposing (Update)
import Util exposing (insertMany)


fetchAlbum : Maybe (Album -> Update Model Msg) -> AlbumId -> Update Model Msg
fetchAlbum maybeCallback =
    fetch
        getRawAlbumId
        "getAlbum"
        decode
        saveSongs
        convert
        { get = .albums
        , set = \repo -> \m -> { m | albums = repo }
        }
        maybeCallback


saveSongs : Socket.DTO.Album.Album -> Model -> Model
saveSongs album model =
    let
        songs =
            convertMany album.songs
    in
    { model
        | songs =
            insertMany
                (.id >> getRawSongId)
                identity
                songs
                model.songs
    }
