package handler

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/request"
	"github.com/aliazizii/IMDB-In-Golang/internal/response"
	"github.com/aliazizii/IMDB-In-Golang/internal/store"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type User struct {
	Store store.IMDB
}

func (admin Admin) Vote(c echo.Context) error {
	var req request.Vote
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	m := model.Vote{
		MovieID: req.MovieID,
		Rating:  req.Vote,
	}
	err = admin.Store.Vote(m)
	if err != nil {
		if errors.Is(err, store.DuplicateVote) {
			// log
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (admin Admin) Comment(c echo.Context) error {
	var req request.Comment
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	m := model.Comment{
		MovieID:   req.MovieID,
		Text:      req.CommentBody,
		CreatedAt: time.Now(),
		Approved:  false,
	}
	err = admin.Store.Comment(m)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (admin Admin) UserRegister(g *echo.Group) {
	g.POST("/vote", admin.Vote)
	g.POST("/comment", admin.Comment)
}
