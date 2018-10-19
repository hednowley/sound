package dto

import (
	"testing"
)

func TestScanningStatus(t *testing.T) {

	DTO := NewScanStatus(true, 45485)

	xml := `<scanStatus scanning="true" count="45485"></scanStatus>`

	json := `
	{
		"scanning":true,
		"count":45485
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestNotScanningStatus(t *testing.T) {

	DTO := NewScanStatus(false, 0)

	xml := `<scanStatus scanning="false" count="0"></scanStatus>`

	json := `
	{
		"scanning":false,
		"count":0
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
