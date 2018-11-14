package dto

import (
	"encoding/json"
	"fmt"
)

type Status int

const (
	Success Status = 0
	Fail    Status = 1
	Error   Status = 2
)

func MakeResponse(status Status, dto interface{}) string {

	var s string
	if status == Success {
		s = "success"
	} else if status == Fail {
		s = "fail"
	} else if status == Error {
		s = "error"
	}

	b, err := json.Marshal(dto)
	if err != nil {

		return fmt.Sprintf(`{ "status": "error", "message": "%v"}`, err)
	}

	return fmt.Sprintf(`{ "status": "%v", "data": %v}`, s, string(b))
}
