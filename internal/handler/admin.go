package handler

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/request"
	"github.com/aliazizii/IMDB-In-Golang/internal/response"
	"github.com/aliazizii/IMDB-In-Golang/internal/store"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Admin struct {
	Store store.IMDB
}

func (admin Admin) AddMovie(c echo.Context) error {
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
	if err := admin.Store.AddMovie(m); err != nil {
		if errors.Is(err, store.DuplicateMovie) {
			// log
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	// log
	return c.JSON(http.StatusNoContent, "")
}

func (admin Admin) DeleteMovie(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	err = admin.Store.DeleteMovie(intID)
	if err != nil {
		if errors.Is(err, store.MovieNotFound) {
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (admin Admin) UpdateMovie(c echo.Context) error {
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
	err = admin.Store.UpdateMovie(intID, m)
	if err != nil {
		if errors.Is(err, store.MovieNotFound) {
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (admin Admin) UpdateComment(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	var req request.UpdateComment
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}

	err = admin.Store.UpdateComment(intID, req.Approved)
	if err != nil {
		if errors.Is(err, store.CommentNotFound) {
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (admin Admin) DeleteComment(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	err = admin.Store.DeleteComment(intID)
	if err != nil {
		if errors.Is(err, store.CommentNotFound) {
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (admin Admin) AdminRegister(g *echo.Group) {
	g.POST("/movie", admin.AddMovie)
	g.PUT("/movie/:id", admin.UpdateMovie)
	g.DELETE("/movie/:id", admin.DeleteMovie)
	g.PUT("/comment/:id", admin.UpdateComment)
	g.DELETE("/comment/:id", admin.DeleteComment)
}
