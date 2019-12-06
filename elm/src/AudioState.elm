module AudioState exposing (State(..))


type State
    = Loading
    | Loaded { duration : Maybe Float }
    | Playing { time : Float, duration : Maybe Float, paused : Bool }
