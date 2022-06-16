package _const

var (
	Brokers = []string{"localhost:9092"} // for local testing
	//Brokers = []string{os.Getenv("kafkaURL")} //for docker
)

const (
	IncomeOrders            = "income_orders"
	IncomePayments          = "income_payments"
	IncomeCreateOrders      = "income_create_orders"
	IncomeNotification      = "income_notification"
	ResetOrders             = "reset_orders"
	RetrySendToCreateOrder  = "retry_send_to_create_order"
	RetrySendToNotification = "retry_send_to_notification"
	RetryNotification       = "retry_notification"
)

var Topics = []string{IncomeOrders, IncomePayments, IncomeCreateOrders, IncomeNotification, ResetOrders, RetryNotification, RetrySendToNotification, RetrySendToCreateOrder}

const (
	YellowBank = "yellowBank"
	BlueBank   = "blueBank"
)

type Product struct {
	ProductsId      int64 `json:"products_id"`
	AmountOfProduct int64 `json:"amount_of_product"`
}

type Order struct {
	Id              int64  `json:"id"`
	UserID          int64  `json:"user_id"`
	PaymentsMethod  string `json:"payment_method"`
	PaymentsData    string `json:"payments_data"`
	ProductsId      int64  `json:"products_id"`
	AmountOfProduct int64  `json:"amount_of_product"`
}

type RetryOrder struct {
	Id              int64  `json:"id"`
	UserID          int64  `json:"user_id"`
	PaymentsMethod  string `json:"payment_method"`
	PaymentsData    string `json:"payments_data"`
	ProductsId      int64  `json:"products_id"`
	AmountOfProduct int64  `json:"amount_of_product"`
	AmountOfRetry   int64  `json:"amount_of_retry"`
	Delay           int64  `json:"delay"`
}
