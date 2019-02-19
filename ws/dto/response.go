package dto

const version = "2.0"

// Reponse follows JSON-RPC syntax.
type Response struct {
	Error   string      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	ID      int         `json:"id,omitempty"`
	Version string      `json:"jsonrpc"`
}

func NewErrorResponse(message string, id int) *Response {
	return &Response{
		Error:   message,
		ID:      id,
		Version: version,
	}
}

func NewErrorNotification(message string) *Response {
	return &Response{
		Error:   message,
		Version: version,
	}
}

func NewResponse(result interface{}, id int) *Response {
	return &Response{
		Result:  result,
		ID:      id,
		Version: version,
	}
}

func NewNotification(result interface{}) *Response {
	return &Response{
		Result:  result,
		Version: version,
	}
}
