package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neverloved-dev/goJWT/db"

	"github.com/golang-jwt/jwt"
)

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

func HandleRefreshTokens(c *gin.Context) {
	guid := c.Param("guid")
}
func main() {
	r := gin.Default()
	db.Connect()
	r.GET("/ping", handleReturnPong)
	r.GET("/:guid", HandleGetTokens)
	r.POST("/:guid/refresh", HandleRefreshTokens)
	r.Run(":9000")
}
