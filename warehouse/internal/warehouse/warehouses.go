package warehouse

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/VeneLooool/homework-3/api/prometheus"
	consts "gitlab.ozon.dev/VeneLooool/homework-3/const"
	"log"
	"net/http"
	"sync"
	"time"
)

type WareHouse struct {
	data           *DB
	producer       sarama.SyncProducer
	incomeConsumer *IncomeHandler
	resetConsumer  *ResetHandler
	metrics        *prometheus.Metrics
	metricsMut     sync.Mutex
}

func New(ctx context.Context) (*WareHouse, error) {
	k := http.NewServeMux()
	metrics := prometheus.NewMetrics("warehouse")
	k.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", k)

	data := NewDB(ctx)

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(consts.Brokers, cfg)
	if err != nil {
		return nil, err
	}

	iHandler := &IncomeHandler{
		prod:    producer,
		data:    data,
		metrics: metrics,
	}
	go incHandler(ctx, iHandler, cfg)

	rHandler := &ResetHandler{
		prod:    producer,
		data:    data,
		metrics: metrics,
	}
	go resetHandler(ctx, rHandler, cfg)

	return &WareHouse{
		data:           data,
		incomeConsumer: iHandler,
		resetConsumer:  rHandler,
		producer:       producer,
		metrics:        metrics,
	}, nil
}

func incHandler(ctx context.Context, handler *IncomeHandler, cfg *sarama.Config) {
	income, err := sarama.NewConsumerGroup(consts.Brokers, "wareHouse", cfg)
	if err != nil {
		log.Fatalf("income handler: %v", err)
	}

	for {
		err := income.Consume(ctx, []string{consts.IncomeOrders}, handler)
		if err != nil {
			log.Printf("income consumer error: %v", err)
			time.Sleep(time.Second * 5)
		}
	}

}
func resetHandler(ctx context.Context, handler *ResetHandler, cfg *sarama.Config) {
	reset, err := sarama.NewConsumerGroup(consts.Brokers, "wareHouseReset", cfg)
	if err != nil {
		log.Fatalf("retry handler: %v", err)
	}

	for {
		err := reset.Consume(ctx, []string{consts.ResetOrders}, handler)
		if err != nil {
			log.Printf("reset orders error: %v", err)
			time.Sleep(time.Second * 5)
		}
	}
}
