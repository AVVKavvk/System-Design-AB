package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDb() (*pgxpool.Pool, error) {
	// Database Connection with Tracer
	connStr := "postgres://myuser:mypassword@localhost:5432/url_shortener"
	config, _ := pgxpool.ParseConfig(connStr)

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {

		log.Fatal(err)
		return nil, err
	}
	return db, nil
}
