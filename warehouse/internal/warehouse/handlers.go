package warehouse

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/VeneLooool/homework-3/api/prometheus"
	consts "gitlab.ozon.dev/VeneLooool/homework-3/const"
	"log"
	"time"
)

type IncomeHandler struct {
	prod    sarama.SyncProducer
	metrics *prometheus.Metrics
	data    *DB
}
type ResetHandler struct {
	prod    sarama.SyncProducer
	metrics *prometheus.Metrics
	data    *DB
}

func (i *IncomeHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (i *IncomeHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (i *IncomeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		i.metrics.OpsProcessed.Inc()
		var order consts.Order
		ctx := context.Background()

		err := json.Unmarshal(msg.Value, &order)
		if err != nil {
			log.Printf("income data %v: %v", string(msg.Value), err)
			continue
		}

		if !i.data.Reserve(ctx, order) {
			log.Printf("Order %v is not reserved", order.Id)
			return nil
		}
		log.Printf("Order %v reserved", order.Id)

		js, err := json.Marshal(order)
		if err != nil {
			log.Printf("income data %v: %v", order.Id, err)
		}
		_, _, err = i.prod.SendMessage(&sarama.ProducerMessage{
			Topic:     consts.IncomePayments,
			Value:     sarama.ByteEncoder(js),
			Timestamp: time.Now(),
		})

		if err != nil {
			log.Printf("order warehouses: %v", err)
			if !i.data.RemoveReserve(ctx, order) {
				log.Printf("Can't remove reserve for order %v", order.Id)
			}
			i.metrics.AmountOfUnsuccessful.Inc()
			continue
		}
		i.metrics.SuccessfulProcessed.Inc()
	}

	return nil
}

func (r *ResetHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (r *ResetHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (r *ResetHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var order consts.Order
		ctx := context.Background()
		err := json.Unmarshal(msg.Value, &order)
		if err != nil {
			log.Printf("data %v: %v", string(msg.Value), err)
			continue
		}
		r.metrics.AmountOfCanceled.Inc()
		if !r.data.RemoveReserve(ctx, order) {
			log.Printf("bad order: %v", order.Id)
		}
		log.Printf("order canceled: %v", order.Id)
	}
	return nil
}
