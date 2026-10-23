package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
)

type User struct {
	ID        int64
	CreatedAt time.Time

	Name  string
	Email string
	Pass  string
}

func NewUser(
	id int64,
	createdAt time.Time,
	name, email, pass string,
) User {
	return User{
		ID:        id,
		CreatedAt: createdAt,
		Name:      name,
		Email:     email,
		Pass:      pass,
	}
}

func NewUserUninitialized(
	name, email, pass string,
) User {
	return User{
		ID:        UninitializedID,
		CreatedAt: time.Now().UTC(),

		Name:  name,
		Email: email,
		Pass:  pass,
	}
}

func (u User) Validate() error {
	if u.ID <= 0 {
		return fmt.Errorf(
			"uninitialized user id: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	u.Name = strings.TrimSpace(u.Name)
	if u.Name == "" {
		return fmt.Errorf(
			"empty user name: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	u.Email = strings.TrimSpace(u.Email)
	if u.Email == "" {
		return fmt.Errorf(
			"empty user email: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	re := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !re.Match([]byte(u.Email)) {
		return fmt.Errorf(
			"invalid user email format: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}
