package handler

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/auth"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/request"
	"github.com/aliazizii/IMDB-In-Golang/internal/response"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/comment"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type Comment struct {
	Store comment.Comment
}

func (icomment Comment) UpdateComment(c echo.Context) error {
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
	var req request.UpdateComment
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}

	err = icomment.Store.UpdateComment(intID, req.Approved)
	if err != nil {
		if errors.Is(err, comment.CommentNotFound) {
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (icomment Comment) DeleteComment(c echo.Context) error {
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
	err = icomment.Store.DeleteComment(intID)
	if err != nil {
		if errors.Is(err, comment.CommentNotFound) {
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (icomment Comment) Comment(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		// log
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	if claims.Role != auth.UserRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only admin can use this endpoint"))
	}

	var req request.Comment
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	m := model.Comment{
		MovieID:   req.MovieID,
		Text:      req.CommentBody,
		CreatedAt: time.Now(),
		Approved:  false,
	}
	err = icomment.Store.Comment(m)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}
