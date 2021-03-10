// moodle/databasse
package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	pool *sql.DB // Database connection pool.
	ctx  context.Context
}

func Open(ctx context.Context, DriverName string, DSN string) (db Database, err error) {
	db.ctx = ctx
	db.pool, err = sql.Open(DriverName, DSN)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		return
	}

	db.pool.SetMaxIdleConns(3)
	db.pool.SetMaxOpenConns(3)

	err = db.Ping()

	return
}

func (db *Database) Close() {
	db.pool.Close()
}

// Ping the database to verify DSN provided by the user is valid and the
// server accessible. If the ping fails exit the program with an error.
func (db *Database) Ping() (err error) {
	ctx, cancel := context.WithTimeout(db.ctx, 1*time.Second)
	defer cancel()

	err = db.pool.PingContext(ctx)

	return
}

func (db *Database) Query(query string) (*sql.Rows, error) {
	return db.pool.QueryContext(db.ctx, query)
}

func (db *Database) QueryRow(query string) *sql.Row {
	return db.pool.QueryRowContext(db.ctx, query)
}

func (db *Database) Exec(query string) (sql.Result, error) {
	return db.pool.ExecContext(db.ctx, query)
}
