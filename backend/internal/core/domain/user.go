package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
	corehttputils "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/utils"
)

type User struct {
	ID        int64
	CreatedAt time.Time

	Name     string
	Email    string
	PassHash string
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
		PassHash:  pass,
	}
}

func NewUserUninitialized(
	name, email string,
) User {
	return User{
		ID:        UninitializedID,
		CreatedAt: time.Now().UTC(),

		Name:     name,
		Email:    email,
		PassHash: "",
	}
}

func (u *User) Validate() error {
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

	re := regexp.MustCompile(
		`^[^\s@]+@[^\s@]+\.[^\s@]+$`,
	)
	if !re.Match([]byte(u.Email)) {
		return fmt.Errorf(
			"invalid user email format: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

type UserPatch struct {
	Name  Nullable[string]
	Email Nullable[string]
	Pass  Nullable[string]
}

func (u *UserPatch) Validate() error {
	if u.Name.Set &&
		u.Name.Value == nil {
		return fmt.Errorf(
			"user name cant be patched to NULL: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if u.Email.Set &&
		u.Email.Value == nil {
		return fmt.Errorf(
			"user email cant be patched to NULL: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if u.Pass.Set &&
		u.Pass.Value == nil {
		return fmt.Errorf(
			"user password cant be patched to NULL: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

func (u *User) ApplyPatch(
	patch UserPatch,
) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf(
			"validate user patch: %w",
			err,
		)
	}

	tmp := *u

	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}
	if patch.Email.Set {
		tmp.Email = *patch.Email.Value
	}
	if patch.Pass.Set {
		salt, hash, err := corehttputils.PasswordHashing(*patch.Pass.Value)
		if err != nil {
			return fmt.Errorf(
				"hashing password: %w",
				err,
			)
		}
		tmp.PassHash = corehttputils.ConcatHash(salt, hash)
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf(
			"validate patched user: %w",
			err,
		)
	}

	*u = tmp

	return nil
}
