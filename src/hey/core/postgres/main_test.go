package postgres

import (
	"hey/storage"
	"hey/storage/postgres"
	"log"
	"os"
	"testing"
)

var db storage.DB

func TestMain(m *testing.M) {
	var err error
	db, err = postgres.SetupPg(
		"localhost",
		"dbname",
		"dbuser",
		"dbuserpassword",
		5432,
		10,
	)

	if err != nil {
		log.Panicln(err)
	}

	err = storage.ExecQueries(db, postgres.SchemaBase)
	if err != nil {
		log.Panicln("[FAIL]", "create schema", err)
	}

	os.Exit(m.Run())
}
