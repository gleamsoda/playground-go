package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T) {
	password := "password"
	hashedPassword, err := Hash(password)
	assert.NoError(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err, "failed to verify hashed password")
}

func TestVerify(t *testing.T) {
	password := "password"
	hashedPassword, _ := Hash(password)

	err := Verify(password, string(hashedPassword))
	assert.NoError(t, err, "Expected password check to succeed")

	err = Verify("wrongpassword", string(hashedPassword))
	assert.Error(t, err, "Expected password check to fail")
}
