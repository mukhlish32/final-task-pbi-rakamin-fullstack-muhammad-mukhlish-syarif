package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Fungsi untuk generates JWT token berdasarkan userID
func GenerateToken(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Fungsi yang digunakan untuk melakukan validasi pada password
func ValidatePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Fungsi untuk mengekstrak token JWT dari authorization header
func GetTokenFromHeader(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

// Fungsi untuk mengekstrak user ID dari JWT token
func GetUserIDFromToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}

	// Extract user ID dari claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token type")
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid token type")
	}
	return uint(userID), nil
}

// Fungsi untuk mengecek apakah user memiliki izin untuk melakukan perubahan pada foto
func CheckPhotoPermission(c *gin.Context, photoUserID uint) error {
	currentUserID, err := GetUserIDFromToken(GetTokenFromHeader(c.GetHeader("Authorization")))
	if err != nil {
		return err
	}

	if currentUserID != photoUserID {
		return errors.New("anda tidak memiliki izin pada foto user ini")
	}

	return nil
}
