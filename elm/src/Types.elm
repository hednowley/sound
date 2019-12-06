module Types exposing (Update, UpdateWithReturn, combine, combineMany, noOp)


type alias Update model msg =
    model -> ( model, Cmd msg )


type alias UpdateWithReturn model msg return =
    model -> ( ( model, Cmd msg ), return )


noOp : Update model msg
noOp model =
    ( model, Cmd.none )


{-| Running one update and then another.
-}
combine : Update model msg -> Update model msg -> Update model msg
combine first second model =
    let
        ( modelA, cmdA ) =
            first model

        ( modelB, cmdB ) =
            second modelA
    in
    ( modelB, Cmd.batch [ cmdB, cmdA ] )


{-| Running all the given updates in the order they appear in the list.
-}
combineMany : List (Update model msg) -> Update model msg
combineMany updates =
    List.foldr combine noOp updates
