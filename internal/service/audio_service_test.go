package service_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func BenchmarkAudioService_ConvertAudio(b *testing.B) {
	service := service.NewAudioService()
	ctx := context.Background()

	file, err := os.Open("../../sample/ninety.m4a")
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()

	inputAudio, err := io.ReadAll(file)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputAudio, err := service.ConvertAudio(ctx, inputAudio, model.AudioFormatTypeWav)
		if err != nil {
			b.Fatal(err)
		}
		if outputAudio == nil {
			b.Fatal("output audio is nil")
		}
	}
}

func BenchmarkAudioService_ConvertAudio_AllSamples(b *testing.B) {
	service := service.NewAudioService()
	ctx := context.Background()

	// Read all .m4a files from sample directory
	files, err := os.ReadDir("../../sample")
	if err != nil {
		b.Fatal(err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".m4a") {
			continue
		}

		b.Run(file.Name(), func(b *testing.B) {
			filePath := filepath.Join("../../sample", file.Name())
			f, err := os.Open(filePath)
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()

			inputAudio, err := io.ReadAll(f)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				outputAudio, err := service.ConvertAudio(ctx, inputAudio, model.AudioFormatTypeWav)
				if err != nil {
					b.Fatal(err)
				}
				if outputAudio == nil {
					b.Fatal("output audio is nil")
				}
			}
		})
	}
}
