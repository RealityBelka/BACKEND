package routing

import (
	"biometry-hack-2024-api/internal/handlers"
	"biometry-hack-2024-api/internal/repository"
	"biometry-hack-2024-api/internal/service"
	"database/sql"
	"log"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

func InitRouting(rg *gin.Engine, db *sql.DB, svc *s3.S3, msg *nats.Conn, log *log.Logger) {
	repo := repository.NewBiometryRepository(db, svc)
	service := service.NewBiometryService(repo, msg, log)
	handlers := handlers.NewBiometryHandlers(service)

	rg.POST("/voice", handlers.CreateVoiceBiometry)
	rg.POST("/face", handlers.CreateFaceBiometry)
}
