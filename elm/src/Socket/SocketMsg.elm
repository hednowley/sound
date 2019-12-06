module Socket.SocketMsg exposing (SocketMsg(..))


type SocketMsg
    = SocketOpened -- The websocket has been successfully opened
    | SocketClosed -- The websocket has been closed
    | SocketIn String -- A message has been received over the websocket
