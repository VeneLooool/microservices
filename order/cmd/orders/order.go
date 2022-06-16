package main

import (
	"context"
	"gitlab.ozon.dev/VeneLooool/homework-3/internal/orders"
	"log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	_, err := orders.New(ctx)
	if err != nil {
		log.Fatalf("New active orders: %v", err)
	}
	log.Println("Start order")
	<-ctx.Done()
}
