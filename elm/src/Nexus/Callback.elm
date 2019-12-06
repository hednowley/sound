module Nexus.Callback exposing (Callback, resolve)

import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg)
import Types exposing (Update)


type alias Callback a =
    a -> Update Model Msg


resolve : Maybe (Callback a) -> Callback a
resolve maybeCallback =
    Maybe.withDefault
        (\a -> \m -> ( m, Cmd.none ))
        maybeCallback
