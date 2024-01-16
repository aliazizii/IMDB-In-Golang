package handler

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/auth"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/request"
	"github.com/aliazizii/IMDB-In-Golang/internal/response"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/movie"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Movie struct {
	Store movie.Movie
}

func (imovie Movie) AddMovie(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		// log
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	if claims.Role != auth.AdminRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only admin can use this endpoint"))
	}
	var req request.Movie
	if err := c.Bind(&req); err != nil {
		// log
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	m := model.Movie{
		Name:        req.Name,
		Description: req.Description,
		Rating:      0,
		NVote:       0,
		Comments:    make([]model.Comment, 0),
	}
	if err = imovie.Store.AddMovie(m); err != nil {
		if errors.Is(err, movie.DuplicateMovie) {
			// log
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	// log
	return c.JSON(http.StatusNoContent, "")
}

func (imovie Movie) DeleteMovie(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		// log
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	if claims.Role != auth.AdminRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only admin can use this endpoint"))
	}
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	err = imovie.Store.DeleteMovie(intID)
	if err != nil {
		if errors.Is(err, movie.MovieNotFound) {
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (imovie Movie) UpdateMovie(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		// log
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	if claims.Role != auth.AdminRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only admin can use this endpoint"))
	}

	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	var req request.Movie
	if err := c.Bind(&req); err != nil {
		// log
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	m := model.Movie{
		Name:        req.Name,
		Description: req.Description,
	}
	err = imovie.Store.UpdateMovie(intID, m)
	if err != nil {
		if errors.Is(err, movie.MovieNotFound) {
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (imovie Movie) AllMovies(c echo.Context) error {
	movies, err := imovie.Store.AllMovies()
	if err != nil {
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
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	mm, err := imovie.Store.Movie(intID)
	if err != nil {
		if errors.Is(err, movie.MovieNotFound) {
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
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
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	comments, err := imovie.Store.AllComments(intID)
	if err != nil {
		if errors.Is(err, movie.MovieNotFound) {
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	mm, err := imovie.Store.Movie(intID)
	if err != nil {
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
