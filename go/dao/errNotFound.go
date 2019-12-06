package dao

// ErrNotFound is a special type of error for when data is missing.
type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "Not found"
}

// IsErrNotFound check is an error is an ErrNotFound.
func IsErrNotFound(e error) bool {
	_, ok := e.(*ErrNotFound)
	return ok
}
