module Socket.NotificationListener exposing (NotificationListener, makeListener, makeListenerWithParams)

import Json.Decode exposing (Decoder, errorToString)
import Socket.Notification exposing (Notification)
import Types exposing (Update)


type alias NotificationListener model msg =
    Notification -> Update model msg


{-| Make a listener which cares about the notification parameters.
-}
makeListenerWithParams : Decoder a -> (a -> Update model msg) -> (String -> Update model msg) -> NotificationListener model msg
makeListenerWithParams decode update onError notification model =
    case notification.params of
        Just params ->
            case Json.Decode.decodeValue decode params of
                Ok body ->
                    update body model

                -- We got the wrong type of parameters
                Err error ->
                    onError (errorToString error) model

        -- We expected parameters but didn't get any
        Nothing ->
            onError "No parameters received" model


{-| Make a listener which doesn't care about parameters.
Since the the notification is pre-routed based on its method
this means it doesn't depend on the notification at all.
-}
makeListener : Update model msg -> NotificationListener model msg
makeListener =
    always
