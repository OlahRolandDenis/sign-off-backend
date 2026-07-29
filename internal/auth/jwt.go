package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret string

func GenerateJWT(agencyID int64) (string, error) {
	claims := jwt.MapClaims{
		"agency_id": agencyID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(JWTSecret))
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)

	if !ok {
		return nil, fmt.Errorf("unexpected signing method")
	}

	return []byte(JWTSecret), nil
}

func ValidateJWT(tokenString string) (int64, error) {
	t, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return 0, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("error getting claims")
	}

	if !t.Valid {
		return 0, fmt.Errorf("error token invalid")
	}

	agencyID, ok := claims["agency_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid agency_id claim")
	}
	return int64(agencyID), nil
}
