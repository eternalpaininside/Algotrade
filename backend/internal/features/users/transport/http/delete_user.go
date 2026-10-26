package usershttp

import (
	"net/http"

	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	corehttpresponse "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/response"
	corehttputils "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/utils"
)

func (u *UserHTTPHandler) DeleteUser(
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
			"failed to get userID path value: %w",
		)
		return
	}

	if err = u.userService.DeleteUser(
		ctx,
		userID,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}

	responseHandler.NoContentResponse()
}
