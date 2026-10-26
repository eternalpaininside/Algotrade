package usershttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type serviceCall struct {
	operation     string
	id            int
	password      string
	user          domain.User
	patch         domain.UserPatch
	limit, offset *int
}

type userServiceStub struct {
	result domain.User
	users  []domain.User
	err    error
	calls  []serviceCall
}

func (s *userServiceStub) CreateUser(_ context.Context, password string, user domain.User) (domain.User, error) {
	s.calls = append(s.calls, serviceCall{operation: "create", password: password, user: user})
	return s.result, s.err
}
func (s *userServiceStub) GetUser(_ context.Context, id int) (domain.User, error) {
	s.calls = append(s.calls, serviceCall{operation: "get", id: id})
	return s.result, s.err
}
func (s *userServiceStub) GetUsers(_ context.Context, limit, offset *int) ([]domain.User, error) {
	s.calls = append(s.calls, serviceCall{operation: "list", limit: limit, offset: offset})
	return s.users, s.err
}
func (s *userServiceStub) PatchUser(_ context.Context, id int, patch domain.UserPatch) (domain.User, error) {
	s.calls = append(s.calls, serviceCall{operation: "patch", id: id, patch: patch})
	return s.result, s.err
}
func (s *userServiceStub) DeleteUser(_ context.Context, id int) error {
	s.calls = append(s.calls, serviceCall{operation: "delete", id: id})
	return s.err
}

type checkedRecorder struct {
	*httptest.ResponseRecorder
	statuses []int
}

func (r *checkedRecorder) WriteHeader(status int) {
	r.statuses = append(r.statuses, status)
	r.ResponseRecorder.WriteHeader(status)
}

