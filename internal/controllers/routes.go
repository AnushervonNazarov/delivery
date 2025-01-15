package controllers

import (
	"delivery/configs"
	repository "delivery/internal/repositories"
	service "delivery/internal/services"
	"delivery/pkg/db"
	"delivery/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	r.GET("/directions", GetDirections)

	r.Use(DatabaseMiddleware(db.GetDBConn()))

	r.POST("/orders", handlePlaceOrder)

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

func handlePlaceOrder(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	// Initialize dependencies here
	itemRepo := repository.NewItemRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(itemRepo, orderRepo)
	orderController := NewOrderController(orderService)

	// Call the controller method
	orderController.PlaceOrder(c)
}
