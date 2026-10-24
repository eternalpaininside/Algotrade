package corehttputils

import (
	"fmt"
	"net/http"
	"strconv"

	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
)

func GetIntQueryParams(
	r *http.Request,
	key string,
) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s', by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	return &val, nil
}

func GetIntPathValues(
	r *http.Request,
	key string,
) (int, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf(
			"no key %s int path values: %w",
			key,
			coreerrors.ErrInvalidArgument,
		)
	}

	val, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf(
			"path value=%s by key=%s not a valid integer: %v: %w",
			pathValue,
			key,
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	return val, nil
}
