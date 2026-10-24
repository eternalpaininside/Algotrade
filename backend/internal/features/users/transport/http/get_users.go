package usershttp

import (
	"bytes"
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	corehttpresponse "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/response"
	corehttputils "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/utils"
)

const (
	LimitKey  = "limit"
	OffsetKey = "offset"
)

type GetUsersResponse []UserDTOResponse

func (r *GetUsersResponse) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')

	for user := range *r {
		if buf.Len() != 1 {
			buf.WriteByte(',')
		}

		data, err := json.Marshal(user)
		if err != nil{
			return nil, fmt.Errorf(
				"marshal user: %v: %w",
				user,
				err,
				)
		}

		buf.Write(data)
	}

	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func (u *UserHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corelogger.FromContext(ctx)
	responseHandler := corehttpresponse.NewHTTPResponseHandler(
		log, rw,
	)

	limit, offset, err := getPaginationQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			fmt.Sprintf("failed to get '%s'/'%s' query params",
				LimitKey,
				OffsetKey,
			),
		)
		return
	}

	userDomains, err := u.userService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get users",
		)
		return
	}

	usersDTO := GetUsersResponse(usersDTOFromDomains(userDomains))

	responseHandler.JSONResponse(usersDTO, http.StatusOK)
}

func getPaginationQueryParams(
	r *http.Request,
) (limit *int, offset *int, err error) {
	if limit, err = corehttputils.GetIntQueryParams(
		r, LimitKey,
	); err != nil {
		return nil, nil, fmt.Errorf(
			"get '%s' query param: %w",
			LimitKey,
			err,
		)
	}

	if offset, err = corehttputils.GetIntQueryParams(
		r, OffsetKey); err != nil {
		return nil, nil, fmt.Errorf(
			"get '%s' query param: %w",
			OffsetKey,
			err,
		)
	}

	return
}

func usersDTOFromDomains(
	users []domain.User,
) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, 0, len(users))
	for _, user := range users {
		usersDTO = append(usersDTO, userDTOFromDomain(user))
	}

	return usersDTO
}