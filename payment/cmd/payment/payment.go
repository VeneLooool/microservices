package main

import (
	"context"
	"gitlab.ozon.dev/VeneLooool/homework-3/internal/payment"
	"log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	_, err := payment.New(ctx)
	if err != nil {
		log.Fatalf("New Payment: %v", err)
	}
	log.Println("Start payments")
	<-ctx.Done()
}
