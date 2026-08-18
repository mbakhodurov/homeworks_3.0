package order

// service реализует бизнес-логику заказов.
type service struct {
	orderRepo       OrderRepository
	orderItemRepo   OrderItemRepository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	txManager       TxManager
	producer        OrderPaidProducer
}

// New создаёт новый сервис заказов.
func New(
	orderRepo OrderRepository,
	orderItemRepo OrderItemRepository,
	inventoryClient InventoryClient,
	paymentClient PaymentClient,
	txManager TxManager,
	producer OrderPaidProducer,
) *service {
	return &service{
		orderRepo:       orderRepo,
		orderItemRepo:   orderItemRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		txManager:       txManager,
		producer:        producer,
	}
}
