package dal

import (
	"github.com/hednowley/sound/database"
)

func NewMock() *DAL {
	return &DAL{
		db: database.NewMock(),
	}
}
