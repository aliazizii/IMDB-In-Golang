package store

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
)

var DuplicateMovie = errors.New("this movie already exists")
var MovieNotFound = errors.New("this movie dose not exist")
var DuplicateVote = errors.New("this vote already exists")
var CommentNotFound = errors.New("this comment dose not exist")

type IMDB interface {
	// Admin

	AddMovie(model.Movie) error
	DeleteMovie(int) error
	UpdateMovie(int, model.Movie) error
	UpdateComment(int, bool) error
	DeleteComment(int) error

	// User

	Vote(model.Vote) error
	Comment(model.Comment) error

	// Public

	AllMovies() ([]model.Movie, error)
	Movie(int) (model.Movie, error)
	AllComments(int) ([]model.Comment, error)

	// Auxiliary
	isExistMovie(string, string) (bool, error)
	isExistVote(int, int) (bool, error)
}
