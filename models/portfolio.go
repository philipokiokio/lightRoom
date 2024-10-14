package models

import (
	"github.com/google/uuid"
	"lightRoom/db"
	"time"
)

type Portfolio struct {
	ID              uuid.UUID `gorm:"primaryKey unique not null" json:"id"`
	Title           string    `json:"name"`
	Description     string    `json:"description"`
	Price           int       `json:"price"`
	Tags            []Tag     `gorm:"many2many:portfolio_tags;" json:"tags"`
	PaywalledImages []string  `gorm:"type:jsonb" json:"paywalled_images"`
	Images          []string  `gorm:"type:jsonb" json:"images"`
	UserID          uuid.UUID `gorm:"foreignKey:User;constraint:OnDelete:CASCADE;" json:"user_id"` // Relationship with User
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Tag struct {
	ID             uuid.UUID `gorm:"primaryKey unique not null" json:"id"`
	Title          string    `json:"title"`
	PortfolioCount int       `json:"portfolio_count"`
}

func CreateTag(tag Tag) error {
	return db.Db.Create(&tag).Error
}
func GetTags(title string) ([]Tag, error) {
	query := db.Db
	var fetchedTags []Tag
	if len(title) > 0 {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	err := query.Find(&fetchedTags).Error
	return fetchedTags, err
}
func GetTag(id uuid.UUID) (Tag, error) {
	var fetchedTag Tag

	err := db.Db.Where("id=?", id).First(&fetchedTag).Error

	return fetchedTag, err

}
func UpdateTagPortfolioCount(id uuid.UUID) error {
	var existingTag Tag
	_ = db.Db.Where("id=?", id).First(&existingTag).Error

	existingTag.PortfolioCount++
	return db.Db.Save(&existingTag).Error
}

func CreatePortfolio(portfolio Portfolio) error {

	return db.Db.Create(&portfolio).Error
}

func GetUserPortfolios(userID uuid.UUID, limit, offset int, tagID ...uuid.UUID) ([]Portfolio, error) {
	var userPortfolios []Portfolio

	query := db.Db.Where("user_id=?", userID)

	if len(tagID) > 0 {
		query.Where("tags IN (?)", tagID)
	}

	err := query.Limit(limit).Offset(offset).Find(&userPortfolios).Error

	return userPortfolios, err
}

func GetPortfolios(limit, offset int, tagID ...uuid.UUID) ([]Portfolio, error) {

	var portfolios []Portfolio
	query := db.Db
	if len(tagID) > 0 {
		query = db.Db.Where("tags in (?)", tagID)
	}
	err := query.Limit(limit).Offset(offset).Find(&portfolios).Error
	return portfolios, err
}

//func UpdateUserPortfolio(userID, portfolioID uuid.UUID, portfolio Portfolio) error {}
