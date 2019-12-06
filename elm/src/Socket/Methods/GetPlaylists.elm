module Socket.Methods.GetPlaylists exposing (getPlaylists)

import Dict
import Json.Decode exposing (field, int, list, string)
import Model exposing (Model)
import Msg exposing (Msg)
import Socket.Listener exposing (Listener, makeIrresponsibleListener)
import Socket.RequestData exposing (RequestData)
import Types exposing (Update)


type alias Body =
    { playlists : List Playlist }


type alias Playlist =
    { id : Int
    , name : String
    }


getPlaylists : RequestData Model
getPlaylists =
    { method = "getPlaylists"
    , params = Nothing
    , listener = Just onResponse
    }


responseDecoder : Json.Decode.Decoder Body
responseDecoder =
    Json.Decode.map Body
        (field "playlists"
            (list <|
                Json.Decode.map2 Playlist
                    (field "id" int)
                    (field "name" string)
            )
        )


onResponse : Listener Model Msg
onResponse =
    makeIrresponsibleListener
        Nothing
        responseDecoder
        setPlaylists


setPlaylists : Body -> Update Model Msg
setPlaylists body model =
    let
        tuples =
            List.map (\a -> ( a.id, a )) body.playlists
    in
    ( { model | playlists = Dict.fromList tuples }, Cmd.none )
