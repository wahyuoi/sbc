package repository

import (
	"context"
	"database/sql"

	"github.com/wahyuoi/sbc/internal/model"
)

var _ PhraseRepository = (*sqlPhraseRepository)(nil)

type sqlPhraseRepository struct {
	db *sql.DB
}

func NewSqlPhraseRepository(db *sql.DB) PhraseRepository {
	return &sqlPhraseRepository{db: db}
}

func (p *sqlPhraseRepository) GetById(ctx context.Context, id int) (*model.Phrase, error) {
	executor := getExecutor(ctx, p.db)

	query := "SELECT id, phrase, created_at, updated_at FROM phrases WHERE id = ?"
	row := executor.QueryRowContext(ctx, query, id)

	var phrase model.Phrase
	err := row.Scan(&phrase.ID, &phrase.Phrase, &phrase.CreatedAt, &phrase.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &phrase, nil
}
