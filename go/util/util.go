package util

import (
	"strconv"
	"strings"
)

// Min returns the smaller of two uints.
func Min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

// Max returns the larger of two uints.
func Max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

// Contains checks whether a slice of uints contains the given value.
func ContainsUint(slice []uint, value uint) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ContainsString(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func FormatUint(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

// ParseUint tries to parse a string into a uint. If this is not possible
// then the given default is returned instead.
func ParseUint(param string, defaultValue uint) uint {

	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return defaultValue
	}
	return uint(id)
}

// ParseBool tries to parse a string into a bool. If this is not possible
// then a nil pointer is returned instead.
func ParseBool(param string) *bool {

	param = strings.ToLower(param)
	if param == "true" || param == "1" {
		b := true
		return &b
	}
	if param == "false" || param == "0" {
		b := false
		return &b
	}

	return nil
}
