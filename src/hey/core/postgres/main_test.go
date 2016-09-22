package postgres

import (
	"hey/storage"
	spostgres "hey/storage/postgres"
	"log"
	"os"
	"testing"
)

var db storage.DB

func TestMain(m *testing.M) {
	var err error
	db, err = spostgres.SetupPgFromENV()

	if err != nil {
		log.Panicln(err)
	}

	err = storage.ExecQueries(db, spostgres.SchemaBase)
	if err != nil {
		log.Panicln("[FAIL]", "create schema", err)
	}

	os.Exit(m.Run())
}
