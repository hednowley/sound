module Socket.Response exposing (Response, decode)

import Json.Decode exposing (Decoder, Value, andThen, field, int, map, oneOf, value)
import Socket.MessageId exposing (MessageId(..))


{-| A reply received through a websocket.
The body represents either the success JSON or error JSON of a response.
-}
type alias Response =
    { id : MessageId
    , body : Result Value Value
    }


decode : Decoder Response
decode =
    field "id" int
        |> andThen decodeInner


decodeInner : Int -> Decoder Response
decodeInner id =
    let
        make =
            Response (MessageId id)
    in
    oneOf
        [ map (Err >> make) (field "error" value)
        , map (Ok >> make) (field "result" value)
        ]
