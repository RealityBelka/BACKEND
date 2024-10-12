package handlers

import (
	"biometry-hack-2024-api/internal/models"
	"biometry-hack-2024-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPhotoExtension = "jpeg"
	DefaultAudioExtension = "aac"
	Megabyte              = 1024 * 1024 // Megabyte = 1024 Kilobytes = 1024 * 1024 bytes

	ExpectedPhotoHeader = "image/jpeg"
	ExpectedAudioHeader = "audio/aac"
)

func NewErrorResponse(err string) models.ErrorResponse {
	return models.ErrorResponse{
		Error: err,
	}
}

type biometryHandlers struct {
	service service.BiometryService
}

func NewBiometryHandlers(service service.BiometryService) BiometryHandlers {
	return biometryHandlers{
		service: service,
	}
}

// CreateFaceBiometry godoc
// @Summary      Создание лицевой биометрии
// @Description  Создает биометрический профиль на основе фотографии
// @Tags         FaceBiometry
// @Accept       image/jpeg
// @Produce      json
// @Param        file  body  string  true  "Фотография лица в формате JPEG" format(binary)
// @Success      201    {object}  models.PhotoAnalyzeResponse  "Успешное создание профиля"
// @Failure      400    {object}  models.ErrorResponse         "Ошибка чтения тела запроса"
// @Failure      413                                           "Размер фотографии превышает 1МБ"
// @Failure      415    {object}  models.ErrorResponse         "Неправильный тип контента (не image/jpeg)"
// @Failure      422    {object}  models.PhotoAnalyzeResponse  "Фотография не удовлетворяет требованиям"
// @Failure      500    {object}  models.ErrorResponse         "Внутренняя ошибка сервера"
// @Router       /face [post]
func (h biometryHandlers) CreateFaceBiometry(c *gin.Context) {
	if c.ContentType() != ExpectedPhotoHeader {
		c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, models.ErrorResponse{
			Error: "request content type is not image/jpeg",
		})
		return
	}

	photo, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if len(photo) > Megabyte {
		c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		return
	}

	response, err := h.service.CreateFaceBiometry(c.Request.Context(), photo, DefaultPhotoExtension)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if !response.Ok {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// CreateVoiceBiometry godoc
// @Summary      Создание голосовой биометрии
// @Description  Создает биометрический профиль на основе голоса
// @Tags         VoiceBiometry
// @Accept       audio/*
// @Produce      json
// @Param        file  body  string  true  "Аудиофайл с голосом" format(binary)
// @Success      201    {object}  models.AudioAnalyzeResponse  "Успешное создание профиля"
// @Failure      400    {object}  models.ErrorResponse         "Ошибка чтения тела запроса"
// @Failure      413                                           "Размер аудио превышает 1МБ"
// @Failure      415    {object}  models.ErrorResponse         "Неправильный тип контента (не аудио)"
// @Failure      422    {object}  models.AudioAnalyzeResponse  "Аудио не удовлетворяет требованиям"
// @Failure      500    {object}  models.ErrorResponse         "Внутренняя ошибка сервера"
// @Router       /voice [post]
func (h biometryHandlers) CreateVoiceBiometry(c *gin.Context) {
	if c.ContentType() != ExpectedAudioHeader {

		c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, models.ErrorResponse{
			Error: "content type must be any correct MIME audio type",
		})
		return
	}

	audio, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if len(audio) > Megabyte {
		c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		return
	}

	response, err := h.service.CreateVoiceBiometry(c.Request.Context(), audio, DefaultAudioExtension)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if !response.Ok {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	c.JSON(http.StatusCreated, response)
}
