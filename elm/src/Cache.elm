module Cache exposing (Cache, makeCache, makeModel, tryDecode)

import Json.Decode as Decode exposing (Decoder)
import Loadable exposing (fromMaybe, toMaybe)
import Model exposing (Model)


{-| A version of the model which can be stored in the browser.
-}
type alias Cache =
    { token : Maybe String }


{-| Create a cache from a model.
-}
makeCache : Model -> Cache
makeCache model =
    { token = toMaybe model.token }


{-| Create a model from a cache.
-}
makeModel : Model -> Maybe Cache -> Model
makeModel default cache =
    case cache of
        Just c ->
            { default | token = fromMaybe c.token }

        Nothing ->
            default


{-| Try to decode a cache from an optional JSON value.
-}
tryDecode : Maybe Decode.Value -> Maybe Cache
tryDecode =
    Maybe.andThen <| Decode.decodeValue decode >> Result.toMaybe


decode : Decoder Cache
decode =
    Decode.map Cache
        (Decode.field "token" (Decode.maybe Decode.string))
