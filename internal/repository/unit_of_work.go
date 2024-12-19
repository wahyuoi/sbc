package repository

import "context"

type UnitOfWork interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	UserRepository() UserRepository
	PhraseRepository() PhraseRepository
	ExerciseRepository() ExerciseRepository
}
