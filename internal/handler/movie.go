package handler

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/auth"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/request"
	"github.com/aliazizii/IMDB-In-Golang/internal/response"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/movie"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Movie struct {
	Store movie.Movie
}

func (imovie Movie) AddMovie(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		logrus.Error("Add movie: failed to extract claims: ", err)
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}

	if claims.Role != auth.AdminRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only admin can use this endpoint"))
	}

	var req request.Movie

	if err := c.Bind(&req); err != nil {
		logrus.Error("Add movie: can not bind the request: ", err)
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}

	if err := req.Validate(); err != nil {
		logrus.Error("Add movie: request validation field fails: ", err)
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
			logrus.Error("Add movie: try to add a duplicate movie: ", err)
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}

		logrus.Error("Add movie: an error happens when storing movie: ", err)
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}

	return c.JSON(http.StatusNoContent, "")
}

func (imovie Movie) DeleteMovie(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		logrus.Error("Delete movie: failed to extract claims: ", err)
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}

	if claims.Role != auth.AdminRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only admin can use this endpoint"))
	}

	id := c.Param("id")

	intID, err := strconv.Atoi(id)
	if err != nil {
		logrus.Error("Delete movie: can not parse path param to int: ", err)
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}

	err = imovie.Store.DeleteMovie(intID)
	if err != nil {
		if errors.Is(err, movie.MovieNotFound) {
			logrus.Error("Delete movie: try to delete movie that is not exist: ", err)
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		logrus.Error("Delete movie: an error happens when deleting movie: ", err)

		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}

	return c.JSON(http.StatusNoContent, "")
}

func (imovie Movie) UpdateMovie(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		logrus.Error("Update movie: failed to extract claims: ", err)
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}

	if claims.Role != auth.AdminRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only admin can use this endpoint"))
	}

	id := c.Param("id")

	intID, err := strconv.Atoi(id)
	if err != nil {
		logrus.Error("Update movie: can not parse path param to int: ", err)
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}

	var req request.Movie

	if err := c.Bind(&req); err != nil {
		logrus.Error("Update movie: can not bind the request: ", err)
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}

	if err := req.Validate(); err != nil {
		logrus.Error("Update movie: request validation field fails: ", err)
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}

	m := model.Movie{
		Name:        req.Name,
		Description: req.Description,
	}

	err = imovie.Store.UpdateMovie(intID, m)
	if err != nil {
		if errors.Is(err, movie.MovieNotFound) {
			logrus.Error("Update movie: try to update movie that is not exist: ", err)
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		logrus.Error("Update movie: an error happens when updating movie: ", err)

		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}

	return c.JSON(http.StatusNoContent, "")
}

func (imovie Movie) AllMovies(c echo.Context) error {
	movies, err := imovie.Store.AllMovies()
	if err != nil {
		logrus.Error("Getting all movie: an error happens when getting all movie: ", err)
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
		logrus.Error("Get a movie: can not parse path param to int: ", err)
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}

	mm, err := imovie.Store.Movie(intID)
	if err != nil {
		if errors.Is(err, movie.MovieNotFound) {
			logrus.Error("Getting a movie: try ao access a movie that is not exist: ", err)
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}

		logrus.Error("Getting a movie: an error happens when getting a movie: ", err)
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
		logrus.Error("Get comments: can not parse path param to int: ", err)
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}

	comments, err := imovie.Store.AllComments(intID)
	if err != nil {
		if errors.Is(err, movie.MovieNotFound) {
			logrus.Error("Get comments: try to access comment of a movie that is not exist: ", err)
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}

		logrus.Error("Get comments: an error happens when trying to access comments: ", err)
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}

	mm, err := imovie.Store.Movie(intID)
	if err != nil {
		logrus.Error("Get comments: an error happens when getting a movie: ", err)
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}

	responseComments := response.Comments{
		Movie:    mm.Name,
		Comments: make([]response.Comment, 0),
	}

	for _, cc := range comments {
		responseComments.Comments = append(responseComments.Comments, response.Comment{
			ID:       cc.ID,
			Author:   cc.AuthorUsername,
			Body:     cc.Text,
			Approved: cc.Approved,
		})
	}

	return c.JSON(http.StatusOK, responseComments)
}
