package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neverloved-dev/goJWT/db"

	"github.com/golang-jwt/jwt"
	gomail "gopkg.in/mail.v2"
)

type User struct {
	GUID          string
	username      string
	email         string
	refresh_token string
}

var secretKey = []byte("myTopSecret-key")

func generateTokenPair(username string, ipaddress string, userguid string) (map[string]string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["ipaddress"] = ipaddress
	claims["name"] = username
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	t, err := token.SignedString([]byte(userguid))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["username"] = username
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rt, err := refreshToken.SignedString([]byte(userguid))
	if err != nil {
		return nil, err
	}
	var user User
	result := db.Database.Where("guid=?", userguid).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user: %v", result.Error)
	}

	user.refresh_token = rt
	if err := db.Database.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user refresh token: %v", err)
	}
	return map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}, nil

}

func refreshTokenPair(ipaddress string, userguid string) {

}

func handleReturnPong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func HandleGetTokens(c *gin.Context) {
	guid := c.Param("guid")
	ip := c.ClientIP()
	// make GORM request to get the username
	username := "username"
	jwtToken, refreshToken := generateTokenPair(username, ip, guid)
	c.JSON(200, gin.H{
		"token":         jwtToken,
		"refresh_token": refreshToken,
	})
}

func HandleRefreshTokens(c *gin.Context, guid string, refreshToken string) {
	// Get the IP address of the incoming request
	requestIP := c.ClientIP()

	// Parse the refresh token
	parsedToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(guid), nil
	})

	if err != nil || !parsedToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Extract claims from the token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	tokenIP := claims["ipaddress"].(string)

	// Compare IP address from token with the request IP
	if tokenIP != requestIP {
		// If IPs don't match, send an email to the user
		err := notifyUser(guid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to notify user"})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "IP address mismatch. Notification sent to user"})
		return
	}

	// If IPs match, generate new access and refresh tokens
	tokenPair, err := generateTokenPair(claims["username"].(string), requestIP, guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
		return
	}

	// Return the new tokens as JSON
	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokenPair["access_token"],
		"refresh_token": tokenPair["refresh_token"],
	})
}

// notifyUser sends an email to the user notifying them of a suspicious login attempt
func notifyUser(userGUID string) error {
	// Fetch user email from the database using GORM
	var user User
	result := db.Database.Where("guid = ?", userGUID).First(&user)
	if result.Error != nil {
		return fmt.Errorf("user not found: %v", result.Error)
	}

	// Use gomail to send an email
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@yourapp.com")
	m.SetHeader("To", user.email)
	m.SetHeader("Subject", "Suspicious Login Attempt")
	m.SetBody("text/html", fmt.Sprintf("Dear user,<br><br>We detected a login attempt from an unknown IP address.<br>If this was not you, please secure your account immediately."))

	// Set up the email dialer
	d := gomail.NewDialer("smtp.example.com", 587, "username", "password")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func main() {
	r := gin.Default()
	db.Connect()
	r.GET("/ping", handleReturnPong)
	r.GET("/:guid", HandleGetTokens)
	// Endpoint to refresh token
	r.POST("/:guid/refresh-token", func(c *gin.Context) {
		var request struct {
			RefreshToken string `json:"refresh_token"`
		}

		// Bind the JSON request body to the struct
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Assume userGUID is retrieved from session or auth context (hardcoded for example)
		userGUID := c.Param("guid")

		// Call the function to check IP and refresh tokens
		HandleRefreshTokens(c, userGUID, request.RefreshToken)
	})
	r.Run(":9000")
}
