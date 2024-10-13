package handlers

import (
	"biometry-hack-2024-api/internal/convertors"
	"biometry-hack-2024-api/internal/models"
	"biometry-hack-2024-api/internal/service"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	defaultPhotoExtension = "jpeg"
	defaultAudioExtension = "ogg"
	megabyte              = 1024 * 1024 // Megabyte = 1024 Kilobytes = 1024 * 1024 bytes

	contentTypeFormData = "multipart/form-data"
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
// @Accept       multipart/form-data
// @Produce      json
// @Param        photo    formData  file    true               "Фотография лица в формате JPEG"
// @Param        numbers  formData  string  true               "Координаты точек овала, в который должно поместиться лицо"
// @Success      201    {object}  models.PhotoAnalyzeResponse  "Успешное создание профиля"
// @Failure      400    {object}  models.ErrorResponse         "Ошибка чтения тела запроса"
// @Failure      413                                           "Размер фотографии превышает 1МБ"
// @Failure      415    {object}  models.ErrorResponse         "Неправильный тип контента (не image/jpeg)"
// @Failure      422    {object}  models.PhotoAnalyzeResponse  "Фотография не удовлетворяет требованиям"
// @Failure      500    {object}  models.ErrorResponse         "Внутренняя ошибка сервера"
// @Router       /face [post]
func (h biometryHandlers) CreateFaceBiometry(c *gin.Context) {
	if c.ContentType() != contentTypeFormData {
		c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, models.ErrorResponse{
			Error: fmt.Sprintf("content type must be '%s'", contentTypeFormData),
		})
		return
	}

	rawPhoto, err := c.FormFile("photo")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if rawPhoto.Size > megabyte {
		c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		return
	}

	photoFile, err := rawPhoto.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	defer photoFile.Close()

	photo, err := io.ReadAll(photoFile)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	rawBorder := c.PostForm("numbers")

	border, err := convertors.FromStringToIntSlice(string(rawBorder))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	response, err := h.service.CreateFaceBiometry(c.Request.Context(), photo, border, defaultPhotoExtension)
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
	if c.ContentType() != contentTypeFormData {
		c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, models.ErrorResponse{
			Error: fmt.Sprintf("content type must be '%s'", contentTypeFormData),
		})
		return
	}

	rawAudio, err := c.FormFile("audio")
	if err != nil {
		fmt.Println("no file with name 'audio'")
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if rawAudio.Size > megabyte {
		c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		return
	}

	audioFile, err := rawAudio.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	defer audioFile.Close()

	audio, err := io.ReadAll(audioFile)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	rawDigits := c.PostForm("numbers")

	digits, err := convertors.FromStringToIntSlice(string(rawDigits))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	response, err := h.service.CreateVoiceBiometry(c.Request.Context(), audio, digits, defaultAudioExtension)
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
