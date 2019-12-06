module Artist.Types exposing (ArtistId(..), getRawArtistId)


type ArtistId
    = ArtistId Int


getRawArtistId : ArtistId -> Int
getRawArtistId albumId =
    let
        (ArtistId raw) =
            albumId
    in
    raw
