module Socket.Model exposing (Model)

import Dict exposing (Dict)
import Msg exposing (Msg)
import Socket.Listener exposing (Listener)
import Socket.MessageId exposing (MessageId)
import Socket.NotificationListener exposing (NotificationListener)
import Socket.RequestData exposing (RequestData)


type alias Model m =
    { listeners : Dict Int (Listener m Msg) -- Everything listening out for a server response, keyed by the id of the response they listen for.
    , notificationListeners : Dict String (NotificationListener m Msg) -- Everything listening out for server notifications, keyed by the notification method they listen for.
    , messageQueue : List ( MessageId, RequestData m ) -- Queue for messages which have arrived while the socket is closed
    , nextMessageId : MessageId -- The next unused ID for a message
    , isOpen : Bool -- True iff the socket is open and authenticated
    , ticket : Maybe String
    }
