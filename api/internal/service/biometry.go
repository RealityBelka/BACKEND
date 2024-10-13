package service

import (
	"biometry-hack-2024-api/internal/convertors"
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
	photoExt = "jpeg"

	photoRequestTimeout = time.Duration(2 * time.Second)
	voiceRequestTimeout = time.Duration(15 * time.Second)
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

func (s biometryService) CreateVoiceBiometry(ctx context.Context, audio []byte, digits []int, ext string) (models.AudioAnalyzeResponse, error) {
	bytes, err := convertors.FromBinNumbersPairToBinary(audio, digits)
	if err != nil {
		return models.AudioAnalyzeResponse{}, err
	}

	reply, err := s.nc.Request(messaging.SubjectAudio, bytes, voiceRequestTimeout)
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

func (s biometryService) CreateFaceBiometry(ctx context.Context, photo []byte, border []int, ext string) (models.PhotoAnalyzeResponse, error) {
	bytes, err := convertors.FromBinNumbersPairToBinary(photo, border)
	if err != nil {
		return models.PhotoAnalyzeResponse{}, err
	}

	reply, err := s.nc.Request(messaging.SubjectPhoto, bytes, photoRequestTimeout)
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
