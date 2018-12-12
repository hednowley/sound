package database

import (
	"database/sql"
	"log"
	"runtime"

	"github.com/hednowley/sound/config"
	testfixtures "gopkg.in/testfixtures.v2"
)

// NewMock makes a new database seeded from test data.
// Note that a real database must exist at the connection string below.
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

	// Get test data path (found relative to this file)
	_, filename, _, _ := runtime.Caller(0)
	dataPath := filename + "/../../testdata/dao"

	// Insert test data
	fixtures, err := testfixtures.NewFolder(db, &testfixtures.PostgreSQL{}, dataPath)
	if err != nil {
		log.Fatal(err)
	}

	err = fixtures.Load()
	if err != nil {
		log.Fatal(err)
	}

	return database
}
