package models

import (
	"github.com/google/uuid"
	"lightRoom/db"
)

type Role string

const (
	user  Role = "User"
	admin Role = "Admin"
)

type User struct {
	ID         uuid.UUID `gorm:"primaryKey unique not null" json:"user_id"`
	Name       string    `json:"name"`
	Email      string    `gorm:"unique not null" json:"email"`
	Password   string    `gorm:"unique not null" json:"password"`
	IsVerified bool      `json:"is_verified"`
}

func CreateUser(user User) error {

	return db.Db.Create(&user).Error
}
func FetchViaMail(email string) (User, error) {

	var user User
	err := db.Db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, err
}
func GetUser(user_id uuid.UUID) (User, error) {

	var fetchedUser = User{ID: user_id}
	err := db.Db.Where("id = ?", user_id).First(&fetchedUser).Error

	return fetchedUser, err
}

func UpdateUser(user_id uuid.UUID, updateUser User) error {

	var existingUser User
	_ = db.Db.Where("id = ?", updateUser.ID).First(&existingUser).Error

	//	Updating the fields of the existing User with the new value
	err := db.Db.Model(&existingUser).Updates(updateUser).Error

	if err != nil {
		return err
	}
	return nil

}

func DeleteUser(user_id uuid.UUID) error {
	return db.Db.Delete(&User{ID: user_id}).Error
}

func GetUsers(limit int, offset int) ([]User, error) {

	var users []User

	err := db.Db.Model(&users).Limit(limit).Offset(offset).Find(&users).Error

	return users, err
}
