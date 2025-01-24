package controllers

import (
	"delivery/errs"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error) {
	if errors.Is(err, errs.ErrUsernameUniquenessFailed) ||
		errors.Is(err, errs.ErrIncorrectUsernameOrPassword) {
		c.JSON(http.StatusBadRequest, newErrorResponse(err.Error()))
	} else if errors.Is(err, errs.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, newErrorResponse(err.Error()))
	} else if errors.Is(err, errs.ErrPermissionDenied) {
		c.JSON(http.StatusForbidden, newErrorResponse(err.Error()))
	} else {
		c.JSON(http.StatusInternalServerError, newErrorResponse(errs.ErrSomethingWentWrong.Error()))
	}
}
