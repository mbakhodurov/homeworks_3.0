package errs

import "errors"

var (
	// ErrPartNotFound деталь не найдена
	ErrPartNotFound = errors.New("деталь не найдена")

	// ErrInvalidUUID неверный формат UUID
	ErrInvalidUUID = errors.New("неверный формат UUID")

	// ErrInvalidPartInfo неверные данные детали
	ErrInvalidPartInfo = errors.New("неверные данные детали")

	// ErrOutOfStock деталь отсутствует на складе
	ErrOutOfStock = errors.New("деталь отсутствует на складе")

	// ErrNothingToRelease нечего освобождать
	ErrNothingToRelease = errors.New("нечего освобождать")

	// ErrIncompatibleParts детали несовместимы
	ErrIncompatibleParts = errors.New("детали несовместимы")

	// ErrPartTypeMismatch тип детали не соответствует слоту корабля
	ErrPartTypeMismatch = errors.New("тип детали не соответствует слоту корабля")

	// ErrInvalidProperties некорректные свойства детали
	ErrInvalidProperties = errors.New("некорректные свойства детали")
)
