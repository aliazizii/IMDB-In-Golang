package user

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"gorm.io/gorm"
	"log"
)

type SQL struct {
	DB *gorm.DB
}

func NewSQL(db *gorm.DB) SQL {
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal(err)
	}
	return SQL{
		DB: db,
	}
}

func (sql SQL) Find(username string) (model.User, error) {
	var user model.User
	query := sql.DB.Where("Username = ?", username).First(&user)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return user, UserNotFound
	}
	if query.Error != nil {
		return user, query.Error
	}
	return user, nil
}

func (sql SQL) Save(user model.User) error {
	_, err := sql.Find(user.Username)
	if err == nil {
		return DuplictateUser
	}
	return sql.DB.Save(&user).Error
}
