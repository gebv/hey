package storage

type ExecDetails interface {
	RowsAffected() int64
}

type DB interface {
	Exec(sql string, arguments ...interface{}) (ExecDetails, error)
	Query(sql string, args ...interface{}) (Rows, error)
	QueryRow(sql string, args ...interface{}) Row
}

type Rows interface {
	Row
	Next() bool
}

type Row interface {
	Scan(dest ...interface{}) error
}
