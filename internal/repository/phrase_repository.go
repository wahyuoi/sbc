package repository

import (
	"context"

	"github.com/wahyuoi/sbc/internal/model"
)

type PhraseRepository interface {
	GetById(ctx context.Context, id int) (*model.Phrase, error)
}
