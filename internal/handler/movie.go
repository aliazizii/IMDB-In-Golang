package handler

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/auth"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/request"
	"github.com/aliazizii/IMDB-In-Golang/internal/response"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/movie"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Movie struct {
	Store  movie.Movie
	Logger *zap.Logger
}

func (imovie Movie) AddMovie(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		imovie.Logger.Error("Add movie: failed to extract claims:", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	if claims.Role != auth.AdminRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only admin can use this endpoint"))
	}
	if claims.StandardClaims.Valid() != nil {
		imovie.Logger.Error("Add movie: claim is invalid:", zap.Error(claims.StandardClaims.Valid()))
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("the token is expired, login again"))
	}
	var req request.Movie
	if err := c.Bind(&req); err != nil {
		imovie.Logger.Error("Add movie: can not bind the request:", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	if err := req.Validate(); err != nil {
		imovie.Logger.Error("Add movie: request validation field fails", zap.Error(err), zap.Any("request", req))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	m := model.Movie{
		Name:        req.Name,
		Description: req.Description,
		Rating:      0,
		NVote:       0,
		// should be removed or not?
		Comments: make([]model.Comment, 0),
	}
	if err = imovie.Store.AddMovie(m); err != nil {
		if errors.Is(err, movie.DuplicateMovie) {
			imovie.Logger.Error("Add movie: try to add a duplicate movie", zap.Error(err), zap.Any("movie", m))
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		imovie.Logger.Error("Add movie: an error happens when storing movie", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (imovie Movie) DeleteMovie(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		imovie.Logger.Error("Delete movie: failed to extract claims:", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	if claims.Role != auth.AdminRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only admin can use this endpoint"))
	}
	if claims.StandardClaims.Valid() != nil {
		imovie.Logger.Error("Delete movie: claim is invalid:", zap.Error(claims.StandardClaims.Valid()))
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("the token is expired, login again"))
	}
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		imovie.Logger.Error("Delete movie: can not parse path param to int:", zap.Error(err), zap.Any("id", id))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	err = imovie.Store.DeleteMovie(intID)
	if err != nil {
		if errors.Is(err, movie.MovieNotFound) {
			imovie.Logger.Error("Delete movie: try to delete movie that is not exist", zap.Error(err))
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		imovie.Logger.Error("Delete movie: an error happens when deleting movie", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (imovie Movie) UpdateMovie(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		imovie.Logger.Error("Update movie: failed to extract claims:", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	if claims.Role != auth.AdminRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only admin can use this endpoint"))
	}
	if claims.StandardClaims.Valid() != nil {
		imovie.Logger.Error("Update movie: claim is invalid:", zap.Error(claims.StandardClaims.Valid()))
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("the token is expired, login again"))
	}

	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		imovie.Logger.Error("Update movie: can not parse path param to int:", zap.Error(err), zap.Any("id", id))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	var req request.Movie
	if err := c.Bind(&req); err != nil {
		imovie.Logger.Error("Update movie: can not bind the request:", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	if err := req.Validate(); err != nil {
		imovie.Logger.Error("Update movie: request validation field fails", zap.Error(err), zap.Any("request", req))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	m := model.Movie{
		Name:        req.Name,
		Description: req.Description,
	}
	err = imovie.Store.UpdateMovie(intID, m)
	if err != nil {
		if errors.Is(err, movie.MovieNotFound) {
			imovie.Logger.Error("Update movie: try to update movie that is not exist", zap.Error(err))
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		imovie.Logger.Error("Update movie: an error happens when updating movie", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (imovie Movie) AllMovies(c echo.Context) error {
	movies, err := imovie.Store.AllMovies()
	if err != nil {
		imovie.Logger.Error("Getting all movie: an error happens when getting all movie", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	responseMovies := response.AllMovies{
		Movies: make([]response.Movie, 0),
	}
	for _, mm := range movies {
		responseMovies.Movies = append(responseMovies.Movies, response.Movie{
			ID:          mm.ID,
			Name:        mm.Name,
			Description: mm.Description,
			Rating:      mm.Rating,
		})
	}
	return c.JSON(http.StatusOK, responseMovies)
}

func (imovie Movie) Movie(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		imovie.Logger.Error("Get a movie: can not parse path param to int:", zap.Error(err), zap.Any("id", id))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	mm, err := imovie.Store.Movie(intID)
	if err != nil {
		if errors.Is(err, movie.MovieNotFound) {
			imovie.Logger.Error("Getting a movie: try ao access a movie that is not exist", zap.Error(err))
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		imovie.Logger.Error("Getting a movie: an error happens when getting a movie", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	responseMovie := response.Movie{
		ID:          mm.ID,
		Name:        mm.Name,
		Description: mm.Description,
		Rating:      mm.Rating,
	}
	return c.JSON(http.StatusOK, responseMovie)
}

func (imovie Movie) Comments(c echo.Context) error {
	id := c.QueryParam("movie")
	intID, err := strconv.Atoi(id)
	if err != nil {
		imovie.Logger.Error("Get comments: can not parse path param to int:", zap.Error(err), zap.Any("id", id))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	comments, err := imovie.Store.AllComments(intID)
	if err != nil {
		if errors.Is(err, movie.MovieNotFound) {
			imovie.Logger.Error("Get comments: try to access comment of a movie that is not exist", zap.Error(err))
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		imovie.Logger.Error("Get comments: an error happens when trying to access comments", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	mm, err := imovie.Store.Movie(intID)
	if err != nil {
		imovie.Logger.Error("Get comments: an error happens when getting a movie", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	responseComments := response.Comments{
		Movie:    mm.Name,
		Comments: make([]response.Comment, 0),
	}
	for _, cc := range comments {
		responseComments.Comments = append(responseComments.Comments, response.Comment{
			ID:     cc.ID,
			Author: cc.User.Username,
			Body:   cc.Text,
		})
	}
	return c.JSON(http.StatusOK, responseComments)
}
