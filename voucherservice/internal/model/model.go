package model

import "time"

type Voucher struct {
	ID        int64
	UserID    int64
	Value     string
	CreatedAt time.Time
}
