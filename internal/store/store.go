package store

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
)

var DuplicateMovie = errors.New("This movie already exists")
var MovieNotFound = errors.New("This movie dosen't exist")
var DuplicateVote = errors.New("This vote already exists")

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
