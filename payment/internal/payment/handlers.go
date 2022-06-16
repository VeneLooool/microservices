package payment

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/VeneLooool/homework-3/api/paymethods/blueBank"
	"gitlab.ozon.dev/VeneLooool/homework-3/api/paymethods/yellowBank"
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

		_, _, err = r.prod.SendMessage(NewAttemptSendMes(&retryOrder))

		if err == nil {
			err = randService.RAND(errors.New("rand error"))
		}
		if err != nil {
			retryOrder.AmountOfRetry++
			retryOrder.Delay = int64(time.Second * 20)
			time.Sleep(time.Second * 5)

			err = r.RetryAgain(&retryOrder, time.Now())
			if err != nil {
				log.Printf("Retry data %v: %v", retryOrder.Id, err)
			}
			continue
		}
		log.Printf("attempt to send to order creation successful %v", retryOrder.Id)
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

		if err = checkOut(&order); err != nil {
			log.Printf("order %v payment: %v", order.Id, err)
			err := i.resetOperation(&order)
			i.metrics.AmountOfCanceled.Inc()
			if err != nil {
				log.Printf("order payment: %v", err)
			}
			continue
		}

		log.Printf("Order %v has been paid", order.Id)

		js, _ := json.Marshal(order)
		_, _, err = i.prod.SendMessage(&sarama.ProducerMessage{
			Topic: consts.IncomeCreateOrders,
			Key:   sarama.StringEncoder("createOrders"),
			Value: sarama.ByteEncoder(js),
		})

		if err == nil {
			err = randService.RAND(errors.New("rand error"))
		}

		if err != nil {
			log.Printf("order %v payments: %v", order.Id, err)
			_, _, err := i.prod.SendMessage(NewRetryMes(&order, time.Now()))
			i.metrics.AmountOfRetry.Inc()
			if err != nil {
				log.Printf("order retry: %v", err)
			}
			continue
		}
		i.metrics.SuccessfulProcessed.Inc()
	}

	return nil
}

func (i *IncomeHandler) resetOperation(order *consts.Order) error {
	marshalMessage, err := json.Marshal(order)
	if err != nil {
		return err
	}

	_, _, err = i.prod.SendMessage(&sarama.ProducerMessage{
		Topic:     consts.ResetOrders,
		Key:       sarama.StringEncoder("reset"),
		Value:     sarama.ByteEncoder(marshalMessage),
		Timestamp: time.Now(),
	})
	if err != nil {
		//TODO нужно пытаться отправить пока не отправим
		return err
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
		Topic:     consts.RetrySendToCreateOrder,
		Key:       sarama.StringEncoder("retry"),
		Value:     sarama.ByteEncoder(marshalMessage),
		Timestamp: timeStamp,
	}
}
func NewAttemptSendMes(retryOrder *consts.RetryOrder) *sarama.ProducerMessage {
	order := consts.Order{
		Id:              retryOrder.Id,
		UserID:          retryOrder.UserID,
		PaymentsMethod:  retryOrder.PaymentsMethod,
		PaymentsData:    retryOrder.PaymentsData,
		ProductsId:      retryOrder.ProductsId,
		AmountOfProduct: retryOrder.AmountOfProduct,
	}
	marshalMessage, _ := json.Marshal(order)
	return &sarama.ProducerMessage{
		Topic: consts.IncomeCreateOrders,
		Key:   sarama.StringEncoder("order"),
		Value: sarama.ByteEncoder(marshalMessage),
	}
}
func (r *RetryHandler) RetryAgain(order *consts.RetryOrder, t time.Time) error {

	marshalMessage, err := json.Marshal(order)
	if err != nil {
		return err
	}

	_, _, err = r.prod.SendMessage(&sarama.ProducerMessage{
		Topic:     consts.RetrySendToCreateOrder,
		Key:       sarama.StringEncoder("retry"),
		Value:     sarama.ByteEncoder(marshalMessage),
		Timestamp: t,
	})
	if err != nil {
		return err
	}

	return nil
}

type CheckOut interface {
	ReserveFound(order *consts.Order) error
	Charge(order *consts.Order) error
}

func checkOut(order *consts.Order) error {
	switch order.PaymentsMethod {
	case consts.YellowBank:
		if err := pay(yellowBank.YellowBank{}, order); err != nil {
			return err
		}
	case consts.BlueBank:
		if err := pay(blueBank.BlueBank{}, order); err != nil {
			return err
		}
	}
	return nil
}

func pay(out CheckOut, order *consts.Order) error {
	if err := out.ReserveFound(order); err != nil {
		return err
	}
	if err := out.Charge(order); err != nil {
		return err
	}
	return nil
}
