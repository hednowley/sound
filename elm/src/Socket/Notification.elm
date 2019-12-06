module Socket.Notification exposing (Notification, decode)

import Json.Decode exposing (Decoder, field, map2, maybe, string, value)


{-| A websocket message which doesn't expect a reply.
-}
type alias Notification =
    { method : String
    , params : Maybe Json.Decode.Value
    }


decode : Decoder Notification
decode =
    map2 Notification
        (field "method" string)
        (maybe <| field "params" value)
