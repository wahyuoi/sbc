package service

import (
	"context"

	"github.com/wahyuoi/sbc/internal/model"
	"github.com/wahyuoi/sbc/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	uow repository.UnitOfWork
}

func NewUserService(uow repository.UnitOfWork) *UserService {
	return &UserService{uow: uow}
}

func (s *UserService) Register(ctx context.Context, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.uow.WithTransaction(ctx, func(ctx context.Context) error {
		return s.uow.UserRepository().Create(ctx, user)
	})
}

func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.uow.UserRepository().GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.uow.UserRepository().GetByEmail(ctx, email)
}

// TODO: if we implement the DeleteUser, we need to make sure that it also delete the exercises record related to the user.
// In this case, transaction using unit of work will be useful.
// Might also need to delete the audio file related to the exercises.
