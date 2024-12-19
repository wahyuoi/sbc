package service

import "github.com/wahyuoi/sbc/internal/repository"

type ExerciseService struct {
	uow repository.UnitOfWork
}

func NewExerciseService(uow repository.UnitOfWork) *ExerciseService {
	return &ExerciseService{uow: uow}
}
