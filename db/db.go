package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lightRoom/utils"
	"log"
)

var Db *gorm.DB

func Init() {
	postgresDSN := utils.Settings.PostgresDsn

	var err error
	Db, err = gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

}
