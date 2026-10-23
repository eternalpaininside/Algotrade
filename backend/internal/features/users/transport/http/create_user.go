package usershttp

import (
	"net/http"
	"time"

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

type CreateUserResponse struct {
	ID        int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	Name  string `json:"name"`
	Email string `json:"email"`
}

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
	}

	response := dtoFromDomain(userDomain)

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

func dtoFromDomain(
	user domain.User,
) CreateUserResponse {
	return CreateUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		Name:      user.Name,
		Email:     user.Email,
	}
}
