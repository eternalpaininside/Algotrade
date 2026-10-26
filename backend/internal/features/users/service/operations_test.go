package usersservice

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
	"golang.org/x/crypto/argon2"
)

type createUserRepositoryMock struct {
	create func(context.Context, domain.User) (domain.User, error)
}

func (s *createUserRepositoryMock) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	return s.create(ctx, user)
}

func (*createUserRepositoryMock) GetUsers(
	context.Context, *int, *int,
) ([]domain.User, error) {
	panic("unexpected GetUsers call")
}

func (*createUserRepositoryMock) GetUser(
	context.Context, int,
) (domain.User, error) {
	panic("unexpected GetUser call")
}

func (*createUserRepositoryMock) DeleteUser(
	context.Context, int,
) error {
	panic("unexpected DeleteUser call")
}

func (*createUserRepositoryMock) PatchUser(
	context.Context, int, domain.User,
) (domain.User, error) {
	panic("unexpected PatchUser call")
}

func TestCreateUser_Success(t *testing.T) {
	const password = "a-long-test-password"

	ctx := context.Background()
	createdAt := time.Date(2026, 10, 26, 10, 0, 0, 0, time.UTC)
	calls := 0

	repository := &createUserRepositoryMock{
		create: func(
			gotCtx context.Context,
			user domain.User,
		) (domain.User, error) {
			calls++

			if gotCtx != ctx {
				t.Fatal("service did not pass the request context")
			}

			if user.Name != "Илья" || user.Email != "ilya@example.com" {
				t.Fatalf(
					"unexpected normalized fields: name=%q email=%q",
					user.Name,
					user.Email,
				)
			}

			assertPasswordHash(t, password, user.PassHash)

			user.ID = 42
			user.CreatedAt = createdAt
			return user, nil
		},
	}

	service := NewUserService(repository)

	user, err := service.CreateUser(
		ctx,
		password,
		domain.NewUserUninitialized("  Илья  ", "  ilya@example.com  "),
	)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if calls != 1 {
		t.Fatalf("repository calls: got %d, want 1", calls)
	}

	if user.ID != 42 || !user.CreatedAt.Equal(createdAt) {
		t.Fatalf(
			"service did not return saved fields: id=%d created_at=%v",
			user.ID,
			user.CreatedAt,
		)
	}
}

func TestCreateUser_InvalidUser(t *testing.T) {
	tests := []struct {
		name  string
		user  string
		email string
	}{
		{"empty name", "", "ilya@example.com"},
		{"spaces-only name", "   ", "ilya@example.com"},
		{"tabs-only name", "\t\t", "ilya@example.com"},
		{"newlines-only name", "\n\r\n", "ilya@example.com"},
		{"mixed whitespace name", " \t\n\r ", "ilya@example.com"},
		{"unicode whitespace name", "\u00a0\u2003", "ilya@example.com"},

		{"empty email", "Илья", ""},
		{"spaces-only email", "Илья", "   "},
		{"whitespace-only email", "Илья", "\t\n"},

		{"missing at", "Илья", "ilya.example.com"},
		{"missing local part", "Илья", "@example.com"},
		{"missing domain", "Илья", "ilya@"},
		{"missing domain dot", "Илья", "ilya@example"},
		{"missing domain suffix", "Илья", "ilya@example."},
		{"multiple at signs", "Илья", "ilya@@example.com"},

		{"space in local part", "Илья", "il ya@example.com"},
		{"space in domain", "Илья", "ilya@exam ple.com"},
		{"tab in local part", "Илья", "il\tya@example.com"},
		{"newline in domain", "Илья", "ilya@exam\nple.com"},

		{"both fields empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &createUserRepositoryMock{
				create: func(
					context.Context,
					domain.User,
				) (domain.User, error) {
					t.Fatal("repository must not be called for invalid input")
					return domain.User{}, nil
				},
			}

			service := NewUserService(repository)

			_, err := service.CreateUser(
				context.Background(),
				"a-long-test-password",
				domain.NewUserUninitialized(tt.user, tt.email),
			)

			if !errors.Is(err, coreerrors.ErrInvalidArgument) {
				t.Fatalf("got %v, want ErrInvalidArgument", err)
			}
		})
	}
}

func TestCreateUser_RepositoryError(t *testing.T) {
	for _, wantErr := range []error{
		coreerrors.ErrConflict,
		errors.New("database unavailable"),
	} {
		t.Run(wantErr.Error(), func(t *testing.T) {
			repository := &createUserRepositoryMock{
				create: func(
					context.Context,
					domain.User,
				) (domain.User, error) {
					return domain.User{}, wantErr
				},
			}

			service := NewUserService(repository)

			_, err := service.CreateUser(
				context.Background(),
				"a-long-test-password",
				domain.NewUserUninitialized("Илья", "ilya@example.com"),
			)

			if !errors.Is(err, wantErr) {
				t.Fatalf("got %v, want wrapped %v", err, wantErr)
			}
		})
	}
}

