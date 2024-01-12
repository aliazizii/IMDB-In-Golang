package handler

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/response"
	"github.com/aliazizii/IMDB-In-Golang/internal/store"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Public struct {
	Store store.IMDB
}

func (p Public) AllMovies(c echo.Context) error {
	movies, err := p.Store.AllMovies()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	responseMovies := response.AllMovies{
		Movies: make([]response.Movie, 0),
	}
	for _, movie := range movies {
		responseMovies.Movies = append(responseMovies.Movies, response.Movie{
			ID:          movie.ID,
			Name:        movie.Name,
			Description: movie.Description,
			Rating:      movie.Rating,
		})
	}
	return c.JSON(http.StatusOK, responseMovies)
}

func (p Public) Movie(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	movie, err := p.Store.Movie(intID)
	if err != nil {
		if errors.Is(err, store.MovieNotFound) {
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	responseMovie := response.Movie{
		ID:          movie.ID,
		Name:        movie.Name,
		Description: movie.Description,
		Rating:      movie.Rating,
	}
	return c.JSON(http.StatusOK, responseMovie)
}

func (p Public) Comments(c echo.Context) error {
	id := c.QueryParam("movie")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	comments, err := p.Store.AllComments(intID)
	if err != nil {
		if errors.Is(err, store.MovieNotFound) {
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	movie, err := p.Store.Movie(intID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	responseComments := response.Comments{
		Movie:    movie.Name,
		Comments: make([]response.Comment, 0),
	}
	for _, comment := range comments {
		responseComments.Comments = append(responseComments.Comments, response.Comment{
			ID:     comment.ID,
			Author: comment.User.Username,
			Body:   comment.Text,
		})
	}
	return c.JSON(http.StatusOK, responseComments)
}

func (p Public) PublicRegister(g *echo.Group) {
	g.GET("/comments", p.Comments)
	g.GET("/movies", p.AllMovies)
	g.GET("/movie/:id", p.Movie)
}
