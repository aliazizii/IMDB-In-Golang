package vote

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
)

var DuplicateVote = errors.New("this vote already exists")

type Vote interface {
	Vote(model.Vote) error
	isExistVote(int, string) (bool, error)
}
