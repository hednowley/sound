package dao

import (
	"database/sql"
	"log"

	"github.com/hednowley/sound/config"

	testfixtures "gopkg.in/testfixtures.v2"
)

func NewMockDatabase() *Database {

	conn := "dbname=sound_test sslmode=disable user=postgres password=sound"

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	// creating the context that hold the fixtures
	// see about all compatible databases in this page below
	fixtures, err := testfixtures.NewFolder(db, &testfixtures.PostgreSQL{}, "../testdata")
	if err != nil {
		log.Fatal(err)
	}

	err = fixtures.Load()
	if err != nil {
		log.Fatal(err)
	}

	c := config.Config{
		Db: conn,
	}

	database, err := NewDatabase(&c)
	if err != nil {
		log.Fatal(err)
	}
	return database
}
