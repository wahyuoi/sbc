package repository

import (
	"context"

	"github.com/wahyuoi/sbc/internal/model"
)

type ExerciseRepository interface {
	Create(ctx context.Context, exercise *model.Exercise) error
	GetByUniqueKey(ctx context.Context, userID int, phraseID int, audioFormat string) (*model.Exercise, error)
}
