module Album.Types exposing (AlbumId(..), getRawAlbumId)


type AlbumId
    = AlbumId Int


getRawAlbumId : AlbumId -> Int
getRawAlbumId albumId =
    let
        (AlbumId raw) =
            albumId
    in
    raw
