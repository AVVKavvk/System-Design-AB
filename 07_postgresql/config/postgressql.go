package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBRegistry struct {
	MyApp *pgxpool.Pool
}

type Tables struct {
	AvailableTables []string
}

func (t *Tables) Add(table string) {
	t.AvailableTables = append(t.AvailableTables, table)
}
func (t *Tables) Remove(table string) {
	for i, v := range t.AvailableTables {
		if v == table {
			t.AvailableTables = append(t.AvailableTables[:i], t.AvailableTables[i+1:]...)
		}
	}
}
func (t *Tables) Exists(table string) bool {
	for _, v := range t.AvailableTables {
		if v == table {
			return true
		}
	}
	return false
}

var MyAppTables *Tables

func GetDBRegistry(ctx context.Context, myAppDbURL string) (*DBRegistry, error) {

	myPool, err := pgxpool.New(ctx, myAppDbURL)
	if err != nil {
		return nil, err
	}

	return &DBRegistry{
		MyApp: myPool,
	}, nil
}

func init() {
	MyAppTables = &Tables{
		AvailableTables: []string{},
	}
}
