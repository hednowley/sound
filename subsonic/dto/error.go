package dto

import (
	"encoding/xml"
)

type Error struct {
	XMLName xml.Name `xml:"error" json:"-"`
	Code    int      `xml:"code,attr" json:"code"`
	Message string   `xml:"message,attr" json:"message"`
}

type ErrorCode int

const (
	Generic              ErrorCode = 0
	MissingParameter     ErrorCode = 10
	ClientTooOld         ErrorCode = 20
	ServerTooOld         ErrorCode = 30
	WrongCredentials     ErrorCode = 40
	TokenAuthUnsupported ErrorCode = 41
	NotAuthorised        ErrorCode = 50
	BadLicense           ErrorCode = 60
	NotFound             ErrorCode = 70
)

func NewError(code ErrorCode, message string) Error {

	return Error{
		Code:    int(code),
		Message: message,
	}
}
