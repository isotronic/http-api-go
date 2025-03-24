package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return "", err
	}
	return string(hashByte), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedJWT, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method used: %v", t.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %v", err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	elements := strings.Split(authHeader, " ")
	if strings.ToLower(elements[0]) != "bearer" {
		return "", fmt.Errorf("wrong authorization method")
	}
	if elements[1] == "" {
		return "", fmt.Errorf("token missing")
	}

	return elements[1], nil
}

func MakeRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	elements := strings.Split(authHeader, " ")
	if strings.ToLower(elements[0]) != "apikey" {
		return "", fmt.Errorf("wrong authorization method")
	}
	if elements[1] == "" {
		return "", fmt.Errorf("api key missing")
	}

	return elements[1], nil
}