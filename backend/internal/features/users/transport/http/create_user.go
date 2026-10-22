package usershttp

import (
	jsonv2 "encoding/json/v2"
	"net/http"
	"time"

	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	"go.uber.org/zap"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

type CreateUserResponse struct {
	ID        int       `json:"user_id"`
	CreatedAT time.Time `json:"created_at"`

	Name  string `json:"name"`
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

func (u *UserHTTPHandler) CreateUser(
	rw http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := corelogger.FromContext(ctx)

	log.Debug("start create user handler")

	var request CreateUserRequest

	r.Body = http.MaxBytesReader(rw, r.Body, 1<<20)
	if err := jsonv2.UnmarshalRead(
		r.Body,
		&request,
		jsonv2.RejectUnknownMembers(true),
	); err != nil {
		log.Error(
			"json unmarshal read",
			zap.Error(err))
		http.Error(
			rw,
			"unexpected errors while decoding",
			http.StatusBadRequest,
		)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
