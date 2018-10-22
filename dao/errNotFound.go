package dao

type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "Not found"
}

func IsErrNotFound(e error) bool {
	_, ok := e.(*ErrNotFound)
	return ok
}
