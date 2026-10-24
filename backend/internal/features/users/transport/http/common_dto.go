package usershttp

import (
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

type UserDTOResponse struct {
	ID        int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	Name  string `json:"name"`
	Email string `json:"email"`
}

func userDTOFromDomain(
	user domain.User,
) UserDTOResponse {
	return UserDTOResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		Name:      user.Name,
		Email:     user.Email,
	}
}
