package model

import "time"

type User struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}
