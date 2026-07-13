package errs

import "errors"

var (
	// ErrPartNotFound деталь не найдена
	ErrPartNotFound = errors.New("деталь не найдена")

	// ErrInvalidUUID неверный формат UUID
	ErrInvalidUUID = errors.New("неверный формат UUID")

	// ErrInvalidPartInfo неверные данные детали
	ErrInvalidPartInfo = errors.New("неверные данные детали")
)
