port module Ports exposing (..)

import Audio
import Cache exposing (Cache)
import Json.Encode



-- Outgoing ports


port setCache : Cache -> Cmd msg


port loadAudio : Audio.LoadRequest -> Cmd msg


port playAudio : Int -> Cmd msg


port pauseAudio : Int -> Cmd msg


port resumeAudio : Int -> Cmd msg


port setAudioTime : { songId : Int, time : Float } -> Cmd msg


port websocketOut : Json.Encode.Value -> Cmd msg


port websocketOpen : String -> Cmd msg


port websocketClose : () -> Cmd msg



-- Incoming ports


port canPlayAudio : (Int -> msg) -> Sub msg


port audioEnded : (Int -> msg) -> Sub msg


port audioPlaying : ({ songId : Int, time : Float, duration : Maybe Float } -> msg) -> Sub msg


port audioPaused : ({ songId : Int, time : Float, duration : Maybe Float } -> msg) -> Sub msg


port audioTimeChanged : ({ songId : Int, time : Float } -> msg) -> Sub msg


port audioNextPressed : (() -> msg) -> Sub msg


port audioPrevPressed : (() -> msg) -> Sub msg


port websocketOpened : (() -> msg) -> Sub msg


port websocketClosed : (() -> msg) -> Sub msg


port websocketIn : (String -> msg) -> Sub msg
