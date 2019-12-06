module Socket.Update exposing (emptyModel, reconnect, update)

import Dict
import Loadable exposing (Loadable(..))
import Model exposing (getSocketModel, setSocketModel)
import Msg exposing (Msg(..))
import Rest.Core as Rest
import Socket.Core exposing (messageIn, sendMessage)
import Socket.Listeners.ScanStatus
import Socket.MessageId exposing (MessageId(..))
import Socket.Methods.Handshake
import Socket.Methods.Start
import Socket.Model
import Socket.SocketMsg exposing (SocketMsg(..))
import Types exposing (Update)


emptyModel : Socket.Model.Model Model.Model
emptyModel =
    { listeners = Dict.empty
    , notificationListeners =
        Dict.fromList
            [ ( "scanStatus", Socket.Listeners.ScanStatus.listener )
            ]
    , messageQueue = []
    , nextMessageId = MessageId 1
    , isOpen = False
    , ticket = Nothing
    }


update : SocketMsg -> Update Model.Model Msg
update msg model =
    let
        socket =
            getSocketModel model
    in
    case msg of
        SocketOpened ->
            case socket.ticket of
                Just ticket ->
                    -- Start the ticket handshake now that websocket is open
                    negotiateSocket ticket model

                Nothing ->
                    ( { model | message = "Can't negotiate websocket as there is no ticket" }, Cmd.none )

        SocketClosed ->
            -- Try to reopen the websocket
            ( setSocketModel model { socket | isOpen = False }, reconnect model )

        SocketIn message ->
            messageIn message model


negotiateSocket : String -> Update Model.Model Msg
negotiateSocket ticket model =
    -- Force send since the socket is unauthenticated!
    sendMessage
        (Socket.Methods.Handshake.prepareRequest ticket Socket.Methods.Start.start)
        True
        model


{-| Tries to connect to the websocket if there is cached token.
-}
reconnect : Model.Model -> Cmd Msg
reconnect model =
    case model.token of
        Loadable.Loaded _ ->
            Rest.getTicket model

        _ ->
            Cmd.none
