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

type myVote model.Vote

func (v myVote) AfterSave(db *gorm.DB) (err error) {
	var m model.Movie
	db.First(&m, v.MovieID)
	nVote := m.NVote
	movieRating := m.Rating
	m.Rating = ((nVote * movieRating) + v.Rating) / (nVote + 1)
	m.NVote += 1
	db.Save(&m)
	return
}

func (sql SQL) Vote(v model.Vote) error {
	ok, err := sql.isExistVote(v.MovieID, v.UserUsername)
	if err != nil {
		return err
	}
	if ok {
		return DuplicateVote
	}
	return sql.DB.Save(&v).Error
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
