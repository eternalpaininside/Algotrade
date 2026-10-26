package usershttp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	corehttprequest "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/request"
	corehttpresponse "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/response"
	corehttputils "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	Name  corehttputils.Nullable[string] `json:"name"`
	Email corehttputils.Nullable[string] `json:"email"`
	Pass  corehttputils.Nullable[string] `json:"pass"`
}

func (p *PatchUserRequest) String() string {
	return fmt.Sprintf(
		"Name: %v\nEmail: %v\nPass: %v",
		p.Name,
		p.Email,
		p.Pass,
	)
}

func (p *PatchUserRequest) Validate() error {
	if p.Name.Set {
		if p.Name.Value == nil {
			return fmt.Errorf("name cannot be null")
		}

		if strings.TrimSpace(*p.Name.Value) == "" {
			return fmt.Errorf("name cannot be empty")
		}
	}

	if p.Email.Set {
		if p.Email.Value == nil {
			return fmt.Errorf("email cannot be null")
		}

		if strings.TrimSpace(*p.Email.Value) == "" {
			return fmt.Errorf("email cannot be empty")
		}
	}

	if p.Pass.Set {
		if p.Pass.Value == nil {
			return fmt.Errorf("password cannot be null")
		}

		if strings.TrimSpace(*p.Pass.Value) == "" {
			return fmt.Errorf("password cannot be empty")
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

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

	userID, err := corehttputils.GetIntPathValues(
		r,
		"id",
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user ID from path",
		)
		return
	}

	var request PatchUserRequest
	if err = corehttprequest.DecodeAndValidateRequest(
		r,
		&request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	userPatch := userPatchFromRequest(request)
	userDomain, err := u.userService.PatchUser(
		ctx,
		userID,
		userPatch,
	)

	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"service failed to patch user",
		)
		return
	}

	response := PatchUserResponse(
		userDTOFromDomain(
			userDomain,
		),
	)

	responseHandler.JSONResponse(
		&response,
		http.StatusAccepted,
	)
}

func userPatchFromRequest(
	req PatchUserRequest,
) domain.UserPatch {
	return domain.UserPatch{
		Name:  req.Name.ToDomain(),
		Email: req.Email.ToDomain(),
		Pass:  req.Pass.ToDomain(),
	}
}
