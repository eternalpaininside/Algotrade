package domain

import "time"

type User struct {
	ID        int
	CreatedAt time.Time

	Name      string
	Email     string
	Pass      string
}
