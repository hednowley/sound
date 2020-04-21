package dal

import (
	"github.com/hednowley/sound/database"
)

func NewMock() *DAL {
	return &DAL{
		Db: database.NewMock(),
	}
}
