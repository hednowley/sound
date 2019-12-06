module DTO.Credentials exposing (credentialsEncoder)

import Json.Encode


credentialsEncoder : String -> String -> Json.Encode.Value
credentialsEncoder username password =
    Json.Encode.object
        [ ( "username", Json.Encode.string username )
        , ( "password", Json.Encode.string password )
        ]
