package main

import (
	"context"
	"gitlab.ozon.dev/VeneLooool/homework-3/internal/db"
	"log"
)

func main() {
	ctx := context.Background()
	pool, err := db.New(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	if err = db.InitDb(ctx, pool); err != nil {
		log.Fatalln(err)
	}
}
