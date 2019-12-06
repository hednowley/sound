module Rest.Core exposing (authenticate, getTicket, gotAuthenticateResponse, gotTicketResponse)

import DTO.Authenticate
import DTO.Credentials
import DTO.Ticket
import Http exposing (Error(..))
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg(..))
import Socket.Core
import Socket.MessageId exposing (MessageId(..))
import String exposing (fromInt)
import Types exposing (Update)
import Updaters exposing (logOut)


{-| Post credentials to the server.
-}
authenticate : String -> Update Model Msg
authenticate password model =
    let
        credentials =
            DTO.Credentials.credentialsEncoder model.username password
    in
    ( { model | token = Loading <| MessageId 1337 }
    , Http.riskyRequest
        { method = "POST"
        , headers = []
        , timeout = Nothing
        , tracker = Nothing
        , body = Http.jsonBody credentials
        , url = "/api/authenticate"
        , expect = Http.expectJson GotAuthenticateResponse DTO.Authenticate.decode
        }
    )


{-| Parse a response from the server to credentials. If it's worked then the response will be a JWT.
-}
gotAuthenticateResponse : Result Http.Error DTO.Authenticate.Response -> Update Model Msg
gotAuthenticateResponse response model =
    case response of
        Ok result ->
            case result of
                -- We got a token
                Nothing ->
                    ( { model | token = Loaded "" }, getTicket model )

                -- Server has told us why we can't have a token
                Just err ->
                    ( { model | message = err, token = Absent }, Cmd.none )

        -- We don't understand what the server said
        Err e ->
            let
                message =
                    case e of
                        BadStatus s ->
                            "BadStatus: " ++ fromInt s

                        Timeout ->
                            "Timeout"

                        NetworkError ->
                            "NetworkError"

                        BadUrl _ ->
                            "BadUrl"

                        BadBody bb ->
                            "BadBody: " ++ bb
            in
            ( { model | message = message, token = Absent }, Cmd.none )


{-| Ask the server for a websocket ticket, using our JWT.
-}
getTicket : Model -> Cmd Msg
getTicket model =
    -- "Risky" request as we need to send JWT cookie
    Http.riskyRequest
        { method = "GET"
        , headers = []
        , body = Http.emptyBody
        , timeout = Nothing
        , tracker = Nothing
        , url = "/api/ticket"
        , expect = Http.expectJson GotTicketResponse DTO.Ticket.decode
        }


{-| The server replied to a request for a websocket ticket.
-}
gotTicketResponse : Result Http.Error String -> Update Model Msg
gotTicketResponse response model =
    case response of
        Ok ticket ->
            Socket.Core.open ticket model

        Err _ ->
            -- Log out as something must be wrong with our session
            logOut { model | message = "Could not retrieve websocket ticket" }
