package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func buildDSN() string {
	return fmt.Sprintf("host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
}

func New(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, buildDSN())
}
