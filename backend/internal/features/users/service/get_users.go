package usersservice

import (
	"context"
	"fmt"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
)

func (u *UserService) GetUsers(
	ctx context.Context,
	limit, offset *int,
) ([]domain.User, error) {
	if limit != nil &&
		*limit < 0 {
		return nil, fmt.Errorf(
			"limit(%d) must be non-negative: %w",
			*limit,
			coreerrors.ErrInvalidArgument,
		)
	}

	if offset != nil &&
		*offset < 0 {
		return nil, fmt.Errorf(
			"offset(%d) must be non-negative: %w",
			*offset,
			coreerrors.ErrInvalidArgument,
		)
	}

	userDomains, err := u.repository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf(
			"get users from repository: %w",
			err,
		)
	}
}
