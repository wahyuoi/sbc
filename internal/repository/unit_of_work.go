package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type UnitOfWork interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	UserRepository() UserRepository
	PhraseRepository() PhraseRepository
	ExerciseRepository() ExerciseRepository
	FileRepository() FileRepository
}

type txKey struct{}

var _ UnitOfWork = (*unitOfWork)(nil)

type unitOfWork struct {
	db             *sql.DB
	userRepo       UserRepository
	phraseRepo     PhraseRepository
	exerciseRepo   ExerciseRepository
	fileRepository FileRepository
}

func (s *unitOfWork) ExerciseRepository() ExerciseRepository {
	return s.exerciseRepo
}

func (s *unitOfWork) PhraseRepository() PhraseRepository {
	return s.phraseRepo
}

func (s *unitOfWork) UserRepository() UserRepository {
	return s.userRepo
}

func (s *unitOfWork) FileRepository() FileRepository {
	return s.fileRepository
}

func (s *unitOfWork) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	defer func() {
		if p := recover(); p != nil {
			// rollback the transaction if panic
			tx.Rollback()
		}
	}()

	if err = fn(txCtx); err != nil {
		// rollback the transaction if error
		if rberr := tx.Rollback(); rberr != nil {
			return fmt.Errorf("failed to rollback transaction: %w, original error: %w", rberr, err)
		}
		return err
	}

	return tx.Commit()
}

func NewUnitOfWork(db *sql.DB, fileRepository FileRepository) UnitOfWork {
	return &unitOfWork{
		db:             db,
		userRepo:       NewSqlUserRepository(db),
		phraseRepo:     NewSqlPhraseRepository(db),
		exerciseRepo:   NewSqlExerciseRepository(db),
		fileRepository: fileRepository,
	}
}
