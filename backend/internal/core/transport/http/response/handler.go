package corehttpresponse

import (
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	"go.uber.org/zap"
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

func (h *HTTPResponseHandler) PanicResponse(
	p any, msg string,
) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))
	h.rw.WriteHeader(statusCode)

	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	data, err := json.Marshal(response)
	if err != nil {
		h.log.Error(
			"json marshal response",
			zap.Error(err),
		)
	}

	if _, err = h.rw.Write(data); err != nil {
		h.log.Error(
			"HTTP response write",
			zap.Error(err),
		)
	}
}
