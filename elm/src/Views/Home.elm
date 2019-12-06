module Views.Home exposing (view)

import Html.Styled exposing (Html, button, div, input, label, span, text)
import Html.Styled.Attributes exposing (checked, class, type_)
import Html.Styled.Events exposing (onClick)
import Model exposing (Model)
import Msg exposing (Msg(..))
import String exposing (fromInt)


view : Model -> Html Msg
view model =
    div [ class "home__wrap" ]
        [ span [] [ text <| "Scanned: " ++ String.fromInt model.scanCount ]
        , button [ onClick LogOut ] [ text "Log out" ]
        , checkboxInput "Update?" model.scanShouldUpdate ToggleScanUpdate
        , checkboxInput "Delete?" model.scanShouldDelete ToggleScanDelete
        , button [ onClick StartScan ] [ text "Start scan" ]
        ]


checkboxInput : String -> Bool -> msg -> Html msg
checkboxInput name isChecked msg =
    label [] [ input [ checked isChecked, type_ "checkbox", onClick msg ] [], text name ]
