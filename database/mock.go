package database

import (
	"database/sql"
	"log"

	"github.com/hednowley/sound/config"
	testfixtures "gopkg.in/testfixtures.v2"
)

func NewMock() *Default {

	conn := "dbname=sound_test sslmode=disable user=postgres password=sound"

	// Initialise the database schema with gorm
	database, err := NewDefault(&config.Config{Db: conn})
	if err != nil {
		log.Fatal(err)
	}

	// Open the database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	// Insert test data
	fixtures, err := testfixtures.NewFolder(db, &testfixtures.PostgreSQL{}, "../testdata/dao")
	if err != nil {
		log.Fatal(err)
	}

	err = fixtures.Load()
	if err != nil {
		log.Fatal(err)
	}

	return database
}
