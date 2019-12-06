module Socket.Listeners.ScanStatus exposing (listener)

import Json.Decode exposing (bool, int)
import Model exposing (Model)
import Msg exposing (Msg)
import Socket.NotificationListener exposing (NotificationListener, makeListenerWithParams)
import Types exposing (Update)


type alias Params =
    { count : Int
    , scanning : Bool
    }


paramsDecoder : Json.Decode.Decoder Params
paramsDecoder =
    Json.Decode.map2
        Params
        (Json.Decode.field "count" int)
        (Json.Decode.field "scanning" bool)


listener : NotificationListener Model Msg
listener =
    makeListenerWithParams paramsDecoder updater onError


onError : String -> Update Model Msg
onError err model =
    ( { model | message = "Couldn't understand scan status: " ++ err }, Cmd.none )


updater : Params -> Update Model Msg
updater params model =
    ( { model | isScanning = params.scanning, scanCount = params.count }, Cmd.none )
