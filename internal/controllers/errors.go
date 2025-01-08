package controllers

import (
	"delivery/errs"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// type ErrorResponse struct {
// 	Error string `json:"error"`
// }

// func newErrorResponse(message string) ErrorResponse {
// 	return ErrorResponse{
// 		Error: message,
// 	}
// }

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

// type defaultResponse struct {
// 	Message string `json:"message"`
// }

// func newDefaultResponse(message string) defaultResponse {
// 	return defaultResponse{
// 		Message: message,
// 	}
// }

// type accessTokenResponse struct {
// 	AccessToken string `json:"access_token"`
// }
