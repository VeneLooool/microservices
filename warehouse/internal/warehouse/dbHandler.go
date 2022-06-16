package warehouse

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	consts "gitlab.ozon.dev/VeneLooool/homework-3/const"
	"gitlab.ozon.dev/VeneLooool/homework-3/internal/db"
	"log"
)

type DB struct {
	conn *pgxpool.Pool
}

func NewDB(ctx context.Context) *DB {
	conn, err := db.New(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return &DB{conn: conn}
}

func (d *DB) GetAmount(ctx context.Context, productsId int64) (amount int, ok bool) {
	const query = `select amount from warehouse_products where productsId = $1;`
	err := d.conn.QueryRow(ctx, query, int(productsId)).Scan(&amount)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false
	}
	if err != nil {
		return 0, false
	}
	return amount, true
}

func (d *DB) SetAmount(ctx context.Context, newAmount int, productsId int) bool {
	const query = `update warehouse_products set amount = $1 where productsId = $2;`
	cmd, err := d.conn.Exec(ctx, query, newAmount, productsId)
	if err != nil {
		log.Println(err)
		for i := 0; i < 3; i++ {
			log.Println(i)
			_, err = d.conn.Exec(ctx, query, newAmount, productsId)
			if err == nil {
				break
			} else {
				log.Println(err)
			}
			if i == 2 {
				return false
			}
		}
	} else if cmd.RowsAffected() == 0 {
		return false
	}
	return true
}

func (d *DB) Reserve(ctx context.Context, order consts.Order) bool {
	amount, ok := d.GetAmount(ctx, order.ProductsId)
	if !ok || order.AmountOfProduct > int64(amount) {
		return false
	}
	if !d.SetAmount(ctx, amount-int(order.AmountOfProduct), int(order.ProductsId)) {
		return false
	}
	return true
}

func (d *DB) RemoveReserve(ctx context.Context, order consts.Order) bool {
	amount, ok := d.GetAmount(ctx, order.ProductsId)
	if !ok {
		return false
	}
	if !d.SetAmount(ctx, amount+int(order.AmountOfProduct), int(order.ProductsId)) {
		return false
	}
	return true
}
