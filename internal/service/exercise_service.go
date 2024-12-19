package service

import (
	"context"

	"github.com/wahyuoi/sbc/internal/model"
	"github.com/wahyuoi/sbc/internal/repository"
)

type ExerciseService struct {
	uow repository.UnitOfWork
}

func NewExerciseService(uow repository.UnitOfWork) ExerciseService {
	return ExerciseService{uow: uow}
}

func (s *ExerciseService) SubmitAudio(ctx context.Context, userID, phraseID int, audioPath, audioFormat string) error {
	return s.uow.WithTransaction(ctx, func(ctx context.Context) error {
		// Check if user and phrase exists.
		// It is being checked here to make sure that both values are valid until the exercise saved.
		_, err := s.uow.UserRepository().GetById(ctx, userID)
		if err != nil {
			return err
		}
		_, err = s.uow.PhraseRepository().GetById(ctx, phraseID)
		if err != nil {
			return err
		}

		return s.uow.ExerciseRepository().Create(ctx, &model.Exercise{
			ID:          userID,
			PhraseID:    phraseID,
			UserID:      userID,
			AudioPath:   audioPath,
			AudioFormat: audioFormat,
		})
	})
}
