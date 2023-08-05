package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 与えられた文字列をbcryptでハッシュ化するだけ
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword 与えられた文字列パスワードがハッシュ化されたパスワードと一致するかをチェックするだけ
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
