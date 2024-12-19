package repository

import (
	"context"

	"github.com/wahyuoi/sbc/internal/model"
)

type ExerciseRepository interface {
	CreateExercise(ctx context.Context, exercise *model.Exercise) error
	GetExerciseById(ctx context.Context, userID int, phraseID int, audioFormat string) (*model.Exercise, error)
}
