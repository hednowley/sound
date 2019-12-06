module Song.Types exposing (SongId(..), getRawSongId)


type SongId
    = SongId Int


getRawSongId : SongId -> Int
getRawSongId songId =
    let
        (SongId raw) =
            songId
    in
    raw
