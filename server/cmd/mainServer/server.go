package main

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	consts "gitlab.ozon.dev/VeneLooool/homework-3/const"
	"gitlab.ozon.dev/VeneLooool/homework-3/internal/db"
	"log"
	"time"
)

func createTopic(name string, admin sarama.ClusterAdmin) {
	err := admin.CreateTopic(name, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	if err != nil {
		log.Fatal("Error while creating topic: ", err.Error())
	}
}

func fullFillWarehouseDb(ctx context.Context) {
	c, err := db.New(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 1; i < 1000; i++ {
		const query = `insert into warehouse_products (
			productsId,
			amount
		) VALUES ($1, $2);`
		_, err := c.Exec(ctx, query,
			i,
			1e5,
		)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	//ctx := context.Background()

	//fullFillWarehouseDb(ctx)
	//log.Printf("table are filled")

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(consts.Brokers, cfg)
	if err != nil {
		log.Fatalf("sync kafka: %v", err)
	}
	var counter, products int64
	for {
		counter++
		products++
		if products == 999 {
			products = 1
		}

		order := consts.Order{
			Id:              counter,
			UserID:          int64(time.Now().Second()) + 1111,
			PaymentsData:    "Bam",
			PaymentsMethod:  consts.YellowBank,
			ProductsId:      products,
			AmountOfProduct: 1,
		}
		marshalMessage, err := json.Marshal(order)
		if err != nil {
			log.Println(err)
			continue
		}

		_, _, err = syncProducer.SendMessage(&sarama.ProducerMessage{
			Topic: consts.IncomeOrders,
			Key:   sarama.StringEncoder("order"),
			Value: sarama.ByteEncoder(marshalMessage),
		})

		if err != nil {
			log.Println(err)
			continue
		}
		time.Sleep(500 * time.Millisecond)
	}

}
