package store

import "github.com/aliazizii/IMDB-In-Golang/internal/model"

type IMDB interface {
	AddMovie(model.Movie) error
	DeleteMovie(int) error
	UpdateMovie(int, model.Movie) error
	Vote(model.Vote) error
	Comment(model.Comment) error
	AllMovies() ([]model.Movie, error)
	Movie(int) (model.Movie, error)
	AllComments() ([]model.Comment, error)
}
