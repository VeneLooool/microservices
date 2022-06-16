package prometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	OpsProcessed         prometheus.Counter
	SuccessfulProcessed  prometheus.Counter
	AmountOfUnsuccessful prometheus.Counter
	AmountOfCanceled     prometheus.Counter
	AmountOfRetry        prometheus.Counter
}

func NewMetrics(serviceName string) *Metrics {
	opsProcessed := promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_processed_ops_total", serviceName),
		Help: "The total number of processed events",
	})
	SuccessfulProcessed := promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_processed_successfule_total", serviceName),
		Help: "The total number of successfully processed events",
	})
	AmountOfUnsuccessful := promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_amount_of_unsuccessful_total", serviceName),
		Help: "The total number of unsuccessfully processed events",
	})
	AmountOfRetry := promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_amount_of_retry_total", serviceName),
		Help: "The total number of retry",
	})
	AmountOfCanceled := promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_amount_of_canceled_total", serviceName),
		Help: "The total number of canceled orders",
	})
	return &Metrics{
		OpsProcessed:         opsProcessed,
		SuccessfulProcessed:  SuccessfulProcessed,
		AmountOfUnsuccessful: AmountOfUnsuccessful,
		AmountOfCanceled:     AmountOfCanceled,
		AmountOfRetry:        AmountOfRetry,
	}
}
