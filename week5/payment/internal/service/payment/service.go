package payment

type service struct{}

// New создаёт сервис бизнес-логики оплаты.
func New() *service {
	return &service{}
}
