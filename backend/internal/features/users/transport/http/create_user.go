package users_transport_http

import (
	jsonv2 "encoding/json/v2"
	"net/http"
	"time"
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
	var request CreateUserRequest

	r.Body = http.MaxBytesReader(rw, r.Body, 1<<20)
	if err := jsonv2.UnmarshalRead(
		r.Body,
		&request,
		jsonv2.RejectUnknownMembers(true),
	); err != nil {
		http.Error(rw)
	}

}
