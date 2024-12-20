package repository

import (
	"context"
)

type FileRepository interface {
	UploadAudio(ctx context.Context, audioBytes []byte, fileExt string) (string, error)
	DownloadAudio(ctx context.Context, audioPath string) ([]byte, error)
}
