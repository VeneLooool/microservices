package main

import (
	"context"
	"gitlab.ozon.dev/VeneLooool/homework-3/internal/notification"
	"log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	_, err := notification.New(ctx)
	if err != nil {
		log.Fatalf("New notification: %v", err)
	}
	log.Println("Start notification")
	<-ctx.Done()
}
