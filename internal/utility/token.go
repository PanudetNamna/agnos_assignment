package utility

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(staffID, hospitalID uuid.UUID, secret string) (string, error) {
	claims := jwt.MapClaims{
		"staff_id":    staffID.String(),
		"hospital_id": hospitalID.String(),
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
