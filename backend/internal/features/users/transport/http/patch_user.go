package usershttp

import (
	"fmt"
	"net/http"

	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	corehttprequest "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/request"
	corehttpresponse "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/response"
)

type PatchUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Pass  *string `json:"pass"`
}

func (p PatchUserRequest) String() string {
	return fmt.Sprintf(
		"Name: %v\nEmail: %v\nPass: %v",
		p.Name,
		p.Email,
		p.Pass,
	)
}

func (u *UserHTTPHandler) PatchUser(
	rw http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := corelogger.FromContext(ctx)
	responseHandler := corehttpresponse.NewHTTPResponseHandler(
		log,
		rw,
	)

	var request PatchUserRequest
	if err := corehttprequest.DecodeAndValidateRequest(
		r,
		&request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
	}

	log.Debug(
		fmt.Sprintf(
			"PatchUserRequest fields:\n%s",
			request.String(),
		),
	)
}
