package movie

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
)

var DuplicateMovie = errors.New("this movie already exists")
var MovieNotFound = errors.New("this movie dose not exist")

type Movie interface {
	AddMovie(model.Movie) error
	DeleteMovie(int) error
	UpdateMovie(int, model.Movie) error
	AllMovies() ([]model.Movie, error)
	Movie(int) (model.Movie, error)
	isExistMovie(string, string) (bool, error)
	AllComments(int) ([]model.Comment, error)
}
