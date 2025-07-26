
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	password := "my-super-secret-password"

	// 1. Test hashing the password
	hashedPassword, err := HashPassword(password)

	assert.NoError(t, err, "Hashing the password should not produce an error")
	assert.NotEmpty(t, hashedPassword, "The hashed password should not be empty")
	assert.NotEqual(t, password, hashedPassword, "The hashed password should not be the same as the original password")

	// 2. Test checking the correct password
	err = CheckPasswordHash(password, hashedPassword)
	assert.NoError(t, err, "Checking the correct password should not produce an error")

	// 3. Test checking an incorrect password
	wrongPassword := "not-my-password"
	err = CheckPasswordHash(wrongPassword, hashedPassword)
	assert.Error(t, err, "Checking the wrong password should produce an error")
}