func TestUserHTTPHandlers(t *testing.T) {
	const password = "a-long-test-password"
	const validCreate = `{"name":"Илья","email":"ilya@example.com","pass":"a-long-test-password"}`
	saved := domain.NewUser(42, time.Date(2026, 10, 26, 10, 0, 0, 0, time.UTC),
		"Илья", "ilya@example.com", "SECRET_HASH_MUST_NOT_LEAK")
	name, email, newPassword := "Новое имя", "new@example.com", "new-long-test-password"
	set := func(value *string) domain.Nullable[string] { return domain.Nullable[string]{Set: true, Value: value} }
	limit, offset := 10, 20
	type testCase struct {
		name, method, path, body string
		status                   int
		operation                string
		serviceErr               error
		errorCode                string
		wantID                   int
		wantPatch                domain.UserPatch
		wantLimit, wantOffset    *int
		emptyList                bool
	}
	tests := []testCase{
		{name: "create success", method: "POST", path: "/users", body: validCreate, status: 201, operation: "create"},
		{name: "create malformed JSON", method: "POST", path: "/users", body: "{", status: 400, errorCode: "BAD_REQUEST"},
		{name: "create missing fields", method: "POST", path: "/users", body: `{}`, status: 400, errorCode: "BAD_REQUEST"},
		{name: "create invalid email", method: "POST", path: "/users", body: `{"name":"Илья","email":"invalid","pass":"a-long-test-password"}`, status: 400, errorCode: "BAD_REQUEST"},
		{name: "create short password", method: "POST", path: "/users", body: `{"name":"Илья","email":"ilya@example.com","pass":"short"}`, status: 400, errorCode: "BAD_REQUEST"},
		{name: "create unknown field", method: "POST", path: "/users", body: strings.TrimSuffix(validCreate, "}") + `,"admin":true}`, status: 400, errorCode: "BAD_REQUEST"},
		{name: "create oversized body", method: "POST", path: "/users", body: `{"name":"` + strings.Repeat("x", 1<<20) + `"}`, status: 413, errorCode: "TOO_LARGE_REQUEST_BODY"},
		{name: "get success", method: "GET", path: "/users/42", status: 200, operation: "get", wantID: 42},
		{name: "get invalid id", method: "GET", path: "/users/abc", status: 400, errorCode: "BAD_REQUEST"},
		{name: "list success", method: "GET", path: "/users", status: 200, operation: "list"},
		{name: "list empty", method: "GET", path: "/users", status: 200, operation: "list", emptyList: true},
		{name: "list pagination", method: "GET", path: "/users?limit=10&offset=20", status: 200, operation: "list", wantLimit: &limit, wantOffset: &offset},
		{name: "list invalid limit", method: "GET", path: "/users?limit=abc", status: 400, errorCode: "BAD_REQUEST"},
		{name: "list invalid offset", method: "GET", path: "/users?offset=abc", status: 400, errorCode: "BAD_REQUEST"},
		{name: "delete success", method: "DELETE", path: "/users/42", status: 204, operation: "delete", wantID: 42},
		{name: "delete invalid id", method: "DELETE", path: "/users/abc", status: 400, errorCode: "BAD_REQUEST"},
		{name: "patch name", method: "PATCH", path: "/users/42", body: `{"name":"Новое имя"}`, status: 202, operation: "patch", wantID: 42, wantPatch: domain.UserPatch{Name: set(&name)}},
		{name: "patch email", method: "PATCH", path: "/users/42", body: `{"email":"new@example.com"}`, status: 202, operation: "patch", wantID: 42, wantPatch: domain.UserPatch{Email: set(&email)}},
		{name: "patch password", method: "PATCH", path: "/users/42", body: `{"pass":"new-long-test-password"}`, status: 202, operation: "patch", wantID: 42, wantPatch: domain.UserPatch{Pass: set(&newPassword)}},
		{name: "patch empty", method: "PATCH", path: "/users/42", body: `{}`, status: 202, operation: "patch", wantID: 42},
		{name: "patch invalid id", method: "PATCH", path: "/users/abc", body: `{}`, status: 400, errorCode: "BAD_REQUEST"},
		{name: "patch malformed JSON", method: "PATCH", path: "/users/42", body: "{", status: 400, errorCode: "BAD_REQUEST"},
		{name: "patch unknown field", method: "PATCH", path: "/users/42", body: `{"admin":true}`, status: 400, errorCode: "BAD_REQUEST"},
	}

	for _, field := range []string{"name", "email", "pass"} {
		for _, value := range []string{"null", `"   "`} {
			tests = append(tests, testCase{
				name: "patch rejects " + field + " " + value, method: "PATCH", path: "/users/42",
				body: fmt.Sprintf(`{"%s":%s}`, field, value), status: 400, errorCode: "BAD_REQUEST",
			})
		}
	}

	for _, op := range []struct{ operation, method, path, body string }{
		{"create", "POST", "/users", validCreate},
		{"get", "GET", "/users/42", ""},
		{"list", "GET", "/users", ""},
		{"patch", "PATCH", "/users/42", `{}`},
		{"delete", "DELETE", "/users/42", ""},
	} {
		for _, failure := range []struct {
			err    error
			status int
			code   string
		}{
			{coreerrors.ErrInvalidArgument, 400, "BAD_REQUEST"},
			{coreerrors.ErrNotFound, 404, "NOT_FOUND"},
			{coreerrors.ErrConflict, 409, "CONFLICT"},
			{errors.New("SECRET_DATABASE_FAILURE"), 500, "INTERNAL_ERROR"},
		} {
			id := 42
			if op.operation == "create" || op.operation == "list" {
				id = 0
			}
			tests = append(tests, testCase{
				name:   fmt.Sprintf("%s service error %d", op.operation, failure.status),
				method: op.method, path: op.path, body: op.body, status: failure.status,
				operation: op.operation, serviceErr: fmt.Errorf("wrapped: %w", failure.err),
				errorCode: failure.code, wantID: id,
			})
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &userServiceStub{result: saved, users: []domain.User{saved}, err: tt.serviceErr}
			if tt.emptyList {
				service.users = nil
			}
			if tt.wantPatch.Name.Set {
				service.result.Name = *tt.wantPatch.Name.Value
			}
			if tt.wantPatch.Email.Set {
				service.result.Email = *tt.wantPatch.Email.Value
			}
			log := &corelogger.Logger{Logger: zap.NewNop()}
			router := chi.NewRouter()
			router.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ctx := context.WithValue(r.Context(), "logger", log)
					next.ServeHTTP(w, r.WithContext(ctx))
				})
			})
			subrouter := chi.NewRouter()
			for _, route := range NewUserHTTPHandler(service).Routes() {
				subrouter.MethodFunc(route.Method, route.Path, route.Handler)
			}
			router.Mount("/api/v1", subrouter)
			request := httptest.NewRequest(tt.method, "/api/v1"+tt.path, strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "application/json")
			recorder := &checkedRecorder{ResponseRecorder: httptest.NewRecorder()}
			router.ServeHTTP(recorder, request)

			if recorder.Code != tt.status {
				t.Errorf("status: got %d, want %d; body=%s", recorder.Code, tt.status, recorder.Body.String())
			}
			if !reflect.DeepEqual(recorder.statuses, []int{tt.status}) {
				t.Errorf("WriteHeader calls: got %v, want [%d]", recorder.statuses, tt.status)
			}
			wantCalls := 0
			if tt.operation != "" {
				wantCalls = 1
			}
			if len(service.calls) != wantCalls {
				t.Errorf("service calls: got %d, want %d", len(service.calls), wantCalls)
			}
			if len(service.calls) == 1 && wantCalls == 1 {
				call := service.calls[0]
				if call.operation != tt.operation || call.id != tt.wantID {
					t.Error("incorrect operation or user ID")
				}
				if call.operation == "create" && (call.password != password || call.user.Name != "Илья" || call.user.Email != "ilya@example.com" || call.user.PassHash != "") {
					t.Error("incorrect registration arguments")
				}
				if call.operation == "patch" && !reflect.DeepEqual(call.patch, tt.wantPatch) {
					t.Error("PATCH fields did not reach the service correctly")
				}
				if call.operation == "list" && (!reflect.DeepEqual(call.limit, tt.wantLimit) || !reflect.DeepEqual(call.offset, tt.wantOffset)) {
					t.Error("incorrect pagination arguments")
				}
			}
			body := recorder.Body.Bytes()
			if strings.Contains(string(body), saved.PassHash) || strings.Contains(string(body), password) || strings.Contains(string(body), "SECRET_DATABASE_FAILURE") {
				t.Error("response leaks secret data")
			}
			if recorder.Code != tt.status {
				return
			}
			if tt.status == http.StatusNoContent {
				if len(body) != 0 {
					t.Error("204 response must have an empty body")
				}
				return
			}
			if recorder.Header().Get("Content-Type") != "application/json" {
				t.Error("missing JSON Content-Type")
			}
			if !json.Valid(body) {
				t.Fatalf("response must contain exactly one valid JSON value: %s", body)
			}
			if tt.errorCode != "" {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatal(err)
				}
				if response["error"] != tt.errorCode {
					t.Errorf("error code: got %q, want %q", response["error"], tt.errorCode)
				}
				if tt.status == 500 && response["message"] != "internal server error" {
					t.Error("500 must use a safe message")
				}
				return
			}
			want := service.result
			if tt.operation == "list" {
				var response []map[string]json.RawMessage
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatal(err)
				}
				if response == nil {
					t.Fatal("empty list must be [], not null")
				}
				if len(response) != len(service.users) {
					t.Fatal("incorrect users count")
				}
				for i, item := range response {
					assertUserJSON(t, item, service.users[i])
				}
			} else {
				var response map[string]json.RawMessage
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatal(err)
				}
				assertUserJSON(t, response, want)
			}
		})
	}
}

func assertUserJSON(t *testing.T, response map[string]json.RawMessage, want domain.User) {
	t.Helper()
	if len(response) != 4 {
		t.Errorf("expected only user_id, created_at, name and email; got %d fields", len(response))
	}
	var id int64
	var createdAt time.Time
	var name, email string
	for key, dest := range map[string]any{"user_id": &id, "created_at": &createdAt, "name": &name, "email": &email} {
		value, exists := response[key]
		if !exists {
			t.Errorf("missing %s", key)
			continue
		}
		if err := json.Unmarshal(value, dest); err != nil {
			t.Errorf("invalid %s: %v", key, err)
		}
	}
	if id != want.ID || !createdAt.Equal(want.CreatedAt) || name != want.Name || email != want.Email {
		t.Error("JSON fields differ from the service result")
	}
}
