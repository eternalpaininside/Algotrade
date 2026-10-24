package usershttp

import (
	"net/http"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	corehttprequest "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/request"
	corehttpresponse "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,max=50"`
	Email string `json:"email" validate:"required,email,max=254"`
	Pass  string `json:"pass" validate:"required,min=15,max=128"`
}

type CreateUserResponse UserDTOResponse

func (u *UserHTTPHandler) CreateUser(
	rw http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := corelogger.FromContext(ctx)
	responseHandler := corehttpresponse.NewHTTPResponseHandler(
		log,
		rw,
	)

	var request CreateUserRequest
	if err := corehttprequest.DecodeAndValidateRequest(
		r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate request",
		)

		return
	}

	userDomain := domainFromDTO(request)
	userDomain, err := u.userService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create user",
		)

		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(
	dto CreateUserRequest,
) domain.User {
	return domain.NewUserUninitialized(
		dto.Name,
		dto.Email,
		dto.Pass,
	)
}
