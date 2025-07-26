
package utils

import (
	"go-web/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTFunctions(t *testing.T) {
	// 1. Setup
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "a-very-secret-key-for-testing",
			Expiration: 60, // 60 seconds
		},
	}
	userID := uint(123)
	role := "admin"

	// 2. Test Token Generation
	tokenString, err := GenerateToken(userID, role, cfg)

	assert.NoError(t, err, "Token generation should not produce an error")
	assert.NotEmpty(t, tokenString, "Generated token string should not be empty")

	// 3. Test Token Parsing (Happy Path)
	claims, err := ParseToken(tokenString, cfg.JWT.Secret)

	assert.NoError(t, err, "Token parsing should not produce an error for a valid token")
	assert.NotNil(t, claims, "Claims should not be nil for a valid token")
	assert.Equal(t, userID, claims.UserID, "UserID in claims should match the original UserID")
	assert.Equal(t, role, claims.Role, "Role in claims should match the original role")

	// 4. Test Token Parsing (Invalid Token)
	invalidTokenString := "this-is-not-a-valid-jwt"
	claims, err = ParseToken(invalidTokenString, cfg.JWT.Secret)

	assert.Error(t, err, "Parsing an invalid token should produce an error")
	assert.Nil(t, claims, "Claims should be nil for an invalid token")

	// 5. Test Token Parsing (Expired Token)
	cfgExpired := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "another-secret",
			Expiration: -1, // Already expired
		},
	}
	expiredToken, _ := GenerateToken(userID, role, cfgExpired)
	// Wait a moment to ensure the timestamp is in the past
	time.Sleep(1 * time.Second)
	claims, err = ParseToken(expiredToken, cfgExpired.JWT.Secret)

	assert.Error(t, err, "Parsing an expired token should produce an error")
	assert.Contains(t, err.Error(), "token is expired", "Error message should indicate token expiration")
}
