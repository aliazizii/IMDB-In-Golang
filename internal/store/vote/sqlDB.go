package vote

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
	if err := db.AutoMigrate(&model.Vote{}); err != nil {
		log.Fatal(err)
	}
	return SQL{
		DB: db,
	}
}

func (sql SQL) Vote(v model.Vote) error {
	ok, err := sql.isExistVote(v.MovieID, v.UserUsername)
	if err != nil {
		return err
	}
	if ok {
		return DuplicateVote
	}
	return sql.DB.Model(&model.Vote{}).Save(&v).Error
}

func (sql SQL) isExistVote(movieID int, userUsername string) (bool, error) {
	var vote model.Vote
	query := sql.DB.Where("movie_id = ? AND user_username = ?", movieID, userUsername).First(&vote)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if query.Error != nil {
		return false, query.Error
	}
	return true, nil
}
