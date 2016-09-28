package postgres

import (
	"log"
	"os"
	"strconv"
	"testing"

	pg "gopkg.in/jackc/pgx.v2"
)

var db *pg.ConnPool

func database() (*pg.ConnPool, error) {
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	pwd := os.Getenv("DB_PWD")
	_port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	port := DefaultPort
	if err == nil {
		port = uint16(_port)
	}
	_maxConns, err := strconv.Atoi(os.Getenv("DB_MAXCONNS"))
	maxConns := DefaultMaxConns
	if err == nil {
		maxConns = _maxConns
	}

	return pg.NewConnPool(pg.ConnPoolConfig{
		ConnConfig: pg.ConnConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: pwd,
			Database: name,
		},
		MaxConnections: maxConns,
	})
}

func TestMain(m *testing.M) {
	var err error
	db, err = database()

	if err != nil {
		log.Panicln(err)
	}

	if err := SetupSchema(db); err != nil {
		log.Panicln(err)
	}

	os.Exit(m.Run())
}
