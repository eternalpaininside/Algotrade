package users_transport_http

type UserHTTPHandler struct {
	userService UserService
}

type UserService interface {
}

func NewUserHTTPHandler(
	service UserService,
) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService: service,
	}
}
