package main

import (
	"context"
	"gitlab.ozon.dev/VeneLooool/homework-3/internal/warehouse"
	"log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	_, err := warehouse.New(ctx)
	if err != nil {
		log.Fatalf("New WareHouse: %v", err)
	}
	log.Println("Start warehouse")
	<-ctx.Done()
}
