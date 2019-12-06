module Socket.Actions exposing (addListener, addListenerExternal, removeListener)

import Dict
import Model exposing (getSocketModel, setSocketModel)
import Msg exposing (Msg)
import Socket.Listener exposing (Listener, combineListeners)
import Socket.MessageId exposing (MessageId, getRawMessageId)
import Socket.Model exposing (Model)
import Socket.Select exposing (getListener)


type alias Model =
    Socket.Model.Model Model.Model


addListenerExternal : MessageId -> Listener Model.Model Msg -> Model.Model -> Model.Model
addListenerExternal id listener model =
    let
        socket =
            getSocketModel model
    in
    addListener id listener socket |> setSocketModel model


{-| Store a new Websocket listener in the model.
-}
addListener : MessageId -> Listener Model.Model Msg -> Model -> Model
addListener id listener model =
    case getListener id model of
        Just existing ->
            let
                combined =
                    combineListeners existing listener
            in
            insertListener id combined model

        Nothing ->
            insertListener id listener model


insertListener : MessageId -> Listener Model.Model Msg -> Model -> Model
insertListener messageId listener model =
    { model
        | listeners =
            Dict.insert (getRawMessageId messageId) listener model.listeners
    }


{-| Remove a stored Websocket listener from the model.
-}
removeListener : Int -> Model -> Model
removeListener id model =
    { model
        | listeners =
            Dict.remove id model.listeners
    }
