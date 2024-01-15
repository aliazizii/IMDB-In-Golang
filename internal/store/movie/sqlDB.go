package movie

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
	if err := db.AutoMigrate(&model.Movie{}); err != nil {
		log.Fatal(err)
	}
	return SQL{
		DB: db,
	}
}

func (sql SQL) isExistMovie(name string, description string) (bool, error) {
	var m model.Movie
	query := sql.DB.Where("Name = ? AND Description = ?", name, description).First(&m)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if query.Error != nil {
		return false, query.Error
	}
	return true, nil
}

func (sql SQL) AddMovie(m model.Movie) error {
	ok, err := sql.isExistMovie(m.Name, m.Description)
	if err != nil {
		return err
	}
	if ok {
		return DuplicateMovie
	}
	return sql.DB.Save(&m).Error
}

func (sql SQL) DeleteMovie(i int) error {
	query := sql.DB.Delete(&model.Movie{}, i)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return MovieNotFound
	}
	if query.Error != nil {
		return query.Error
	}
	return nil
}

func (sql SQL) UpdateMovie(i int, m model.Movie) error {
	var movie model.Movie
	query := sql.DB.First(&movie, i)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return MovieNotFound
	}
	if query.Error != nil {
		return query.Error
	}
	return sql.DB.Save(&model.Movie{
		ID:          i,
		Name:        m.Name,
		Description: m.Description,
		Rating:      movie.Rating,
		NVote:       movie.NVote,
	}).Error
}

func (sql SQL) AllMovies() ([]model.Movie, error) {
	var movies []model.Movie
	if err := sql.DB.Find(&movies).Error; err != nil {
		return movies, err
	}
	return movies, nil
}

func (sql SQL) Movie(i int) (model.Movie, error) {
	var movie model.Movie
	query := sql.DB.First(&movie, i)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return movie, MovieNotFound
	}
	return movie, query.Error
}

func (sql SQL) AllComments(i int) ([]model.Comment, error) {
	comments := make([]model.Comment, 0)
	var movie model.Movie
	_, err := sql.Movie(i)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return comments, MovieNotFound
	}
	if err != nil {
		return comments, err
	}
	query := sql.DB.Preload("Comments").First(&movie, i)
	for _, comment := range movie.Comments {
		comments = append(comments, comment)
	}
	return comments, query.Error
}
