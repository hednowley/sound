package database

import (
	"database/sql"
	"log"
	"path/filepath"
	"runtime"

	"github.com/hednowley/sound/config"
	testfixtures "gopkg.in/testfixtures.v2"
)

// NewMock makes a new database seeded from test data.
// Note that a real database must exist at the connection string below.
func NewMock() *Default {

	conn := "host=localhost port=5432 user=sound password=sound dbname=sound_test sslmode=disable"

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
	dir := filepath.Dir(filename)
	dataPath := filepath.Join(dir, "..", "testdata", "dao")

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
