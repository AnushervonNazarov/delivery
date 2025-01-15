package service

import (
	"context"
	"delivery/errs"
	"delivery/internal/models"
	repository "delivery/internal/repositories"
	"delivery/pkg/db"
	"delivery/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"googlemaps.github.io/maps"
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

// GetDirections fetches directions from Google Maps and returns the route details
func GetDirections(origin, destination string) (map[string]interface{}, error) {
	// Initialize the Google Maps client
	client, err := maps.NewClient(maps.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		return nil, fmt.Errorf("failed to create maps client: %v", err)
	}

	// Make the directions request
	req := &maps.DirectionsRequest{
		Origin:      origin,
		Destination: destination,
		Mode:        maps.TravelModeDriving,
	}

	// Request directions
	resp, _, err := client.Directions(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get directions: %v", err)
	}

	// Structure the detailed response
	if len(resp) > 0 {
		route := resp[0]
		steps := []map[string]interface{}{}

		// Iterate through steps to create a clearer output
		for _, step := range route.Legs[0].Steps {
			// Convert the duration to seconds
			durationInSeconds := int(step.Duration.Seconds())

			// Check for short durations (less than 1 minute)
			if durationInSeconds < 60 {
				steps = append(steps, map[string]interface{}{
					"instruction": step.HTMLInstructions,
					"distance":    step.Distance.HumanReadable,
					"duration":    fmt.Sprintf("%d seconds", durationInSeconds),
				})
			} else {
				// Calculate hours and minutes
				hours := durationInSeconds / 3600
				minutes := (durationInSeconds % 3600) / 60

				steps = append(steps, map[string]interface{}{
					"instruction": step.HTMLInstructions,
					"distance":    step.Distance.HumanReadable,
					"duration":    fmt.Sprintf("%02d hours %02d minutes", hours, minutes),
				})
			}
		}

		// Calculate total duration for the entire route
		totalDurationInSeconds := int(route.Legs[0].Duration.Seconds())
		totalHours := totalDurationInSeconds / 3600
		totalMinutes := (totalDurationInSeconds % 3600) / 60

		// Return a structured and clear response
		return map[string]interface{}{
			"summary":        route.Summary,
			"total_distance": route.Legs[0].Distance.HumanReadable,
			"total_duration": fmt.Sprintf("%02d hours %02d minutes", totalHours, totalMinutes),
			"steps":          steps,
		}, nil
	}

	return nil, fmt.Errorf("no routes found")
}
