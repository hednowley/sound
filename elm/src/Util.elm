module Util exposing (insertMany)

import Dict exposing (Dict)


insertMany : (a -> comparable) -> (a -> value) -> List a -> Dict comparable value -> Dict comparable value
insertMany toKey toValue list dict =
    List.foldl
        (\a -> \d -> Dict.insert (toKey a) (toValue a) d)
        dict
        list
