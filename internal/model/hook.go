package model

import (
	"gorm.io/gorm"
)

func (v Vote) AfterSave(db *gorm.DB) (err error) {
	var m Movie
	db.First(&m, v.MovieID)
	nVote := m.NVote
	movieRating := m.Rating
	m.Rating = ((nVote * movieRating) + v.Rating) / (nVote + 1)
	m.NVote += 1
	db.Save(&m)
	return
}
