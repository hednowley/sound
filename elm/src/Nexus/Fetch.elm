module Nexus.Fetch exposing (fetch)

import Dict exposing (Dict)
import Json.Decode exposing (Decoder)
import Json.Encode
import Loadable exposing (Loadable(..))
import Model exposing (Model)
import Msg exposing (Msg)
import Nexus.Callback exposing (Callback, resolve)
import Nexus.Model
import Socket.Actions exposing (addListenerExternal)
import Socket.Core exposing (sendMessageWithId)
import Socket.Listener exposing (Listener, makeIrresponsibleListener)
import Socket.RequestData exposing (RequestData)
import Types exposing (Update)


{-| A place to store things.
-}
type alias Repo obj =
    Dict Int (Loadable obj)


{-| Tells you how to get and set a repo.
-}
type alias RepoAccessor obj =
    { get : Nexus.Model.Model -> Repo obj
    , set : Repo obj -> Nexus.Model.Model -> Nexus.Model.Model
    }


type alias SaveObject obj =
    Loadable obj -> Model -> Model


type alias OnDecode obj =
    obj -> Model -> Model


noDecode : OnDecode obj
noDecode _ model =
    model


fetch :
    (id -> Int)
    -> String
    -> Decoder dto
    -> OnDecode dto
    -> (dto -> obj)
    -> RepoAccessor obj
    -> Maybe (Callback obj)
    -> id
    -> Update Model Msg
fetch extractId method decoder onDecode convert accessRepo maybeCallback id model =
    let
        extractedId =
            extractId id

        stored =
            Dict.get extractedId (accessRepo.get model.nexus)
    in
    case Maybe.withDefault Absent stored of
        {- This thing has never been fetched. -}
        Absent ->
            let
                listener =
                    makeListener
                        accessRepo
                        decoder
                        onDecode
                        convert
                        (resolve maybeCallback)
                        extractedId

                ( ( sentModel, cmd ), messageId ) =
                    sendMessageWithId
                        (makeRequest extractedId method listener)
                        False
                        model

                saved =
                    saveObject accessRepo extractedId (Loading messageId) sentModel
            in
            ( saved, cmd )

        {- There is already an in-flight fetch for this thing. -}
        Loading requestId ->
            case maybeCallback of
                Just callback ->
                    let
                        -- Don't include onDecode as we can assume it's in the existing listener
                        listener =
                            makeListener
                                accessRepo
                                decoder
                                noDecode
                                convert
                                callback
                                extractedId
                    in
                    ( addListenerExternal requestId listener model, Cmd.none )

                Nothing ->
                    ( model, Cmd.none )

        {- This thing has already been fetched -}
        Loaded playlist ->
            case maybeCallback of
                Just callback ->
                    callback playlist model

                Nothing ->
                    ( model, Cmd.none )


{-| Makes a socket request.
-}
makeRequest : Int -> String -> Listener Model Msg -> RequestData Model
makeRequest id method listener =
    { method = method
    , params = Just <| Json.Encode.object [ ( "id", Json.Encode.int id ) ]
    , listener = Just listener
    }


saveObject : RepoAccessor obj -> Int -> SaveObject obj
saveObject accessRepo id loadable model =
    let
        repo =
            Dict.insert id loadable (accessRepo.get model.nexus)
    in
    { model | nexus = accessRepo.set repo model.nexus }


{-| Makes a socket listener which
-}
makeListener :
    RepoAccessor obj
    -> Decoder dto
    -> OnDecode dto
    -> (dto -> obj)
    -> Callback obj
    -> Int
    -> Listener Model Msg
makeListener accessRepo decoder onDecode convert callback id =
    makeIrresponsibleListener
        Nothing
        decoder
        (onSuccess
            onDecode
            convert
            (saveObject accessRepo id)
            callback
        )



-- TODO : FIX THIS


onSuccess :
    OnDecode dto
    -> (dto -> obj)
    -> SaveObject obj
    -> Callback obj
    -> dto
    -> Update Model Msg
onSuccess onDecode convert save callback dto model =
    let
        m1 =
            onDecode dto model

        thing =
            convert dto

        saved =
            save (Loaded thing) m1
    in
    callback thing saved
