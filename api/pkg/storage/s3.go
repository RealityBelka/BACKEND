package storage

import (
	"biometry-hack-2024-api/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

func NewS3Client() *s3.S3 {
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentials(
			viper.GetString(config.S3AccessKey),
			viper.GetString(config.S3SecretKey),
			"",
		),
		Endpoint: aws.String(viper.GetString(config.S3URL)),
		Region:   aws.String("ru-1"),

		S3ForcePathStyle: aws.Bool(true),
	}

	sess := session.Must(session.NewSession(&cfg))

	svc := s3.New(sess)

	_, err := svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(viper.GetString(config.S3BucketName)),
	})
	if err != nil {
		_, err = svc.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(viper.GetString(config.S3BucketName)),
		})
		if err != nil {
			panic(err)
		}

	}

	return svc
}
