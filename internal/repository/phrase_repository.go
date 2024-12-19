package repository

import (
	"context"

	"github.com/wahyuoi/sbc/internal/model"
)

type PhraseRepository interface {
	GetPhraseById(ctx context.Context, id int) (*model.Phrase, error)
}
