package orders

import (
	"context"
	consts "gitlab.ozon.dev/VeneLooool/homework-3/const"
)

type Cache interface {
	Get(context.Context, int) (consts.Order, error)
	Set(context.Context, consts.Order) error
	Delete(ctx context.Context, id int) error
}