func assertPasswordHash(t *testing.T, password, encoded string) {
	t.Helper()

	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" {
		t.Fatal("invalid Argon2id record format")
	}

	if parts[2] != fmt.Sprintf("v=%d", argon2.Version) {
		t.Fatal("unexpected Argon2 version")
	}

	wantParams := fmt.Sprintf(
		"m=%d,t=%d,p=%d",
		domain.PassMemory,
		domain.PassTime,
		domain.PassThreads,
	)
	if parts[3] != wantParams {
		t.Fatalf("unexpected Argon2 parameters: %q", parts[3])
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) != domain.PassSaltSize {
		t.Fatal("invalid salt")
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(storedHash) != domain.PassKeySize {
		t.Fatal("invalid hash")
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		domain.PassTime,
		domain.PassMemory,
		domain.PassThreads,
		domain.PassKeySize,
	)

	if subtle.ConstantTimeCompare(actualHash, storedHash) != 1 {
		t.Fatal("stored hash does not match the input password")
	}
}

type repositoryStub struct {
	get    func(context.Context, int) (domain.User, error)
	list   func(context.Context, *int, *int) ([]domain.User, error)
	patch  func(context.Context, int, domain.User) (domain.User, error)
	delete func(context.Context, int) error
}

func (*repositoryStub) CreateUser(context.Context, domain.User) (domain.User, error) {
	panic("unexpected CreateUser call")
}
func (s *repositoryStub) GetUser(ctx context.Context, id int) (domain.User, error) {
	return s.get(ctx, id)
}
func (s *repositoryStub) GetUsers(ctx context.Context, limit, offset *int) ([]domain.User, error) {
	return s.list(ctx, limit, offset)
}
func (s *repositoryStub) PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error) {
	return s.patch(ctx, id, user)
}
func (s *repositoryStub) DeleteUser(ctx context.Context, id int) error {
	return s.delete(ctx, id)
}

func serviceTestUser() domain.User {
	return domain.NewUser(42, time.Date(2026, 10, 26, 10, 0, 0, 0, time.UTC),
		"Илья", "ilya@example.com", "existing-password-hash")
}

