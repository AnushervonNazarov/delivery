package controllers

import (
	"delivery/configs"
	service "delivery/internal/services"
	"delivery/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunRouts() *gin.Engine {
	utils.InitGoogleOAuth()

	r := gin.Default()
	gin.SetMode(configs.AppSettings.AppParams.GinMode)

	r.GET("/ping", PingPong)

	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", SignUp)
		auth.POST("/sign-in", SignIn)
	}

	r.GET("/google/login", service.GoogleLogin)
	r.GET("/callback/google", service.GoogleCallback)

	apiG := r.Group("/api", checkUserAuthentication)

	userG := apiG.Group("/users")
	{
		userG.GET("", GetAllUsers)
		userG.GET("/:id", GetUserByID)
		userG.POST("", CreateUser)
		userG.PUT("/:id", EditUserByID)
		userG.DELETE("/:id", DeleteUserByID)
	}
	return r
}

func PingPong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
