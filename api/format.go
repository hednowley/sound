package api

type responseFormat int

const (
	jsonFormat responseFormat = 0
	xmlFormat  responseFormat = 1
)

var defaultFormat = xmlFormat
