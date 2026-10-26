package corehttpresponse

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"

	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	"go.uber.org/zap"
)

var (
	ErrBadRequest      = errors.New("BAD_REQUEST")
	ErrNotFound        = errors.New("NOT_FOUND")
	ErrConflict        = errors.New("CONFLICT")
	ErrTooLargeRequest = errors.New("TOO_LARGE_REQUEST_BODY")
	ErrValidation      = errors.New("VALIDATION_ERROR")
	ErrInternal        = errors.New("INTERNAL_ERROR")
)

type HTTPResponseHandler struct {
	log *corelogger.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(
	logger *corelogger.Logger,
	w http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: logger,
		rw:  w,
	}
}

func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	data, err := json.Marshal(responseBody)
	if err != nil {
		h.log.Error(
			"json marshal response body",
			zap.Error(err),
		)
		http.Error(
			h.rw,
			"failed to serializing response body",
			http.StatusInternalServerError,
		)
		return
	}

	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(statusCode)

	if _, err = h.rw.Write(data); err != nil {
		h.log.Error(
			"response write data",
			zap.Error(err),
		)
		return
	}

}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

func (h *HTTPResponseHandler) PanicResponse(
	p any, msg string,
) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(
		statusCode,
		err,
		msg,
	)
}

func (h *HTTPResponseHandler) ErrorResponse(
	err error,
	msg string,
) {
	var (
		statusCode  int
		logFunc     func(string, ...zap.Field)
		responseErr error
	)

	switch {
	case errors.Is(err, coreerrors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
		responseErr = ErrBadRequest

	case errors.Is(err, coreerrors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug
		responseErr = ErrNotFound

	case errors.Is(err, coreerrors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn
		responseErr = ErrConflict

	case errors.Is(err, coreerrors.ErrTooLarge):
		statusCode = http.StatusRequestEntityTooLarge
		logFunc = h.log.Warn
		responseErr = ErrTooLargeRequest

	case errors.Is(err, coreerrors.ErrValidation):
		statusCode = http.StatusUnprocessableEntity
		logFunc = h.log.Warn
		responseErr = ErrValidation

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
		responseErr = ErrInternal
	}

	logFunc(msg, zap.Error(err))

	if statusCode >= 500 {
		msg = "internal server error"
	}

	h.errorResponse(
		statusCode,
		responseErr,
		msg,
	)
}

func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {
	h.JSONResponse(
		map[string]string{
			"message": msg,
			"error":   err.Error(),
		},
		statusCode,
	)
}
