package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type txKey struct{}

var _ UnitOfWork = (*sqlUnitOfWork)(nil)

type sqlUnitOfWork struct {
	db           *sql.DB
	userRepo     UserRepository
	phraseRepo   PhraseRepository
	exerciseRepo ExerciseRepository
}

func (s *sqlUnitOfWork) ExerciseRepository() ExerciseRepository {
	return s.exerciseRepo
}

func (s *sqlUnitOfWork) PhraseRepository() PhraseRepository {
	return s.phraseRepo
}

func (s *sqlUnitOfWork) UserRepository() UserRepository {
	return s.userRepo
}

func (s *sqlUnitOfWork) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
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

func NewSqlUnitOfWork(db *sql.DB) UnitOfWork {
	return &sqlUnitOfWork{
		db:           db,
		userRepo:     NewSqlUserRepository(db),
		phraseRepo:   NewSqlPhraseRepository(db),
		exerciseRepo: NewSqlExerciseRepository(db),
	}
}
