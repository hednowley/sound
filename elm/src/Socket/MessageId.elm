module Socket.MessageId exposing (MessageId(..), getRawMessageId)


type MessageId
    = MessageId Int


getRawMessageId : MessageId -> Int
getRawMessageId messageId =
    let
        (MessageId id) =
            messageId
    in
    id
