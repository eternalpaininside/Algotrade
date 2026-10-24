package userpostgres

import "time"

type UserModel struct {
	ID       int64
	Name     string
	Email    string
	Password string
	CreateAt time.Time
}
