module Player.Select exposing (getCurrentSongId, getCurrentSongState, getSongId, isPlaying, shuffleIsOn)

import Array
import Audio.Select exposing (getSongState)
import Audio.State exposing (State(..))
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Routing exposing (Route(..))
import Song.Types exposing (SongId(..))


{-| Gets the ID of the song at the given position in the playlist.
-}
getSongId : Model -> Int -> Maybe SongId
getSongId model index =
    Array.get index model.player.playlist


{-| Gets the ID of the currently playing song.
-}
getCurrentSongId : Model -> Maybe SongId
getCurrentSongId model =
    model.player.playing
        |> Maybe.andThen (getSongId model)


getCurrentSongState : Model -> Maybe State
getCurrentSongState model =
    getCurrentSongId model |> Maybe.andThen (getSongState model)


shuffleIsOn : Model -> Bool
shuffleIsOn model =
    model.player.shuffle


isPlaying : Model -> Bool
isPlaying model =
    case getCurrentSongState model of
        Just (Playing _) ->
            True

        _ ->
            False
