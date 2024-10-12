package service

import (
	"biometry-hack-2024-api/internal/models"
	"biometry-hack-2024-api/internal/repository"
	"biometry-hack-2024-api/pkg/messaging"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

const (
	photoExt       = "jpeg"
	requestTimeout = time.Duration(2 * time.Second)
)

type biometryService struct {
	repo repository.BiometryRepository
	nc   *nats.Conn
	log  *log.Logger
}

func NewBiometryService(repo repository.BiometryRepository, nc *nats.Conn, log *log.Logger) BiometryService {
	return biometryService{
		repo: repo,
		nc:   nc,
		log:  log,
	}
}

func (s biometryService) CreateVoiceBiometry(ctx context.Context, audio []byte, ext string) (models.AudioAnalyzeResponse, error) {
	reply, err := s.nc.Request(messaging.SubjectAudio, audio, requestTimeout)
	if err != nil {
		return models.AudioAnalyzeResponse{}, err
	}

	filename := uuid.NewString()

	if err := s.repo.CreateVoiceBiometry(ctx, fmt.Sprintf("%s.%s", filename, ext), audio); err != nil {
		return models.AudioAnalyzeResponse{}, err
	}

	var response models.AudioAnalyzeResponse
	if err := json.Unmarshal(reply.Data, &response); err != nil {
		return models.AudioAnalyzeResponse{}, err
	}

	if response.Ok {
		response.Message = "Голосовая биометрия создана успешно"
	}

	return response, nil
}

func (s biometryService) CreateFaceBiometry(ctx context.Context, photo []byte, ext string) (models.PhotoAnalyzeResponse, error) {
	reply, err := s.nc.Request(messaging.SubjectPhoto, photo, requestTimeout)
	if err != nil {
		return models.PhotoAnalyzeResponse{}, err
	}

	filename := uuid.NewString()

	if err := s.repo.CreateFaceBiometry(ctx, fmt.Sprintf("%s.%s", filename, photoExt), photo); err != nil {
		return models.PhotoAnalyzeResponse{}, err
	}

	var response models.PhotoAnalyzeResponse
	if err := json.Unmarshal(reply.Data, &response); err != nil {
		return models.PhotoAnalyzeResponse{}, err
	}

	if response.Ok {
		response.Message = "Лицевая биометрия создана успешно"
	}

	return response, nil
}
