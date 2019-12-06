module Loadable exposing (Loadable(..), fromMaybe, toMaybe)

import Socket.MessageId exposing (MessageId)


{-| Something which takes time to load.
-}
type Loadable value
    = Absent
    | Loading MessageId
    | Loaded value


{-| Convert a Loadable to a Maybe.
-}
toMaybe : Loadable value -> Maybe value
toMaybe loadable =
    case loadable of
        Loaded v ->
            Just v

        _ ->
            Nothing


{-| Convert a Maybe into a Loadable.
-}
fromMaybe : Maybe value -> Loadable value
fromMaybe maybe =
    case maybe of
        Just v ->
            Loaded v

        _ ->
            Absent
