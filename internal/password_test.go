package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err, "failed to verify hashed password")
}

func TestCheckPassword(t *testing.T) {
	password := "password"
	hashedPassword, _ := HashPassword(password)

	err := CheckPassword(password, string(hashedPassword))
	assert.NoError(t, err, "Expected password check to succeed")

	err = CheckPassword("wrongpassword", string(hashedPassword))
	assert.Error(t, err, "Expected password check to fail")
}
