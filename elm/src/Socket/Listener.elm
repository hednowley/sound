module Socket.Listener exposing (Listener, combineListeners, makeIrresponsibleListener, makeResponsibleListener)

import Json.Decode exposing (Decoder, Value, decodeValue)
import Socket.Response exposing (Response)
import Types exposing (Update, combine, noOp)


{-| Describes how to transform the model and dispatch commands with an incoming websocket message.
-}
type alias Listener model msg =
    Response -> Update model msg


combineListeners : Listener model msg -> Listener model msg -> Listener model msg
combineListeners first second response =
    let
        u1 =
            first response

        u2 =
            second response
    in
    combine u1 u2


{-| Make a new listener which has error handling.
-}
makeResponsibleListener :
    Maybe (Response -> model -> model)
    -> Decoder a
    -> (a -> Update model msg)
    -> Decoder b
    -> (b -> Update model msg)
    -> Listener model msg
makeResponsibleListener maybeCleanup successDecoder onSuccess errorDecoder onError response model =
    let
        cleaned =
            case maybeCleanup of
                Just cleanup ->
                    cleanup response model

                Nothing ->
                    model
    in
    case response.body of
        Ok success ->
            processSuccess success successDecoder onSuccess cleaned

        Err error ->
            processError error errorDecoder onError cleaned


{-| Make a new listener which has no error handling.
-}
makeIrresponsibleListener :
    Maybe (Response -> model -> model)
    -> Decoder a
    -> (a -> Update model msg)
    -> Listener model msg
makeIrresponsibleListener maybeCleanup successDecoder onSuccess response model =
    let
        cleaned =
            case maybeCleanup of
                Just cleanup ->
                    cleanup response model

                Nothing ->
                    model
    in
    case response.body of
        Ok success ->
            processSuccess success successDecoder onSuccess cleaned

        Err _ ->
            ( cleaned, Cmd.none )


processSuccess : Value -> Decoder a -> (a -> Update model msg) -> Update model msg
processSuccess json decoder makeUpdate model =
    case decodeValue decoder json of
        Ok body ->
            makeUpdate body model

        Err _ ->
            ( model, Cmd.none )


processError : Value -> Decoder a -> (a -> Update model msg) -> Update model msg
processError json decoder _ =
    case decodeValue decoder json of
        Ok _ ->
            noOp

        Err _ ->
            noOp
