package util

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
func Contains(slice []uint, value uint) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
