package part

// service реализует бизнес-логику работы с деталями.
type service struct {
	txManager            TxManager
	partRepo             PartRepository
	compatibilityChecker CompatibilityChecker
}

// New создаёт новый сервис деталей.
func New(txManager TxManager, partRepo PartRepository, compatibilityChecker CompatibilityChecker) *service {
	return &service{
		txManager:            txManager,
		partRepo:             partRepo,
		compatibilityChecker: compatibilityChecker,
	}
}
