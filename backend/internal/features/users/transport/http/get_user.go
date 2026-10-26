package usershttp

import (
	"net/http"

	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	corehttpresponse "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/response"
	corehttputils "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/utils"
)

type GetUserResponse UserDTOResponse

func (u *UserHTTPHandler) GetUser(
	rw http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := corelogger.FromContext(ctx)
	responseHandler := corehttpresponse.NewHTTPResponseHandler(
		log,
		rw,
	)

	userID, err := corehttputils.GetIntPathValues(
		r,
		UserIDKey,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	user, err := u.userService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)
		return
	}

	response := GetUserResponse(
		userDTOFromDomain(
			user,
		),
	)
	responseHandler.JSONResponse(
		response,
		http.StatusOK,
	)
}
