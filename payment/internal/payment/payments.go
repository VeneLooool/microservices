package payment

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

type Payments struct {
	producer       sarama.SyncProducer
	incomeConsumer *IncomeHandler
	retryConsumer  *RetryHandler
	metrics        *prometheus.Metrics
	metricsMut     sync.Mutex
}

func New(ctx context.Context) (*Payments, error) {
	k := http.NewServeMux()
	metrics := prometheus.NewMetrics("payment")
	k.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8081", k)

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(consts.Brokers, cfg)
	if err != nil {
		return nil, err
	}

	iHandler := &IncomeHandler{
		prod:    producer,
		metrics: metrics,
	}
	go incHandler(ctx, iHandler, cfg)

	retryHandler := &RetryHandler{
		prod:    producer,
		metrics: metrics,
	}
	go rHandler(ctx, retryHandler, cfg)

	return &Payments{
		incomeConsumer: iHandler,
		retryConsumer:  retryHandler,
		producer:       producer,
		metrics:        metrics,
	}, nil
}

func incHandler(ctx context.Context, handler *IncomeHandler, cfg *sarama.Config) {
	income, err := sarama.NewConsumerGroup(consts.Brokers, "payments", cfg)
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
func rHandler(ctx context.Context, handler *RetryHandler, cfg *sarama.Config) {
	retry, err := sarama.NewConsumerGroup(consts.Brokers, "retryPayments", cfg)
	if err != nil {
		log.Fatalf("retry handler: %v", err)
	}

	for {
		err := retry.Consume(ctx, []string{consts.RetrySendToCreateOrder}, handler)
		if err != nil {
			log.Printf("retry consumer error: %v", err)
			time.Sleep(time.Second * 5)
		}
	}
}