func TestGetUser(t *testing.T) {
	for _, wantErr := range []error{nil, coreerrors.ErrNotFound, errors.New("storage unavailable")} {
		name := "success"
		if wantErr != nil {
			name = wantErr.Error()
		}
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			want := serviceTestUser()
			calls := 0
			repo := &repositoryStub{get: func(gotCtx context.Context, id int) (domain.User, error) {
				calls++
				if gotCtx != ctx || id != 42 {
					t.Fatal("incorrect repository arguments")
				}
				return want, wantErr
			}}
			got, err := NewUserService(repo).GetUser(ctx, 42)
			if !errors.Is(err, wantErr) {
				t.Fatalf("got error %v, want %v", err, wantErr)
			}
			if calls != 1 {
				t.Fatalf("repository calls: %d", calls)
			}
			if err == nil && got != want {
				t.Fatal("returned user differs from repository result")
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	for _, wantErr := range []error{nil, coreerrors.ErrNotFound, errors.New("storage unavailable")} {
		name := "success"
		if wantErr != nil {
			name = wantErr.Error()
		}
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			calls := 0
			repo := &repositoryStub{delete: func(gotCtx context.Context, id int) error {
				calls++
				if gotCtx != ctx || id != 42 {
					t.Fatal("incorrect repository arguments")
				}
				return wantErr
			}}
			err := NewUserService(repo).DeleteUser(ctx, 42)
			if !errors.Is(err, wantErr) {
				t.Fatalf("got error %v, want %v", err, wantErr)
			}
			if calls != 1 {
				t.Fatalf("repository calls: %d", calls)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	limit, offset, zero, negative := 10, 20, 0, -1
	storageErr := errors.New("storage unavailable")
	for _, tt := range []struct {
		name             string
		limit, offset    *int
		result           []domain.User
		repoErr, wantErr error
		wantCalls        int
	}{
		{"without pagination", nil, nil, []domain.User{serviceTestUser()}, nil, nil, 1},
		{"with pagination", &limit, &offset, []domain.User{serviceTestUser()}, nil, nil, 1},
		{"zero values", &zero, &zero, nil, nil, nil, 1},
		{"empty result", nil, nil, []domain.User{}, nil, nil, 1},
		{"negative limit", &negative, nil, nil, nil, coreerrors.ErrInvalidArgument, 0},
		{"negative offset", nil, &negative, nil, nil, coreerrors.ErrInvalidArgument, 0},
		{"repository error", nil, nil, nil, storageErr, storageErr, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			calls := 0
			repo := &repositoryStub{list: func(gotCtx context.Context, gotLimit, gotOffset *int) ([]domain.User, error) {
				calls++
				if gotCtx != ctx || !reflect.DeepEqual(gotLimit, tt.limit) || !reflect.DeepEqual(gotOffset, tt.offset) {
					t.Fatal("incorrect repository arguments")
				}
				return tt.result, tt.repoErr
			}}
			got, err := NewUserService(repo).GetUsers(ctx, tt.limit, tt.offset)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
			if calls != tt.wantCalls {
				t.Fatalf("repository calls: got %d, want %d", calls, tt.wantCalls)
			}
			if err == nil && !reflect.DeepEqual(got, tt.result) {
				t.Fatal("unexpected users")
			}
		})
	}
}

func TestPatchUser(t *testing.T) {
	name, email, password, blank, invalidEmail := "Новое имя", "new@example.com", "new-long-test-password", "   ", "invalid"
	set := func(value *string) domain.Nullable[string] { return domain.Nullable[string]{Set: true, Value: value} }
	storageErr := errors.New("storage unavailable")
	for _, tt := range []struct {
		name                      string
		patch                     domain.UserPatch
		getErr, patchErr, wantErr error
		wantPatchCalls            int
	}{
		{"name only", domain.UserPatch{Name: set(&name)}, nil, nil, nil, 1},
		{"email only", domain.UserPatch{Email: set(&email)}, nil, nil, nil, 1},
		{"password only", domain.UserPatch{Pass: set(&password)}, nil, nil, nil, 1},
		{"empty patch", domain.UserPatch{}, nil, nil, nil, 1},
		{"null name", domain.UserPatch{Name: set(nil)}, nil, nil, coreerrors.ErrInvalidArgument, 0},
		{"null email", domain.UserPatch{Email: set(nil)}, nil, nil, coreerrors.ErrInvalidArgument, 0},
		{"null password", domain.UserPatch{Pass: set(nil)}, nil, nil, coreerrors.ErrInvalidArgument, 0},
		{"blank name", domain.UserPatch{Name: set(&blank)}, nil, nil, coreerrors.ErrInvalidArgument, 0},
		{"invalid email", domain.UserPatch{Email: set(&invalidEmail)}, nil, nil, coreerrors.ErrInvalidArgument, 0},
		{"not found", domain.UserPatch{}, coreerrors.ErrNotFound, nil, coreerrors.ErrNotFound, 0},
		{"read error", domain.UserPatch{}, storageErr, nil, storageErr, 0},
		{"conflict", domain.UserPatch{Name: set(&name)}, nil, coreerrors.ErrConflict, coreerrors.ErrConflict, 1},
		{"write error", domain.UserPatch{}, nil, storageErr, storageErr, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			original := serviceTestUser()
			getCalls, patchCalls := 0, 0
			saved := original
			repo := &repositoryStub{
				get: func(gotCtx context.Context, id int) (domain.User, error) {
					getCalls++
					if gotCtx != ctx || id != 42 {
						t.Fatal("incorrect GetUser arguments")
					}
					return original, tt.getErr
				},
				patch: func(gotCtx context.Context, id int, user domain.User) (domain.User, error) {
					patchCalls++
					if gotCtx != ctx || id != 42 {
						t.Fatal("incorrect PatchUser arguments")
					}
					wantName, wantEmail := original.Name, original.Email
					if tt.patch.Name.Set {
						wantName = *tt.patch.Name.Value
					}
					if tt.patch.Email.Set {
						wantEmail = *tt.patch.Email.Value
					}
					if user.Name != wantName || user.Email != wantEmail || user.ID != original.ID || !user.CreatedAt.Equal(original.CreatedAt) {
						t.Fatal("patch changed incorrect user fields")
					}
					if tt.patch.Pass.Set {
						assertPasswordHash(t, password, user.PassHash)
					} else if user.PassHash != original.PassHash {
						t.Fatal("patch changed an omitted password")
					}
					saved = user
					return saved, tt.patchErr
				},
			}
			got, err := NewUserService(repo).PatchUser(ctx, 42, tt.patch)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
			if getCalls != 1 || patchCalls != tt.wantPatchCalls {
				t.Fatalf("calls: get=%d patch=%d, want get=1 patch=%d", getCalls, patchCalls, tt.wantPatchCalls)
			}
			if err == nil && got != saved {
				t.Fatal("service did not return saved user")
			}
		})
	}
}
