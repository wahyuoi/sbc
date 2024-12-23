package repository

import (
	"context"

	"github.com/wahyuoi/sbc/internal/model"
)

type ExerciseRepository interface {
	// Upsert is used to create or update an exercise.
	//
	// WARN: This method will override existing audio path if the exercise already exists.
	// But it will not delete the existing audio file. Probably needs to check with the team,
	// whether it is safe to delete the existing audio file, or we need to keep it for audit purposes.
	// If we need to keep the history of the exercise, we can remove the UniqueKey constraint in the database,
	// and add a version column.
	Upsert(ctx context.Context, exercise *model.Exercise) error

	// GetByUniqueKey is used to get an exercise by unique key.
	// The unique key is a combination of userID, phraseID, and audioFormat.
	// Because it uses unique constraint in the mysql, it will automatically has index on the database.
	// So if the unique key is removed, please create new index on the database.
	GetByUniqueKey(ctx context.Context, userID int, phraseID int, audioFormat string) (*model.Exercise, error)
}
