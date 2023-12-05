package repository

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/malkev1ch/observability/voucherservice/internal/model"
)

type Voucher struct {
	lastID atomic.Int64
	m      map[int64]*model.Voucher
}

func NewVoucher() *Voucher {
	return &Voucher{
		lastID: atomic.Int64{},
		m:      make(map[int64]*model.Voucher),
	}
}

func (r *Voucher) GetByUserID(_ context.Context, id int64) (*model.Voucher, error) {
	user, ok := r.m[id]
	if !ok {
		return nil, fmt.Errorf("voucher with user id %v not found", id)
	}

	return user, nil
}

func (r *Voucher) Create(_ context.Context, voucher *model.Voucher) (*model.Voucher, error) {
	id := r.lastID.Load()
	defer r.lastID.Add(1)

	voucher.ID = id
	voucher.CreatedAt = time.Now().UTC()
	r.m[voucher.UserID] = voucher

	return voucher, nil
}
