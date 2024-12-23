package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/wahyuoi/sbc/internal/model"
)

type AudioService struct {
}

func NewAudioService() *AudioService {
	return &AudioService{}
}

func (s *AudioService) ConvertAudio(ctx context.Context, audioBytes []byte, audioFormat model.AudioFormatType) ([]byte, error) {

	filename := uuid.New().String()

	// Create temporary input file
	inputFile, err := os.CreateTemp("", fmt.Sprintf("%s*.m4a", filename))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp input file: %w", err)
	}
	defer os.Remove(inputFile.Name())

	// Write input bytes to temp file
	if _, err := inputFile.Write(audioBytes); err != nil {
		return nil, fmt.Errorf("failed to write to temp file: %w", err)
	}
	inputFile.Close()

	// Create temporary output file
	outputFile, err := os.CreateTemp("", fmt.Sprintf("%s*.wav", filename))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp output file: %w", err)
	}
	defer os.Remove(outputFile.Name())
	outputFile.Close()

	// Convert audio using	ffmpeg command
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", inputFile.Name(),
		"-acodec", "pcm_s16le", // WAV codec
		"-ar", "44100", // Sample rate
		"-ac", "2", // Mono channel
		"-y", outputFile.Name()) // Overwrite existing file

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to convert audio: %w", err)
	}

	// Read converted file
	convertedBytes, err := os.ReadFile(outputFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read converted file: %w", err)
	}

	return convertedBytes, nil
}
