package v1

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	errs "github.com/mbakhodurov/homeworks2/week5/order/internal/errors"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler — глобальный hook ogen. Подключается через orderv1.WithErrorHandler.
func ErrorHandler(ctx context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	code, message := mapError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if encErr := json.NewEncoder(w).Encode(errorResponse{Code: code, Message: message}); encErr != nil {
		slog.ErrorContext(ctx, "ошибка кодирования ответа", "error", encErr)
	}
}

func mapError(err error) (int, string) {
	var decodeParams *ogenerrors.DecodeParamsError
	var decodeRequest *ogenerrors.DecodeRequestError

	switch {
	case errors.As(err, &decodeParams), errors.As(err, &decodeRequest):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, errs.ErrPartTypeMismatch):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, errs.ErrOrderNotFound), errors.Is(err, errs.ErrPartNotFound):
		return http.StatusNotFound, err.Error()

	case errors.Is(err, errs.ErrOrderAlreadyPaid),
		errors.Is(err, errs.ErrOrderCancelled),
		errors.Is(err, errs.ErrOutOfStock),
		errors.Is(err, errs.ErrIncompatibleParts):
		return http.StatusConflict, err.Error()

	default:
		return http.StatusInternalServerError, "внутренняя ошибка"
	}
}
