package memcached

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	consts "gitlab.ozon.dev/VeneLooool/homework-3/const"
	"os"
)

type cache struct {
	client *memcache.Client
}

func New() *cache {
	mc := memcache.New(os.Getenv("memcachedURL"))
	return &cache{mc}
}

func (c *cache) Get(ctx context.Context, id int) (res consts.Order, err error) {
	item, err := c.client.Get(fmt.Sprint(id))
	if err != nil {
		return
	}

	err = json.Unmarshal(item.Value, &res)
	if err != nil {
		return
	}
	return
}

func (c *cache) Set(ctx context.Context, p consts.Order) (err error) {

	data, err := json.Marshal(p)
	if err != nil {
		return
	}
	return c.client.Set(&memcache.Item{
		Key:   fmt.Sprint(p.Id),
		Value: data,
	})
}

func (c *cache) Delete(ctx context.Context, id int) error {
	return c.client.Delete(fmt.Sprint(id))
}
