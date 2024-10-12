package repository

import "context"

type BiometryRepository interface {
	CreateVoiceBiometry(ctx context.Context, filename string, audio []byte) error
	CreateFaceBiometry(ctx context.Context, filename string, photo []byte) error
}
