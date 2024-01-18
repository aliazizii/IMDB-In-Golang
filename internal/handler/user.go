package handler

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/auth"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/request"
	"github.com/aliazizii/IMDB-In-Golang/internal/response"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/user"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type User struct {
	Store     user.User
	JwtSecret string
	Logger    *zap.Logger
}

func (iuser User) Login(c echo.Context) error {
	var req request.User
	if err := c.Bind(&req); err != nil {
		iuser.Logger.Error("Login: can not bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	if err := req.Validate(); err != nil {
		iuser.Logger.Error("Login: request validation field fails", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	req.Username = strings.ToLower(req.Username)

	u, err := iuser.Store.Find(req.Username)
	if err != nil {
		if errors.Is(err, user.UserNotFound) {
			iuser.Logger.Error("Login: username not found", zap.Error(err))
			return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("the username is incorrect"))
		}
		iuser.Logger.Error("Login: an error happens when finding user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}

	if !auth.CheckPassword(req.Password, u.Password) {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("the password is incorrect"))
	}
	isAdmin := false
	if u.Role == auth.AdminRoleCode {
		isAdmin = true
	}
	token, err := auth.GenerateJWT(iuser.JwtSecret, u.Username, isAdmin)
	if err != nil {
		iuser.Logger.Error("Login: an error happens when generating jwt", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	resUser := response.User{
		ID:       u.ID,
		Username: u.Username,
		IsAdmin:  isAdmin,
		JWT:      token,
	}
	return c.JSON(http.StatusOK, resUser)
}

func (iuser User) SignUp(c echo.Context) error {
	var req request.User
	if err := c.Bind(&req); err != nil {
		iuser.Logger.Error("Signup: can not bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	if err := req.Validate(); err != nil {
		iuser.Logger.Error("Signup: request validation field fails", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	req.Username = strings.ToLower(req.Username)
	req.Password = auth.Hash(req.Password)
	u := model.User{
		Username: req.Username,
		Password: req.Password,
		Role:     auth.UserRoleCode,
	}
	err := iuser.Store.Save(u)
	if err != nil {
		if errors.Is(err, user.DuplictateUser) {
			iuser.Logger.Error("Signup: try to use duplicate username", zap.Error(err))
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("a user with this username already exist"))
		}
		iuser.Logger.Error("Signup: an error happens when storing user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	token, err := auth.GenerateJWT(iuser.JwtSecret, u.Username, false)
	if err != nil {
		iuser.Logger.Error("Signup: an error happens when generating jwt", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	resUser := response.User{
		ID:       u.ID,
		Username: u.Username,
		IsAdmin:  false,
		JWT:      token,
	}
	return c.JSON(http.StatusCreated, resUser)
}
