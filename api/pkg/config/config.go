package config

import "github.com/spf13/viper"

const (
	DBName     = "DB_NAME"
	DBHost     = "DB_HOST"
	DBPort     = "DB_PORT"
	DBUser     = "DB_USER"
	DBPassword = "DB_PASSWORD"

	S3URL        = "S3_URL"
	S3AccessKey  = "S3_ACCESS_KEY"
	S3SecretKey  = "S3_SECRET_KEY"
	S3BucketName = "S3_BUCKET_NAME"

	NATSHost = "NATS_HOST"
	NATSPort = "NATS_PORT"
)

func NewConfig() {
	viper.SetConfigFile("cfgs/api.env")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
