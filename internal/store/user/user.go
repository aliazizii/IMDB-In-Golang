package user

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
)

var UserNotFound = errors.New("this user does not exist")
var DuplictateUser = errors.New("this user already exist")

type User interface {
	Find(username string) (model.User, error)
	Save(user model.User) error
}
