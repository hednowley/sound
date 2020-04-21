package database

import (
	"database/sql"
	"log"
	"path/filepath"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/projectpath"
)

// NewMock makes a new database seeded from test data.
// Note that a real database must exist at the connection string below.
func NewMock() *Default {

	conn := "host=localhost port=5432 user=sound password=sound dbname=sound_test sslmode=disable"

	database, err := NewDefault(&config.Config{Db: conn})
	if err != nil {
		log.Fatal(err)
	}

	// Open the database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	dataDir := filepath.Join(projectpath.Root, "testdata", "dao")

	// Insert test data
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.UseAlterConstraint(),
		testfixtures.Directory(dataDir),
		testfixtures.ResetSequencesTo(10000),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = fixtures.Load()
	if err != nil {
		log.Fatal(err)
	}

	return database
}
