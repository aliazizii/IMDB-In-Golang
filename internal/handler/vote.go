package handler

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/request"
	"github.com/aliazizii/IMDB-In-Golang/internal/response"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/vote"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Vote struct {
	Store vote.Vote
}

func (ivote Vote) Vote(c echo.Context) error {
	var req request.Vote
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	m := model.Vote{
		MovieID: req.MovieID,
		Rating:  req.Vote,
	}
	err = ivote.Store.Vote(m)
	if err != nil {
		if errors.Is(err, vote.DuplicateVote) {
			// log
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}
