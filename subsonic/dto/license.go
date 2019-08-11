package dto

import (
	"encoding/xml"
)

type License struct {
	XMLName xml.Name `xml:"license" json:"-"`
	Valid   bool     `xml:"valid,attr" json:"valid"`
}

func NewLicense() *License {
	return &License{
		xml.Name{},
		true,
	}
}
