package repository

import (
	"context"
	"database/sql"

	"github.com/wahyuoi/sbc/internal/model"
)

var _ ExerciseRepository = (*sqlExerciseRepository)(nil)

type sqlExerciseRepository struct {
	db *sql.DB
}

func NewSqlExerciseRepository(db *sql.DB) ExerciseRepository {
	return &sqlExerciseRepository{db: db}
}

func (e *sqlExerciseRepository) CreateExercise(ctx context.Context, exercise *model.Exercise) error {
	executor := getExecutor(ctx, e.db)
	query := "INSERT INTO exercises (phrase_id, user_id, audio_path, audio_format, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())"
	_, err := executor.ExecContext(ctx, query, exercise.PhraseID, exercise.UserID, exercise.AudioPath, exercise.AudioFormat)
	return err
}

func (e *sqlExerciseRepository) GetExerciseById(ctx context.Context, userID int, phraseID int, audioFormat string) (*model.Exercise, error) {
	executor := getExecutor(ctx, e.db)
	query := "SELECT id, phrase_id, user_id, audio_path, audio_format, created_at, updated_at FROM exercises WHERE user_id = ? AND phrase_id = ? AND audio_format = ?"
	row := executor.QueryRowContext(ctx, query, userID, phraseID, audioFormat)

	var exercise model.Exercise
	err := row.Scan(&exercise.ID, &exercise.PhraseID, &exercise.UserID, &exercise.AudioPath, &exercise.AudioFormat, &exercise.CreatedAt, &exercise.UpdatedAt)
	return &exercise, err
}
