package repository

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/malkev1ch/observability/userservice/internal/model"
)

type User struct {
	lastID atomic.Int64
	m      map[int64]*model.User
}

func NewUser() *User {
	return &User{
		lastID: atomic.Int64{},
		m:      make(map[int64]*model.User),
	}
}

func (r *User) GetByID(_ context.Context, id int64) (*model.User, error) {
	user, ok := r.m[id]
	if !ok {
		return nil, fmt.Errorf("user with id %v not found", id)
	}

	return user, nil
}

func (r *User) Create(_ context.Context, user *model.User) (*model.User, error) {
	id := r.lastID.Load()
	defer r.lastID.Add(1)

	user.ID = id
	user.CreatedAt = time.Now().UTC()
	r.m[id] = user

	return user, nil
}
