module Player.Actions exposing
    ( finishShufflePlaylist
    , goNext
    , goPrev
    , onSongEnded
    , pauseCurrent
    , playItem
    , queueAndPlaySong
    , queueSong
    , replacePlaylist
    , resumeCurrent
    , setCurrentTime
    , setShuffle
    , shufflePlaylist
    )

import Array exposing (Array, append, fromList, length, push, slice)
import Audio.Actions exposing (loadSong, playSong)
import AudioState exposing (State(..))
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg)
import Player.Msg exposing (PlayerMsg(..))
import Player.Select exposing (getCurrentSongId, getCurrentSongState, getSongId)
import Ports
import Random
import Random.Array exposing (shuffle)
import Routing exposing (Route(..))
import Song.Types exposing (SongId(..), getRawSongId)
import Types exposing (Update, combine)


setPlaylist : Array SongId -> Model -> Model
setPlaylist playlist model =
    let
        player =
            model.player
    in
    { model
        | player =
            { player
                | playlist = playlist
            }
    }


setPlaying : Maybe Int -> Model -> Model
setPlaying index model =
    let
        player =
            model.player
    in
    { model
        | player =
            { player
                | playing = index
            }
    }


goNext : Update Model Msg
goNext model =
    case model.player.playing of
        Just index ->
            if Array.length model.player.playlist - 1 == index then
                -- No next song
                ( model, Cmd.none )

            else
                case getCurrentSongState model of
                    Just (Playing _) ->
                        playItem (index + 1) model

                    _ ->
                        case getCurrentSongId model of
                            Just songId ->
                                loadSong songId model

                            Nothing ->
                                ( model, Cmd.none )

        Nothing ->
            ( model, Cmd.none )


goPrev : Update Model Msg
goPrev model =
    case model.player.playing of
        Just index ->
            if index == 0 then
                playItem 0 model

            else
                case getCurrentSongState model of
                    Just (Playing { time }) ->
                        if time > 2 then
                            playItem index model

                        else
                            playItem (index - 1) model

                    _ ->
                        case getCurrentSongId model of
                            Just songId ->
                                loadSong songId model

                            Nothing ->
                                ( model, Cmd.none )

        Nothing ->
            ( model, Cmd.none )


setShuffle : Bool -> Update Model Msg
setShuffle on model =
    let
        player =
            model.player

        updated =
            { model
                | player =
                    { player
                        | shuffle = on
                    }
            }
    in
    if on then
        shufflePlaylist updated

    else
        ( updated, Cmd.none )


{-| Shuffles the remaining tracks in the current playlist.
-}
shufflePlaylist : Update Model Msg
shufflePlaylist model =
    case model.player.playing of
        Just index ->
            let
                upcoming =
                    slice (index + 1) (length model.player.playlist) model.player.playlist
            in
            ( model
            , Random.generate
                (Shuffled >> Msg.PlayerMsg)
                (Random.Array.shuffle upcoming)
            )

        Nothing ->
            ( model, Cmd.none )


{-| Uses the random generator output to complete the playlist shuffle.
-}
finishShufflePlaylist : Array SongId -> Model -> Model
finishShufflePlaylist s model =
    case model.player.playing of
        Just index ->
            let
                past =
                    slice 0 (index + 1) model.player.playlist
            in
            setPlaylist (append past s) model

        Nothing ->
            model


pauseCurrent : Update Model Msg
pauseCurrent model =
    case getCurrentSongId model of
        Just songId ->
            ( model, Ports.pauseAudio (getRawSongId songId) )

        Nothing ->
            ( model, Cmd.none )


setCurrentTime : Float -> Update Model Msg
setCurrentTime time model =
    case getCurrentSongId model of
        Just songId ->
            ( model, Ports.setAudioTime { songId = getRawSongId songId, time = time } )

        Nothing ->
            ( model, Cmd.none )


resumeCurrent : Update Model Msg
resumeCurrent model =
    case getCurrentSongId model of
        Just songId ->
            ( model, Ports.resumeAudio (getRawSongId songId) )

        Nothing ->
            ( model, Cmd.none )


playItem : Int -> Update Model Msg
playItem index model =
    case getSongId model index of
        Just songId ->
            model
                |> combine
                    pauseCurrent
                    (\m -> playSong songId (setPlaying (Just index) m))

        Nothing ->
            ( model, Cmd.none )


replacePlaylistWithoutPausing : List SongId -> Update Model Msg
replacePlaylistWithoutPausing playlist model =
    let
        m =
            setPlaylist (fromList playlist) model
    in
    case playlist of
        [] ->
            ( m, Cmd.none )

        first :: _ ->
            playSong first (setPlaying (Just 0) m)


replacePlaylist : List SongId -> Update Model Msg
replacePlaylist playlist =
    combine
        pauseCurrent
        (replacePlaylistWithoutPausing playlist)


queueAndPlaySong : SongId -> Update Model Msg
queueAndPlaySong songId model =
    let
        m =
            queueSong songId model
    in
    playItem (Array.length m.player.playlist - 1) m


queueSong : SongId -> Model -> Model
queueSong songId model =
    setPlaylist
        (push songId model.player.playlist)
        model


onSongEnded : Update Model Msg
onSongEnded model =
    case model.player.playing of
        Just index ->
            if Array.length model.player.playlist - 1 == index then
                -- The last song has finished
                ( setPlaying Nothing model, Cmd.none )

            else
                playItem (index + 1) model

        Nothing ->
            ( model, Cmd.none )
