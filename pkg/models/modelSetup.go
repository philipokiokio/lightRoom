package models

import (
	"lightRoom/pkg/db"
)

func Init() {
	// Auto Migrate
	db.Db.AutoMigrate(&User{}, &Tag{}, &Portfolio{})
}

func GetStringPointer(pointedString *string) string {
	if pointedString != nil && *pointedString != "" {
		return *pointedString
	}
	return ""
}
