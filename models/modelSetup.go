package models

import "lightRoom/db"

func Init() {
	// Auto Migrate
	db.Db.AutoMigrate(&User{})
}
