package model

import "time"

type User struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

type Voucher struct {
	ID        int64
	Value     string
	UserID    int64
	UserName  string
	CreatedAt time.Time
}
