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
	db, err = SetupPgFromENV()

	if err != nil {
		log.Panicln(err)
	}

	err = storage.ExecQueries(db, SchemaBase)
	if err != nil {
		log.Panicln("[FAIL]", "create schema", err)
	}

	os.Exit(m.Run())
}
