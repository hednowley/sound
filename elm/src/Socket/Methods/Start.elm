module Socket.Methods.Start exposing (start)

import Model exposing (Model, getSocketModel, setSocketModel)
import Msg exposing (Msg)
import Socket.Core exposing (sendMessage, sendQueuedMessage)
import Socket.Methods.GetArtists exposing (getArtists)
import Types exposing (Update, combineMany)


{-| This should be run once the websocket handshake is complete.
-}
start : Update Model Msg
start =
    combineMany
        [ setWebsocketOpen
        , processQueue
        , sendMessage getArtists False
        ]


setWebsocketOpen : Update Model Msg
setWebsocketOpen model =
    let
        socket =
            getSocketModel model

        updated =
            setSocketModel model { socket | isOpen = True }
    in
    ( updated, Cmd.none )


processQueue : Update Model Msg
processQueue model =
    combineMany
        (List.map sendQueuedMessage (getSocketModel model).messageQueue)
        model
