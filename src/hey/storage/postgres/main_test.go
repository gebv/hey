package postgres

import (
	"hey/storage"
	"log"
	"os"
	"testing"
)

var db storage.DB

func TestMain(m *testing.M) {
	var err error
	db, err = SetupPg(
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

	err = storage.ExecQueries(db, schemaBase)
	if err != nil {
		log.Panicln("[FAIL]", "create schema", err)
	}

	os.Exit(m.Run())
}
