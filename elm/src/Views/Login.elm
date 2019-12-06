module Views.Login exposing (view)

import Html.Styled exposing (Html, button, div, form, input, text)
import Html.Styled.Attributes exposing (class, disabled, name, placeholder, type_, value)
import Html.Styled.Events exposing (onInput)
import Json.Decode
import Model exposing (Model)
import Msg exposing (Msg(..))


view : Model -> Html Msg
view model =
    div [ class "login__wrap" ]
        [ form [ class "login__container" ]
            [ div [ class "login__logo " ] [ text "Sound." ]
            , viewInput "username" "text" "Username" model.username UsernameChanged
            , viewInput "password" "password" "Password" model.password PasswordChanged
            , button [ onClickNoBubble SubmitLogin, class "login__submit", disabled (model.username == "") ] [ text "Login" ]
            ]
        ]


viewInput : String -> String -> String -> String -> (String -> msg) -> Html msg
viewInput n t p v toMsg =
    input [ name n, type_ t, placeholder p, value v, onInput toMsg, class "login__input" ] []


onClickNoBubble : msg -> Html.Styled.Attribute msg
onClickNoBubble message =
    Html.Styled.Events.custom "click" (Json.Decode.succeed { message = message, stopPropagation = True, preventDefault = True })
