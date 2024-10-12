package storage

import (
	"biometry-hack-2024-api/pkg/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const driverName = "postgres"

func NewDBClient() *sql.DB {
	addr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString(config.DBHost),
		viper.GetInt(config.DBPort),
		viper.GetString(config.DBUser),
		viper.GetString(config.DBPassword),
		viper.GetString(config.DBName),
	)

	db, err := sql.Open(driverName, addr)
	if err != nil {
		panic(err)
	}

	return db
}
