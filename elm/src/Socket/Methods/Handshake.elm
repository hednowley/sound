module Socket.Methods.Handshake exposing (makeRequest, prepareRequest)

import Json.Decode
import Json.Encode
import Model exposing (Model)
import Msg exposing (Msg)
import Socket.Listener exposing (Listener, makeIrresponsibleListener)
import Socket.RequestData exposing (RequestData)
import Types exposing (Update)


type alias Body =
    { accepted : Bool }


responseDecoder : Json.Decode.Decoder Body
responseDecoder =
    Json.Decode.map Body
        (Json.Decode.field "accepted" Json.Decode.bool)


makeRequest : String -> Json.Encode.Value
makeRequest ticket =
    Json.Encode.object
        [ ( "ticket", Json.Encode.string ticket ) ]


{-| Make a message which starts the websocket handshake.
-}
prepareRequest : String -> Update Model Msg -> RequestData Model
prepareRequest ticket onHandshakeSuccess =
    { method = "handshake"
    , params = makeRequest ticket |> Just
    , listener = onResponse onHandshakeSuccess |> Just
    }


onResponse : Update Model Msg -> Listener Model Msg
onResponse onHandshakeSuccess =
    makeIrresponsibleListener
        Nothing
        responseDecoder
        (onSuccess onHandshakeSuccess)


onSuccess : Update Model Msg -> Body -> Update Model Msg
onSuccess onHandshakeSuccess response model =
    if response.accepted then
        onHandshakeSuccess { model | message = "Websocket handshake succeeded" }

    else
        ( { model | message = "Websocket handshake failed" }, Cmd.none )
