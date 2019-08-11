package dto

type Directory struct {
	ID     string `xml:"id,attr" json:"id"`
	IsDir  bool   `xml:"isDir,attr" json:"isDir"`
	Parent string `xml:"parent,attr,omitempty" json:"parent,omitempty"`
}
