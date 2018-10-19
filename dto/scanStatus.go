package dto

import (
	"encoding/xml"
)

type ScanStatus struct {
	XMLName  xml.Name `xml:"scanStatus" json:"-"`
	Scanning bool     `xml:"scanning,attr" json:"scanning"`
	Count    int64    `xml:"count,attr" json:"count"`
}

func NewScanStatus(scanning bool, count int64) ScanStatus {
	return ScanStatus{
		Scanning: scanning,
		Count:    count,
	}
}
