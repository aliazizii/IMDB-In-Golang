package comment

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
	if err := db.AutoMigrate(&model.Comment{}); err != nil {
		log.Fatal(err)
	}
	return SQL{
		DB: db,
	}
}

func (sql SQL) UpdateComment(i int, approved bool) error {
	var comment model.Comment
	query := sql.DB.First(&comment, i)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return CommentNotFound
	}
	if query.Error != nil {
		return query.Error
	}
	comment.Approved = approved
	return sql.DB.Save(&comment).Error
}

func (sql SQL) DeleteComment(i int) error {
	var comment model.Comment
	query := sql.DB.First(&comment, i)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return CommentNotFound
	}
	if query.Error != nil {
		return query.Error
	}
	return sql.DB.Delete(&model.Comment{}, i).Error
}

func (sql SQL) Comment(c model.Comment) error {
	return sql.DB.Save(&c).Error
}
