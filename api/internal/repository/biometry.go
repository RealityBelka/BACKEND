package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
)

type biometryRepository struct {
	db  *sql.DB
	svc *s3.S3
}

func NewBiometryRepository(db *sql.DB, svc *s3.S3) BiometryRepository {
	return biometryRepository{
		db:  db,
		svc: svc,
	}
}

func (b biometryRepository) CreateVoiceBiometry(ctx context.Context, filename string, audio []byte) error {
	dbTX, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO voice_biometry (filename, audio) VALUES ($1, $2)
	`

	res, err := dbTX.ExecContext(ctx, query, filename, audio)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("no rows affected")
	}

	/* :)))
	_, err = b.svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(config.S3BucketName),
		Key:    aws.String(filename),
		Body:   body,
	})
	if err != nil {
		dbTX.Rollback()
		return err
	}
	*/

	if err = dbTX.Commit(); err != nil {
		dbTX.Rollback()

		/* :)))
		b.svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(config.S3BucketName),
			Key:    aws.String(filename),
		})
		*/

		return err
	}

	return nil
}

func (b biometryRepository) CreateFaceBiometry(ctx context.Context, filename string, photo []byte) error {
	dbTX, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Println("photo len:", len(photo))

	query := `
	INSERT INTO face_biometry (filename, photo) VALUES ($1, $2)
	`

	res, err := dbTX.ExecContext(ctx, query, filename, photo)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("no rows affected")
	}

	/* :)))
	_, err = b.svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(config.S3BucketName),
		Key:    aws.String(filename),
		Body:   body,
	})
	if err != nil {
		dbTX.Rollback()
		return err
	}
	*/

	if err = dbTX.Commit(); err != nil {
		dbTX.Rollback()

		/* :)))
		b.svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(config.S3BucketName),
			Key:    aws.String(filename),
		})
		*/

		return err
	}

	return nil
}
