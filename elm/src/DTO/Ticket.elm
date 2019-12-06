module DTO.Ticket exposing (decode)

import Json.Decode exposing (Decoder, field, string)


decode : Decoder String
decode =
    field "data" (field "ticket" string)
