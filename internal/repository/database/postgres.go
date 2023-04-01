package database

import (
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() (*gorm.DB, error) {
	host := viper.GetString("database.postgres.host")
	port := viper.GetString("database.postgres.port")
	username := viper.GetString("database.postgres.username")
	password := viper.GetString("database.postgres.password")
	dbname := viper.GetString("database.postgres.name")
	timezone := viper.GetString("database.postgres.timezone")

	dsn := "host=" + host + " user=" + username + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable TimeZone=" + timezone
	db, err := gorm.Open(postgres.Open(dsn))

	return db, err
}
