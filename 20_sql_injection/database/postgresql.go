package database

import (
	"context"
	"log"

	tracer "github.com/AVVKavvk/sql-injection/tracer"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDb() (*pgxpool.Pool, error) {
	// Database Connection with Tracer
	connStr := "postgres://myuser:mypassword@localhost:5432/test"
	config, _ := pgxpool.ParseConfig(connStr)
	config.ConnConfig.Tracer = &tracer.QueryLogger{}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {

		log.Fatal(err)
		return nil, err
	}
	return db, nil
}
