module Socket.Methods.GetAlbums exposing (getAlbums)

import Dict
import Json.Decode exposing (field, list)
import Model exposing (Model)
import Msg exposing (Msg)
import Socket.DTO.AlbumSummary exposing (AlbumSummary, convert, decode)
import Socket.Listener exposing (Listener, makeIrresponsibleListener)
import Socket.RequestData exposing (RequestData)
import Types exposing (Update)


type alias Body =
    { albums : List AlbumSummary }


getAlbums : RequestData Model
getAlbums =
    { method = "getAlbums"
    , params = Nothing
    , listener = Just onResponse
    }


responseDecoder : Json.Decode.Decoder Body
responseDecoder =
    Json.Decode.map Body
        (field "albums"
            (list decode)
        )


onResponse : Listener Model Msg
onResponse =
    makeIrresponsibleListener
        Nothing
        responseDecoder
        setAlbums


setAlbums : Body -> Update Model Msg
setAlbums body model =
    let
        tuples =
            List.map (\a -> ( a.id, convert a )) body.albums
    in
    ( { model | albums = Dict.fromList tuples }, Cmd.none )
