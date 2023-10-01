package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hash 与えられた文字列をbcryptでハッシュ化するだけ
func Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// Verify 与えられた文字列パスワードがハッシュ化されたパスワードと一致するかをチェックするだけ
func Verify(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
