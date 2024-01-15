package comment

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
)

var CommentNotFound = errors.New("this comment dose not exist")

type Comment interface {
	UpdateComment(int, bool) error
	DeleteComment(int) error
	Comment(model.Comment) error
}
