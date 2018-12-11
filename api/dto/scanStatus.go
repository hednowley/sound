package dto

type ScanStatus struct {
	Scanning bool  `xml:"scanning,attr" json:"scanning"`
	Count    int64 `xml:"count,attr" json:"count"`
}

func NewScanStatus(scanning bool, count int64) ScanStatus {
	return ScanStatus{
		Scanning: scanning,
		Count:    count,
	}
}
