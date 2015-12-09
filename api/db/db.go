package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	db := Get()
	buildTables(db)
}

func buildTables(db *sqlx.DB) {
	userTable :=
		`CREATE TABLE users (
	id text NOT NULL,
	name text,
	email text,
	jwt text,
	passwordhash text);`

	db.MustExec(userTable)
}

// Get loads the db from file
func Get() *sqlx.DB {
	// _, filename, _, _ := runtime.Caller(0)
	// dbPath := path.Join(path.Dir(filename), "./tmp/gorm.db")
	return sqlx.MustConnect("sqlite3", ":memory:")
}
