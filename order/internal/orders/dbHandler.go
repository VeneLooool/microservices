package orders

import (
	"context"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	consts "gitlab.ozon.dev/VeneLooool/homework-3/const"
	"gitlab.ozon.dev/VeneLooool/homework-3/internal/db"
	"gitlab.ozon.dev/VeneLooool/homework-3/tools/memcached"
	"log"
)

type DB struct {
	conn  *pgxpool.Pool
	cache Cache
}

func NewDB(ctx context.Context) *DB {
	conn, err := db.New(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	cache := memcached.New()
	return &DB{conn: conn, cache: cache}
}

func (d *DB) GetOrderFromDb(ctx context.Context, id int) (order consts.Order, ok bool) {

	const query = `select id, 
							userId, 
							paymentsMethod,
							paymentsData,
							productsId,
							amount,
					from order_active_orders where id = $1;`

	err := d.conn.QueryRow(ctx, query, int(id)).Scan(
		&order.Id,
		&order.UserID,
		&order.PaymentsMethod,
		&order.PaymentsData,
		&order.ProductsId,
		&order.AmountOfProduct,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return consts.Order{}, false
	}
	if err != nil {
		return consts.Order{}, false
	}
	return order, true
}

func (d *DB) GetOrder(ctx context.Context, id int) (*consts.Order, bool) {
	var (
		order consts.Order
		err   error
		ok    bool
	)
	order, err = d.cache.Get(ctx, id)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			order, ok = d.GetOrderFromDb(ctx, id)
			if !ok {
				return nil, false
			}

			if err = d.cache.Set(ctx, order); err != nil {
				log.Println(err)
			}
		} else {
			log.Println(err)
			order, ok = d.GetOrderFromDb(ctx, id)
			if !ok {
				return nil, false
			}
		}
	}

	return &order, true
}

func (d *DB) UpdateOrderInDb(ctx context.Context, order consts.Order) bool {
	const query = `update order_active_orders 
                   set		userId = $1, 
							paymentsMethod = $2,
							paymentsData = $3,
							productsId = $4,
							amount = $5, where id= $6;`

	cmd, err := d.conn.Exec(ctx, query,
		order.UserID,
		order.PaymentsMethod,
		order.PaymentsData,
		order.ProductsId,
		order.AmountOfProduct,
	)

	if cmd.RowsAffected() == 0 || err != nil {
		return false
	}
	return true
}

func (d *DB) UpdateOrder(ctx context.Context, order consts.Order) bool {
	err := d.cache.Delete(ctx, int(order.Id))
	if err != nil {
		if !errors.Is(err, memcache.ErrCacheMiss) {
			log.Println(err)
		}
	}
	if ok := d.UpdateOrderInDb(ctx, order); !ok {
		return false
	}
	return true
}

func (d *DB) CreateOrderInDb(ctx context.Context, order consts.Order) bool {
	const query = `insert into order_active_orders (
			id,	
			userId,
			paymentsMethod,
			paymentsData,
			productsId,
			amount
		) VALUES ($1, $2, $3, $4, $5, $6);`
	_, err := d.conn.Exec(ctx, query,
		order.Id,
		order.UserID,
		order.PaymentsMethod,
		order.PaymentsData,
		order.ProductsId,
		order.AmountOfProduct,
	)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (d *DB) CreateOrder(ctx context.Context, order consts.Order) (ok bool) {
	if ok = d.CreateOrderInDb(ctx, order); !ok {
		return false
	}
	if err := d.cache.Set(ctx, order); err != nil {
		log.Println(err)
	}
	return true
}

func (d *DB) DeleteOrderFromDb(ctx context.Context, id int) bool {
	const query = `
		delete from order_active_orders
		where id = $1;`

	cmd, err := d.conn.Exec(ctx, query, id)
	if cmd.RowsAffected() == 0 || err != nil {
		return false
	}
	return true
}

func (d *DB) DeleteOrder(ctx context.Context, id int) bool {

	err := d.cache.Delete(ctx, id)
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		log.Println(err)
	}

	if ok := d.DeleteOrderFromDb(ctx, id); !ok {
		return false
	}
	return true
}
