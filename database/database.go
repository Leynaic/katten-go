package database

import (
	"github.com/Leynaic/katten-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var instance *gorm.DB

func New(host string, user string, password string, name string, port string) {
	dsn := "host=" + host + " user=" + user + " password='" + password + "' dbname=" + name + " port=" + port + " sslmode=disable TimeZone=Europe/Paris"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		QueryFields: true,
		Logger:      logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	if db.AutoMigrate(&models.Cat{}) != nil {
		log.Fatal(err)
	}

	instance = db
}

func GetInstance() *gorm.DB {
	return instance
}
