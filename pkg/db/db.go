package db

import (
	"database/sql"
    _ "github.com/lib/pq"
)

func OpenDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

    query := `CREATE TABLE IF NOT EXISTS requests(email TEXT, name TEXT, state TEXT, city TEXT, plan_id TEXT, installer_id TEXT)`
    _, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	return db, nil
}
