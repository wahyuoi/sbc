package repository

import (
	"context"
	"os"
	"path"

	"github.com/google/uuid"
)

var _ FileRepository = (*fileRepositoryLocal)(nil)

type fileRepositoryLocal struct {
	audioDir string
}

func NewFileRepositoryLocal() FileRepository {
	return &fileRepositoryLocal{
		audioDir: "/tmp/sbc/exercises/audio/",
	}
}

func (s *fileRepositoryLocal) UploadAudio(ctx context.Context, audioBytes []byte, fileExt string) (string, error) {
	fileName := uuid.New().String() + fileExt
	audioPath := path.Join(s.audioDir, fileName)

	if err := os.MkdirAll(s.audioDir, 0755); err != nil {
		return "", err
	}

	audioFile, err := os.Create(audioPath)
	if err != nil {
		return "", err
	}
	defer audioFile.Close()

	if _, err := audioFile.Write(audioBytes); err != nil {
		return "", err
	}

	return audioPath, nil
}

func (s *fileRepositoryLocal) DownloadAudio(ctx context.Context, audioPath string) ([]byte, error) {
	return os.ReadFile(audioPath)
}
