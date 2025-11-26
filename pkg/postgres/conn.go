package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

func NewPostgresConn(host string, port int, user, password, dbname string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	dbx, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = dbx.Ping()
	if err != nil {
		return nil, err
	}
	dbx.SetMaxOpenConns(25)
	dbx.SetMaxIdleConns(5)
	dbx.SetConnMaxLifetime(time.Minute * 5)
	return dbx, nil
}
