package service

import (
	"blogging_platform/errs"
	"blogging_platform/internal/models"
	repository "blogging_platform/internal/repositories"
	"blogging_platform/pkg/db"
	"blogging_platform/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB, username, password, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func SignIn(username, password string) (accessToken string, err error) {
	password = utils.GenerateHash(password)
	user, err := repository.GetUserByUsernameAndPassword(username, password)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return "", errs.ErrIncorrectUsernameOrPassword
		}
		return "", err
	}

	accessToken, err = GenerateToken(user.ID, user.Username, user.Role, user.Email)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func GoogleLogin(c *gin.Context) {
	// Generate the URL for Google login
	url := utils.GoogleOauthConfig.AuthCodeURL("randomstate", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	// Get the authorization code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
		return
	}

	// Exchange the code for a token
	token, err := utils.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Fetch user info using the token
	client := utils.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}
	defer resp.Body.Close()

	// Parse the user info
	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Check if the user exists in the database
	var user models.User
	err = db.GetDBConn().Where("email = ?", userInfo.Email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: User not registered"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	accessToken, err := GenerateToken(user.ID, user.Username, user.Role, user.Email)
	if err != nil {
		return
	}

	// User is registered, respond with success
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"name":        user.Username,
			"email":       user.Email,
			"accesstoken": accessToken,
		},
	})
}
