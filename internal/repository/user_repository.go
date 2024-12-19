package repository

import (
	"context"

	"github.com/wahyuoi/sbc/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetById(ctx context.Context, id int) (*model.User, error)
}
