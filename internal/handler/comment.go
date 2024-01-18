package handler

import (
	"errors"
	"github.com/aliazizii/IMDB-In-Golang/internal/auth"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/request"
	"github.com/aliazizii/IMDB-In-Golang/internal/response"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/comment"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type Comment struct {
	Store  comment.Comment
	Logger *zap.Logger
}

func (icomment Comment) UpdateComment(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		icomment.Logger.Error("Update comment: failed to extract claims:", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	if claims.Role != auth.AdminRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only admin can use this endpoint"))
	}
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		icomment.Logger.Error("Update comment: can not parse path param to int:", zap.Error(err), zap.Any("id", id))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	var req request.UpdateComment
	err = c.Bind(&req)
	if err != nil {
		icomment.Logger.Error("Update comment: can not bind the request:", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	if err := req.Validate(); err != nil {
		icomment.Logger.Error("Update comment: request validation field fails", zap.Error(err), zap.Any("request", req))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	err = icomment.Store.UpdateComment(intID, req.Approved)
	if err != nil {
		if errors.Is(err, comment.CommentNotFound) {
			icomment.Logger.Error("Update comment: try to update comment that is not exist:", zap.Error(err))
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		icomment.Logger.Error("Update comment: an error happens when changing situation of comment:", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (icomment Comment) DeleteComment(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		icomment.Logger.Error("Delete comment: failed to extract claims:", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	if claims.Role != auth.AdminRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only admin can use this endpoint"))
	}
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		icomment.Logger.Error("Delete comment: can not parse path param to int:", zap.Error(err), zap.Any("id", id))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	err = icomment.Store.DeleteComment(intID)
	if err != nil {
		if errors.Is(err, comment.CommentNotFound) {
			icomment.Logger.Error("Delete comment: try to delete comment that is not exist:", zap.Error(err))
			return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
		}
		icomment.Logger.Error("Delete comment: an error happens when deleting comment:", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}

func (icomment Comment) Comment(c echo.Context) error {
	claims, err := auth.ExtractJWT(c)
	if err != nil {
		icomment.Logger.Error("Comment: failed to extract claims:", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	if claims.Role != auth.UserRoleCode {
		return c.JSON(http.StatusUnauthorized, response.CreateErrMessageResponse("only users can use this endpoint"))
	}
	var req request.Comment
	err = c.Bind(&req)
	if err != nil {
		icomment.Logger.Error("Comment: can not bind the request:", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	if err := req.Validate(); err != nil {
		icomment.Logger.Error("Comment: request validation field fails", zap.Error(err), zap.Any("request", req))
		return c.JSON(http.StatusBadRequest, response.CreateErrMessageResponse("bad request"))
	}
	m := model.Comment{
		MovieID:        req.MovieID,
		Text:           req.CommentBody,
		CreatedAt:      time.Now(),
		Approved:       false,
		AuthorUsername: claims.ID,
	}
	err = icomment.Store.Comment(m)
	if err != nil {
		icomment.Logger.Error("Comment: an error happens when commenting:", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.CreateErrMessageResponse("there is an internal issue"))
	}
	return c.JSON(http.StatusNoContent, "")
}
