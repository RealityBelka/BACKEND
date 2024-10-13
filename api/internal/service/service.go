package service

import (
	"biometry-hack-2024-api/internal/models"
	"context"
)

type BiometryService interface {
	CreateVoiceBiometry(ctx context.Context, audio []byte, digits []int, ext string) (models.AudioAnalyzeResponse, error)
	CreateFaceBiometry(ctx context.Context, photo []byte, border []int, ext string) (models.PhotoAnalyzeResponse, error)
}
