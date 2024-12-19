package model

import "time"

type Phrase struct {
	ID        int       `json:"id"`
	Phrase    string    `json:"phrase"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Exercise struct {
	ID          int       `json:"id"`
	PhraseID    int       `json:"phrase_id"`
	UserID      int       `json:"user_id"`
	AudioPath   string    `json:"audio_path"`
	AudioFormat string    `json:"audio_format"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
