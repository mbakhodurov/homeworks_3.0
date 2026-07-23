package errs

import "errors"

var (
	// Ошибки заказов
	ErrOrderNotFound    = errors.New("заказ не найден")
	ErrOrderAlreadyPaid = errors.New("заказ уже оплачен")
	ErrOrderCancelled   = errors.New("заказ отменён")

	// Ошибки деталей
	ErrPartNotFound = errors.New("деталь не найдена")
	ErrOutOfStock   = errors.New("деталь отсутствует на складе")

	// ErrIncompatibleParts детали несовместимы
	ErrIncompatibleParts = errors.New("детали несовместимы")

	// ErrPartTypeMismatch тип детали не соответствует слоту корабля
	ErrPartTypeMismatch = errors.New("тип детали не соответствует слоту корабля")
)
