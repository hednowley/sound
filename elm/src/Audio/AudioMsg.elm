module Audio.AudioMsg exposing (AudioMsg(..))

import Song.Types exposing (SongId)


type AudioMsg
    = CanPlay SongId -- A song is ready to be played
    | Ended SongId
    | SetTime Float
    | TimeChanged { songId : Int, time : Float }
    | Playing { songId : Int, time : Float, duration : Maybe Float }
    | Paused { songId : Int, time : Float, duration : Maybe Float }
