module DTO.Authenticate exposing (Response, decode)

import Json.Decode exposing (Decoder, andThen, fail, field, map, string, succeed)


{-| A potential error message. Absense means success.
-}
type alias Response =
    Maybe String


{-| Decode the message.
-}
decode : Decoder Response
decode =
    field "status" string
        |> andThen decodeData


{-| Decode the data portion of the message.
-}
decodeData : String -> Decoder Response
decodeData status =
    case status of
        "success" ->
            succeed Nothing

        "error" ->
            map Just (field "data" string)

        _ ->
            fail "Unknown status"
