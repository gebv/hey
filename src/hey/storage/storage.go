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

type TX interface {
	DB

	Commit() error
	Rollback() error
}

type BeginTX interface {
	Begin() (TX, error)
}

func ExecQueries(_conn DB, queries []string) (err error) {
	var conn DB = _conn

	if _, ok := _conn.(BeginTX); ok {
		conn, err = _conn.(BeginTX).Begin()
		if err != nil {
			return
		}
	}

	for _, query := range queries {
		if _, err = conn.Exec(query); err != nil {

			if _, ok := conn.(TX); ok {
				conn.(TX).Rollback()
			}

			return
		}
	}

	if _, ok := conn.(TX); ok {
		conn.(TX).Commit()
	}

	return
}
