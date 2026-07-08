package part

// service реализует бизнес-логику работы с деталями.
type service struct {
	partRepo PartRepository
}

// New создаёт новый сервис деталей.
func New(partRepo PartRepository) *service {
	return &service{
		partRepo: partRepo,
	}
}
