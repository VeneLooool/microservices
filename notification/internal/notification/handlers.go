package notification

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/VeneLooool/homework-3/api/prometheus"
	"gitlab.ozon.dev/VeneLooool/homework-3/api/randService"
	consts "gitlab.ozon.dev/VeneLooool/homework-3/const"
	"log"
	"time"
)

type IncomeHandler struct {
	prod    sarama.SyncProducer
	metrics *prometheus.Metrics
}

type RetryHandler struct {
	prod    sarama.SyncProducer
	metrics *prometheus.Metrics
}

func (r *RetryHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (r *RetryHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (r *RetryHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var retryOrder consts.RetryOrder

		err := json.Unmarshal(msg.Value, &retryOrder)
		if err != nil {
			log.Printf("Retry data %v: %v", string(msg.Value), err)
			continue
		}

		if int64(msg.Timestamp.Second()-time.Now().Second()) < retryOrder.Delay {
			time.Sleep(time.Second * 5)
			err = r.RetryAgain(&retryOrder, msg.Timestamp)
			if err != nil {
				log.Printf("Retry data %v: %v", string(msg.Value), err)
			}
			continue
		}

		if retryOrder.AmountOfRetry >= 3 {
			log.Printf("Retry data more than 3 times retry %v: %v", string(msg.Value), err)
			r.metrics.AmountOfUnsuccessful.Inc()
			continue
		}

		retryOrder.AmountOfRetry++
		retryOrder.Delay = int64(time.Second * 20)
		time.Sleep(time.Second * 5)

		err = r.RetryAgain(&retryOrder, time.Now())
		if err != nil {
			log.Printf("Retry data %v: %v", retryOrder.Id, err)
		}
		log.Printf("attempt to send notification succesfule for order %v", retryOrder.Id)
		r.metrics.SuccessfulProcessed.Inc()
	}
	return nil
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

		err := json.Unmarshal(msg.Value, &order)
		if err != nil {
			log.Printf("income data %v: %v", string(msg.Value), err)
			continue
		}
		err = randService.RAND(errors.New("rand broke notification"))

		if err != nil {
			_, _, err := i.prod.SendMessage(NewRetryMes(&order, time.Now()))
			i.metrics.AmountOfRetry.Inc()
			if err != nil {
				log.Printf("order retry: %v", err)
			}
		}
		log.Printf("Notification for order %v, has been sended", order.Id)
		i.metrics.SuccessfulProcessed.Inc()
	}

	return nil
}

func NewRetryMes(order *consts.Order, timeStamp time.Time) *sarama.ProducerMessage {
	retry := consts.RetryOrder{
		Id:              order.Id,
		UserID:          order.UserID,
		PaymentsMethod:  order.PaymentsMethod,
		PaymentsData:    order.PaymentsData,
		ProductsId:      order.ProductsId,
		AmountOfProduct: order.AmountOfProduct,
		AmountOfRetry:   0,
		Delay:           int64(time.Second * 5),
	}
	marshalMessage, _ := json.Marshal(retry)
	return &sarama.ProducerMessage{
		Topic:     consts.RetryNotification,
		Key:       sarama.StringEncoder("retry"),
		Value:     sarama.ByteEncoder(marshalMessage),
		Timestamp: timeStamp,
	}
}

func (r *RetryHandler) RetryAgain(order *consts.RetryOrder, t time.Time) error {

	marshalMessage, err := json.Marshal(order)
	if err != nil {
		return err
	}

	_, _, err = r.prod.SendMessage(&sarama.ProducerMessage{
		Topic:     consts.RetryNotification,
		Key:       sarama.StringEncoder("retry"),
		Value:     sarama.ByteEncoder(marshalMessage),
		Timestamp: t,
	})
	if err != nil {
		return err
	}

	return nil
}
