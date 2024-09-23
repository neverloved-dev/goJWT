package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTokenPair(t *testing.T) {
	// Arrange
	username := "testuser"
	ipaddress := "127.0.0.1"
	userguid := "myUserGUID"

	// ACT
	tokenPair, err := generateTokenPair(username, ipaddress, userguid)

	// Assert
	assert.NoError(t, err, "Expected no error while generating token pair")
	assert.NotEmpty(t, tokenPair["access_token"], "Expected access token to be generated")
	assert.NotEmpty(t, tokenPair["refresh_token"], "Expected refresh token to be generated")
}

func TestAccessTokenClaims(t *testing.T) {
	username := "testuser"
	ipaddress := "127.0.0.1"
	userguid := "myUserGUID"

	// Arrange
	tokenPair, _ := generateTokenPair(username, ipaddress, userguid)

	// ACT
	accessToken, err := parseToken(tokenPair["access_token"], userguid, t)
	//ASSERT
	assert.NoError(t, err, "Expected no error while parsing the access token")
	if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
		assert.Equal(t, username, claims["name"], "Expected username in access token to match")
		assert.Equal(t, ipaddress, claims["ipaddress"], "Expected IP address in access token to match")

		// Check expiration time
		exp := claims["exp"].(float64)
		assert.True(t, int64(exp) > time.Now().Unix(), "Expected access token to be valid (not expired)")
	} else {
		t.Error("Invalid access token claims")
	}
}

func TestRefreshTokenClaims(t *testing.T) {
	//ARRANGE
	username := "testuser"
	userguid := "myUserGUID"
	ipaddress := "127.0.0.1"

	// ACT
	tokenPair, _ := generateTokenPair(username, ipaddress, userguid)

	// ASSERT
	refreshToken, err := parseToken(tokenPair["refresh_token"], userguid, t)

	assert.NoError(t, err, "Expected no error while parsing the refresh token")

	// Verify refresh token claims
	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		assert.Equal(t, username, claims["username"], "Expected username in refresh token to match")
		// Check expiration time
		exp := claims["exp"].(float64)
		assert.True(t, int64(exp) > time.Now().Unix(), "Expected refresh token to be valid (not expired)")
	} else {
		t.Error("Invalid refresh token claims")
	}
}

// parseToken is a helper function to parse and validate JWT tokens
func parseToken(tokenString, userguid string, t *testing.T) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(userguid), nil
	})

	return token, err
}
