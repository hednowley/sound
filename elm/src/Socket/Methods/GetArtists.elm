module Socket.Methods.GetArtists exposing (getArtists)

import Dict
import Json.Decode exposing (field, list)
import Model exposing (Model)
import Msg exposing (Msg)
import Socket.DTO.ArtistSummary exposing (convert, decode)
import Socket.Listener exposing (Listener, makeIrresponsibleListener)
import Socket.RequestData exposing (RequestData)
import Types exposing (Update)


type alias Body =
    { artists : List Artist }


type alias Artist =
    { id : Int
    , name : String
    }


getArtists : RequestData Model
getArtists =
    { method = "getArtists"
    , params = Nothing
    , listener = Just onResponse
    }


responseDecoder : Json.Decode.Decoder Body
responseDecoder =
    Json.Decode.map Body
        (field "artists"
            (list decode)
        )


onResponse : Listener Model Msg
onResponse =
    makeIrresponsibleListener
        Nothing
        responseDecoder
        setArtists


setArtists : Body -> Update Model Msg
setArtists body model =
    let
        tuples =
            List.map (\a -> ( a.id, convert a )) body.artists
    in
    ( { model | artists = Dict.fromList tuples }, Cmd.none )
