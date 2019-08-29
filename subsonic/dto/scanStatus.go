package dto

import (
	"encoding/xml"
)

// ScanStatus describes the progress of all music scans.
type ScanStatus struct {
	XMLName  xml.Name `xml:"scanStatus" json:"-"`
	Scanning bool     `xml:"scanning,attr" json:"scanning"`
	Count    int64    `xml:"count,attr" json:"count"`
}

// NewScanStatus makes a new ScanStatus DTO.
func NewScanStatus(scanning bool, count int64) ScanStatus {
	return ScanStatus{
		Scanning: scanning,
		Count:    count,
	}
}
