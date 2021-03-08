// moodle/databasse
package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	pool *sql.DB // Database connection pool.
	ctx  context.Context
}

func Open(DriverName string, DSN string, ctx context.Context) (db Database) {
	var err error
	db.ctx = ctx
	db.pool, err = sql.Open(DriverName, DSN)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal("unable to use data source name", err)
	}

	db.pool.SetMaxIdleConns(3)
	db.pool.SetMaxOpenConns(3)

	db.Ping()

	return
}

func (db *Database) Close() {
	db.pool.Close()
}

// Ping the database to verify DSN provided by the user is valid and the
// server accessible. If the ping fails exit the program with an error.
func (db *Database) Ping() {
	ctx, cancel := context.WithTimeout(db.ctx, 1*time.Second)
	defer cancel()

	if err := db.pool.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
}

func (db *Database) Query(query string) (*sql.Rows, error) {
	return db.pool.QueryContext(db.ctx, query)
}
