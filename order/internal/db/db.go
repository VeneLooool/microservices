package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"strings"
)

func buildDSN() string {
	hostAndPort := os.Getenv("POSTGRES_URL")
	host := hostAndPort[:strings.IndexRune(hostAndPort, ':')]
	port := hostAndPort[strings.IndexRune(hostAndPort, ':')+1:]
	log.Println("url:", hostAndPort)
	user := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbName)
}

func New(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, buildDSN())
}

func InitDb(ctx context.Context, d *pgxpool.Pool) (err error) {
	hostFirst := os.Getenv("SHARD1_URL")
	hostSecond := os.Getenv("SHARD2_URL")

	if _, err = d.Exec(ctx, fmt.Sprintf("create server shard1 foreign data wrapper postgres_fdw options (host '%s', port '5432', dbname 'example');", hostFirst)); err != nil {
		return err
	}
	if _, err = d.Exec(ctx, fmt.Sprintf("create server shard2 foreign data wrapper postgres_fdw options (host '%s', port '5432', dbname 'example');", hostSecond)); err != nil {
		return err
	}

	if _, err = d.Exec(ctx, "create foreign table warehouse_products_1 partition of warehouse_products for values with (modulus 2, remainder 0) server shard1;"); err != nil {
		return err
	}
	if _, err = d.Exec(ctx, "create foreign table warehouse_products_2 partition of warehouse_products for values with (modulus 2, remainder 1) server shard2;"); err != nil {
		return err
	}

	if _, err = d.Exec(ctx, "create foreign table order_active_orders_1 partition of order_active_orders for values with (modulus 2, remainder 0) server shard1;"); err != nil {
		return err
	}
	if _, err = d.Exec(ctx, "create foreign table order_active_orders_2 partition of order_active_orders for values with (modulus 2, remainder 1) server shard2;"); err != nil {
		return err
	}

	if _, err = d.Exec(ctx, "create user mapping for postgres server shard1 options (user 'postgres', password 'postgres');"); err != nil {
		return err
	}
	if _, err = d.Exec(ctx, "create user mapping for postgres server shard2 options (user 'postgres', password 'postgres');"); err != nil {
		return err
	}

	return nil
}
