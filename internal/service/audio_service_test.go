package service_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wahyuoi/sbc/internal/model"
	"github.com/wahyuoi/sbc/internal/service"
)

func TestAudioService_ConvertAudio(t *testing.T) {
	service := service.NewAudioService()
	ctx := context.Background()

	t.Run("success convert audio", func(t *testing.T) {
		file, err := os.Open("../../sample/eleven.m4a")
		assert.NoError(t, err)
		defer file.Close()

		// Create a small test WAV file in memory
		inputAudio, err := io.ReadAll(file)
		assert.NoError(t, err)

		outputAudio, err := service.ConvertAudio(ctx, inputAudio, model.AudioFormatTypeWav)
		assert.NoError(t, err)
		assert.NotNil(t, outputAudio)
		assert.True(t, len(outputAudio) > len(inputAudio))
		fmt.Println(len(outputAudio))
	})

	t.Run("error when ffmpeg fails", func(t *testing.T) {
		// Invalid audio data should cause ffmpeg to fail
		invalidAudio := []byte("not valid audio")

		outputAudio, err := service.ConvertAudio(ctx, invalidAudio, model.AudioFormatTypeWav)
		assert.Error(t, err)
		assert.Nil(t, outputAudio)
	})

	t.Run("error with nil input", func(t *testing.T) {
		outputAudio, err := service.ConvertAudio(ctx, nil, model.AudioFormatTypeWav)
		assert.Error(t, err)
		assert.Nil(t, outputAudio)
	})
}
