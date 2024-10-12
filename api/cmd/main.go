package main

import (
	"biometry-hack-2024-api/internal/routing"
	"biometry-hack-2024-api/pkg/config"
	"biometry-hack-2024-api/pkg/logging"
	"biometry-hack-2024-api/pkg/messaging"
	"biometry-hack-2024-api/pkg/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	s := gin.Default()

	s.StaticFile("/docs/swagger.json", "./docs/swagger.json")

	s.GET("/api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

	log := logging.NewLogger()
	log.Println("logger initialized successfully")

	config.NewConfig()
	log.Println("configuration initialized successfully")

	db := storage.NewDBClient()
	log.Println("database client initialized successfully")
	svc := storage.NewS3Client()
	log.Println("s3 client initialized successfully")

	msg := messaging.NewNATSConn()
	defer msg.Drain()
	log.Println("message broker initialized successfully")

	routing.InitRouting(s, db, svc, msg, log)

	if err := s.Run(":8080"); err != nil {
		panic(err)
	}
}
