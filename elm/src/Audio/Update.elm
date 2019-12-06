module Audio.Update exposing (update)

import Audio.Actions exposing (onSongLoaded, onTimeChanged, updateSongState)
import Audio.AudioMsg exposing (AudioMsg(..))
import AudioState
import Model exposing (Model)
import Msg exposing (Msg)
import Player.Actions
    exposing
        ( onSongEnded
        , setCurrentTime
        )
import Song.Types exposing (SongId(..))
import Types exposing (Update)


update : AudioMsg -> Update Model Msg
update msg model =
    case msg of
        CanPlay songId ->
            onSongLoaded songId model

        Ended _ ->
            onSongEnded model

        SetTime time ->
            setCurrentTime time model

        Playing { songId, time, duration } ->
            ( updateSongState
                (SongId songId)
                (AudioState.Playing { paused = False, time = time, duration = duration })
                model
            , Cmd.none
            )

        Paused { songId, time, duration } ->
            ( updateSongState
                (SongId songId)
                (AudioState.Playing { paused = True, time = time, duration = duration })
                model
            , Cmd.none
            )

        TimeChanged args ->
            ( onTimeChanged
                (SongId args.songId)
                args.time
                model
            , Cmd.none
            )
