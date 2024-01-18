package handler

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/auth"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/request"
	"github.com/aliazizii/IMDB-In-Golang/internal/response"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/vote"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Vote struct {
	Store  vote.Vote
	Logger *zap.Logger
}

func (ivote Vote) Vote(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		ivote.Logger.Error("Vote: failed to extract claim", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	if claims.Role != auth.UserRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only users can use this endpoint"))
	}
	var req request.Vote
	err = c.Bind(&req)
	if err != nil {
		ivote.Logger.Error("Vote: can not bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	if err := req.Validate(); err != nil {
		ivote.Logger.Error("Vote: request validation field fails", zap.Error(err), zap.Any("request", req))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	m := model.Vote{
		MovieID:      req.MovieID,
		Rating:       float64(req.Vote),
		UserUsername: claims.ID,
	}
	err = ivote.Store.Vote(m)
	if err != nil {
		if errors.Is(err, vote.DuplicateVote) {
			ivote.Logger.Error("Vote: already vote this movie", zap.Error(err))
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		ivote.Logger.Error("Vote: an error happens when voting a movie", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}
