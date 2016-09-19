package postgres

import (
	"hey/storage"

	pg "gopkg.in/jackc/pgx.v2"
)

var (
	_ storage.DB      = (*Conn)(nil)
	_ storage.BeginTX = (*Conn)(nil)
	_ storage.TX      = (*ConnTx)(nil)
)

// SetupPg
func SetupPg(
	host, name, user, pwd string,
	port uint16,
	maxConns int,
) (storage.DB, error) {
	conn, err := pg.NewConnPool(pg.ConnPoolConfig{
		ConnConfig: pg.ConnConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: pwd,
			Database: name,
		},
		MaxConnections: maxConns,
	})
	return &Conn{db: conn}, err
}

type Conn struct {
	db *pg.ConnPool
}

func (c *Conn) Exec(sql string, arguments ...interface{}) (storage.ExecDetails, error) {
	return c.db.Exec(sql, arguments...)
}

func (c *Conn) Query(sql string, arguments ...interface{}) (storage.Rows, error) {
	return c.db.Query(sql, arguments...)
}

func (c *Conn) QueryRow(sql string, arguments ...interface{}) storage.Row {
	return c.db.QueryRow(sql, arguments...)
}

func (c *Conn) Begin() (storage.TX, error) {
	tx, err := c.db.Begin()
	return &ConnTx{tx: tx}, err
}

type ConnTx struct {
	tx *pg.Tx
}

func (c *ConnTx) Exec(sql string, arguments ...interface{}) (storage.ExecDetails, error) {
	return c.tx.Exec(sql, arguments...)
}

func (c *ConnTx) Query(sql string, arguments ...interface{}) (storage.Rows, error) {
	return c.tx.Query(sql, arguments...)
}

func (c *ConnTx) QueryRow(sql string, arguments ...interface{}) storage.Row {
	return c.tx.QueryRow(sql, arguments...)
}

func (c *ConnTx) Commit() error {
	return c.tx.Commit()
}

func (c *ConnTx) Rollback() error {
	return c.tx.Rollback()
}
