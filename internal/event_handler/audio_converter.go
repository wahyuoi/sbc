package event_handler

import (
	"context"

	"github.com/wahyuoi/sbc/internal/model"
	"github.com/wahyuoi/sbc/internal/service"
)

type AudioConverter struct {
	audioService    *service.AudioService
	exerciseService *service.ExerciseService
}

func NewAudioConverter(
	audioService *service.AudioService,
	exerciseService *service.ExerciseService,
) *AudioConverter {
	return &AudioConverter{
		audioService:    audioService,
		exerciseService: exerciseService,
	}
}

// ConvertAudio converts the audio file to the new format and stores it in the file storage.
func (h *AudioConverter) ConvertAudio(ctx context.Context, userID int, phraseID int, originalAudioFormat, newAudioFormat model.AudioFormatType) error {
	originalAudioBytes, err := h.exerciseService.GetAudio(ctx, userID, phraseID, originalAudioFormat)
	if err != nil {
		return err
	}
	newAudioBytes, err := h.audioService.ConvertAudio(ctx, originalAudioBytes, newAudioFormat)
	if err != nil {
		return err
	}
	err = h.exerciseService.SubmitAudio(ctx, userID, phraseID, newAudioBytes, newAudioFormat)
	if err != nil {
		return err
	}
	return nil
}
