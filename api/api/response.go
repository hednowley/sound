package api

type Status int

const (
	Success Status = 0
	Fail    Status = 1
	Error   Status = 2
)

// Reponse follows JSend syntax.
type Response struct {
	Body   interface{}
	Status Status
}

func NewErrorReponse(message string) *Response {
	return &Response{
		Body:   message,
		Status: Error,
	}
}

func NewSuccessfulReponse(body interface{}) *Response {
	return &Response{
		Body:   body,
		Status: Success,
	}
}
