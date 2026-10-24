package userpostgres

import (
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

type UserModel struct {
	ID       int64
	Name     string
	Email    string
	Password string
	CreateAt time.Time
}

func userDomainsFromModels(
	users []UserModel,
) []domain.User {
	userDomains := make([]domain.User, 0, len(users))
	for _, user := range users {
		userDomains = append(userDomains, domain.NewUser(
			user.ID,
			user.CreateAt,
			user.Name,
			user.Email,
			user.Password,
		))
	}

	return userDomains
}
