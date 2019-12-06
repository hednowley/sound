module Socket.RequestData exposing (RequestData)

import Json.Encode
import Msg exposing (Msg)
import Socket.Listener


{-| Describes a message to send down the websocket and optionally how to handle a response to that message.
-}
type alias RequestData model =
    { method : String
    , params : Maybe Json.Encode.Value
    , listener : Maybe (Socket.Listener.Listener model Msg) -- How any replies to the message should be handled.
    }
