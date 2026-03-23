package tracer

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type QueryLogger struct{}

func (t *QueryLogger) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	fmt.Printf("\n[DATABASE LOG] Executing SQL: %s\n", data.SQL)
	fmt.Printf("[DATABASE LOG] With Arguments: %v\n", data.Args)
	return ctx
}
func (t *QueryLogger) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {}
