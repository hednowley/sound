package dto

import (
	"encoding/xml"
)

// License is a minimal Subsonic software license.
type License struct {
	XMLName xml.Name `xml:"license" json:"-"`
	Valid   bool     `xml:"valid,attr" json:"valid"`
}

// NewLicense makes a new License.
func NewLicense() *License {
	return &License{
		xml.Name{},
		true,
	}
}
