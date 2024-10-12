package messaging

import (
	"biometry-hack-2024-api/pkg/config"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

const (
	SubjectPhoto = "photo_subject"
	SubjectAudio = "audio_subject"
)

func NewNATSConn() *nats.Conn {
	url := fmt.Sprintf(
		"nats://%s:%d",
		viper.GetString(config.NATSHost),
		viper.GetInt(config.NATSPort),
	)

	conn, err := nats.Connect(url)
	if err != nil {
		panic(err)
	}

	return conn
}
