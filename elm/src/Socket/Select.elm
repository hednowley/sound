module Socket.Select exposing (getListener, getNotificationListener)

import Dict
import Model
import Msg exposing (Msg)
import Socket.Listener exposing (Listener)
import Socket.MessageId exposing (MessageId, getRawMessageId)
import Socket.Model
import Socket.NotificationListener exposing (NotificationListener)


type alias Model =
    Socket.Model.Model Model.Model


{-| Try and retrieve the listener with the given ID.
-}
getListener : MessageId -> Model -> Maybe (Listener Model.Model Msg)
getListener id model =
    Dict.get (getRawMessageId id) model.listeners


{-| Try and retrieve the notification listener for the given method.
-}
getNotificationListener : String -> Model -> Maybe (NotificationListener Model.Model Msg)
getNotificationListener method model =
    Dict.get method model.notificationListeners
