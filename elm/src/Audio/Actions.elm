module Audio.Actions exposing
    ( loadSong
    , onSongLoaded
    , onTimeChanged
    , playSong
    , updateSongState
    )

import Audio.Request exposing (makeLoadRequest)
import Audio.Select exposing (getSongState)
import Audio.State exposing (State(..))
import Dict
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg)
import Player.Select exposing (getSongId)
import Ports
import Routing exposing (Route(..))
import Song.Types exposing (SongId(..), getRawSongId)
import Types exposing (Update)


updateSongState : SongId -> Audio.State.State -> Model -> Model
updateSongState songId state model =
    let
        (SongId s) =
            songId

        audio =
            model.audio
    in
    { model | audio = { audio | songs = Dict.insert s state audio.songs } }


onSongLoaded : SongId -> Update Model Msg
onSongLoaded songId model =
    let
        m =
            updateSongState
                songId
                (Audio.State.Loaded { duration = Nothing })
                model

        playing =
            Maybe.andThen (getSongId model) model.player.playing
    in
    if playing == Just songId then
        playSong songId m

    else
        ( m, Cmd.none )


onTimeChanged : SongId -> Float -> Model -> Model
onTimeChanged songId time model =
    case getSongState model songId of
        Just (Playing p) ->
            updateSongState songId (Playing { p | time = time }) model

        _ ->
            model


playSong : SongId -> Update Model Msg
playSong songId model =
    case getSongState model songId of
        Just Audio.State.Loading ->
            ( model, Cmd.none )

        Just (Audio.State.Loaded _) ->
            ( model, Ports.playAudio (getRawSongId songId) )

        Just (Audio.State.Playing _) ->
            ( model, Ports.playAudio (getRawSongId songId) )

        _ ->
            loadSong songId model


loadSong : SongId -> Update Model Msg
loadSong songId model =
    let
        m =
            updateSongState songId Audio.State.Loading model
    in
    ( m, Ports.loadAudio <| makeLoadRequest songId )
