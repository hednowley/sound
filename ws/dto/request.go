package dto

import "encoding/json"

type Request struct {
	Method  string                      `json:"method"`
	Params  map[string]*json.RawMessage `json:"params"`
	ID      int                         `json:"id"`
	Version string                      `json:"jsonrpc"`
}
