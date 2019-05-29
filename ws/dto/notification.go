package dto

type Notification struct {
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Version string                 `json:"jsonrpc"`
}

func NewNotification(method string, params map[string]interface{}) *Notification {
	return &Notification{
		Method:  method,
		Params:  params,
		Version: "2.0",
	}
}

func NewScanStatusNotification(scanning bool, count int64) *Notification {
	return NewNotification("scanStatus", map[string]interface{}{
		"scanning": scanning,
		"count":    count,
	})
}
