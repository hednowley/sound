module Msg exposing (Msg(..))

import Audio.Msg exposing (AudioMsg)
import Browser
import DTO.Authenticate
import Http
import Player.Msg exposing (PlayerMsg)
import Socket.SocketMsg exposing (SocketMsg)
import Url


type Msg
    = OnUrlChange Url.Url
    | OnUrlRequest Browser.UrlRequest
    | UsernameChanged String
    | PasswordChanged String
    | SubmitLogin
    | LogOut
    | ToggleScanUpdate
    | ToggleScanDelete
    | StartScan -- Ask for a scan to be started
    | GotAuthenticateResponse (Result Http.Error DTO.Authenticate.Response) -- Server has replied to posting of credentials
    | GotTicketResponse (Result Http.Error String) -- Server has replied to a request for a websocket ticket
    | SocketMsg SocketMsg
    | PlayerMsg PlayerMsg
    | AudioMsg AudioMsg
