package api

import (
	"github.com/hednowley/sound/dto"
)

type Response struct {
	Body      interface{}
	IsSuccess bool
}

func NewErrorReponse(code dto.ErrorCode, message string) *Response {
	return &Response{
		Body:      dto.NewError(code, message),
		IsSuccess: false,
	}
}

func NewSuccessfulReponse(body interface{}) *Response {
	return &Response{
		Body:      body,
		IsSuccess: true,
	}
}
