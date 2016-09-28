package postgres

import (
	"errors"

	pg "gopkg.in/jackc/pgx.v2"
)

var (
	ErrWantTx = errors.New("is only supported in transaction")
)

func SetupSchema(pg *pg.ConnPool) error {
	tx, err := pg.Begin()

	if err != nil {
		return err
	}

	for _, query := range SchemaBase {
		if _, err = tx.Exec(query); err != nil {

			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
