package controllers

import (
	"blogging_platform/errs"
	"blogging_platform/internal/models"
	service "blogging_platform/internal/services"
	"blogging_platform/logger"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	userRole := c.GetString(userRoleCtx)
	if userRole != "admin" {
		handleError(c, errs.ErrPermissionDenied)
		return
	}

	users, err := service.GetAllUsers()
	if err != nil {
		logger.Error.Printf("[controllers.GetAllUsers] error getting all user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func GetUserByID(c *gin.Context) {
	userRole := c.GetString(userRoleCtx)
	if userRole != "admin" {
		handleError(c, errs.ErrPermissionDenied)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error.Printf("[controllers.GetUserByID] error getting user: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})

		return
	}

	user, err := service.GetUserByID(uint(id))
	if err != nil {
		logger.Error.Printf("[controllers.GetUserByID] error getting user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)

}

func CreateUser(c *gin.Context) {
	userRole := c.GetString(userRoleCtx)
	if userRole != "admin" {
		handleError(c, errs.ErrPermissionDenied)
		return
	}

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		logger.Error.Printf("[controllers.CreateUser] error creating user: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err := service.CreateUser(user)
	if err != nil {
		logger.Error.Printf("[controllers.CreateUser] error creating user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
	})

}

func EditUserByID(c *gin.Context) {
	userRole := c.GetString(userRoleCtx)
	if userRole != "admin" {
		handleError(c, errs.ErrPermissionDenied)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error.Printf("[controllers.EditUserByID] error editing user: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var userInput models.User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		logger.Error.Printf("[controllers.EditUserByID] error editing user: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input data",
		})
		return
	}

	updatedUser, err := service.EditUserByID(uint(id), userInput)
	if err != nil {
		logger.Error.Printf("[controllers.EditUserByID] error editing user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func DeleteUserByID(c *gin.Context) {
	userRole := c.GetString(userRoleCtx)
	if userRole != "admin" {
		handleError(c, errs.ErrPermissionDenied)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error.Printf("[controllers.DeleteUserByID] error deleating user: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	err = service.DeleteUserByID(uint(id))
	if err != nil {
		logger.Error.Printf("[controllers.DeleteUserByID] error deleating user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})
}
