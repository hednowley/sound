module Main exposing (main)

import Audio.Model
import Audio.Msg exposing (AudioMsg(..))
import Audio.Update
import Browser
import Browser.Navigation as Nav exposing (Key)
import Cache exposing (makeCache, makeModel, tryDecode)
import Config exposing (Config)
import Dict
import Html.Styled exposing (div, text, toUnstyled)
import Json.Decode
import Loadable exposing (Loadable(..))
import Model exposing (Model, SocketModelWrap(..), getSocketModel)
import Msg exposing (Msg(..))
import Nexus.Model
import Player.Model
import Player.Msg exposing (PlayerMsg(..))
import Player.Update
import Ports
import Rest.Core as Rest
import Routing exposing (Route(..))
import Socket.Core as Socket
import Socket.Methods.StartScan
import Socket.SocketMsg exposing (SocketMsg(..))
import Socket.Update
import Song.Types exposing (SongId(..))
import Types exposing (Update, combine)
import Updaters exposing (logOut, onUrlChange)
import Url exposing (Url)
import Views.Login
import Views.Root


{-| This is the object passed in by the JS bootloader.
-}
type alias Flags =
    { config : Config
    , model : Maybe Json.Decode.Value
    }


{-| Application entry point.
-}
main : Program Flags Model Msg
main =
    Browser.application
        { init = init
        , view = view
        , update = updateWithStorage
        , subscriptions = subscriptions
        , onUrlRequest = OnUrlRequest
        , onUrlChange = OnUrlChange
        }


{-| A regular update with a "middleware" for storing the model in local storage.
-}
updateWithStorage : Msg -> Update Model Msg
updateWithStorage msg model =
    let
        ( newModel, cmds ) =
            update msg model
    in
    ( newModel
    , Cmd.batch [ Ports.setCache (makeCache newModel), cmds ]
    )


{-| Start the application, passing in the optional serialised model.
-}
init : Flags -> Url -> Key -> ( Model, Cmd Msg )
init flags url navKey =
    let
        model =
            makeModel
                (emptyModel navKey url flags.config)
                (tryDecode flags.model)
    in
    combine (onUrlChange url)
        (\m -> ( m, Socket.Update.reconnect m ))
        model


emptyModel : Key -> Url -> Config -> Model
emptyModel key url config =
    { key = key
    , url = url
    , username = ""
    , password = ""
    , message = ""
    , token = Absent
    , isScanning = False
    , scanCount = 0
    , scanShouldUpdate = False
    , scanShouldDelete = False
    , playlists = Dict.empty
    , nexus = Nexus.Model.empty
    , artists = Dict.empty
    , albums = Dict.empty
    , songs = Dict.empty
    , config = config
    , route = Nothing
    , audio = Audio.Model.emptyModel
    , socket = SocketModelWrap Socket.Update.emptyModel
    , player = Player.Model.emptyModel
    }


{-| Dispatches messages when events are received from ports.
-}
subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.batch
        [ Ports.websocketOpened <| always (SocketMsg SocketOpened)
        , Ports.websocketClosed <| always (SocketMsg SocketClosed)
        , Ports.websocketIn <| SocketIn >> Msg.SocketMsg
        , Ports.canPlayAudio <| SongId >> CanPlay >> AudioMsg
        , Ports.audioEnded <| SongId >> Ended >> AudioMsg
        , Ports.audioPlaying <| Audio.Msg.Playing >> AudioMsg
        , Ports.audioPaused <| Audio.Msg.Paused >> AudioMsg
        , Ports.audioTimeChanged <| Audio.Msg.TimeChanged >> AudioMsg
        , Ports.audioNextPressed <| always (Msg.PlayerMsg Next)
        , Ports.audioPrevPressed <| always (Msg.PlayerMsg Prev)
        ]


update : Msg -> Update Model Msg
update msg model =
    case msg of
        OnUrlRequest request ->
            case request of
                Browser.Internal url ->
                    ( model, Nav.pushUrl model.key (Url.toString url) )

                Browser.External href ->
                    ( model, Nav.load href )

        OnUrlChange url ->
            onUrlChange url model

        UsernameChanged name ->
            ( { model | username = name }, Cmd.none )

        PasswordChanged password ->
            ( { model | password = password }, Cmd.none )

        SubmitLogin ->
            Rest.authenticate model.password { model | password = "" }

        LogOut ->
            logOut model

        GotAuthenticateResponse response ->
            Rest.gotAuthenticateResponse response model

        GotTicketResponse response ->
            Rest.gotTicketResponse response model

        StartScan ->
            Socket.sendMessage
                (Socket.Methods.StartScan.prepareRequest
                    model.scanShouldUpdate
                    model.scanShouldDelete
                )
                False
                model

        ToggleScanUpdate ->
            ( { model | scanShouldUpdate = not model.scanShouldUpdate }, Cmd.none )

        ToggleScanDelete ->
            ( { model | scanShouldDelete = not model.scanShouldDelete }, Cmd.none )

        SocketMsg socketMsg ->
            Socket.Update.update socketMsg model

        AudioMsg a ->
            Audio.Update.update a model

        PlayerMsg p ->
            Player.Update.update p model



-- VIEWS


view : Model -> Browser.Document Msg
view model =
    { title = "Sound"
    , body =
        [ toUnstyled <|
            div []
                [ div [] [ text model.message ]
                , case model.token of
                    Absent ->
                        Views.Login.view model

                    Loadable.Loading _ ->
                        div [] [ text "Getting token..." ]

                    _ ->
                        if (getSocketModel model).isOpen then
                            Views.Root.view model

                        else
                            div [] [ text "Websocket not open" ]
                ]
        ]
    }
