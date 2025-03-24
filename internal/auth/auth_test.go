package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"

	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mysecretpassword"

	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	err = CheckPasswordHash(password, hashedPassword)
	assert.NoError(t, err)
}

func TestCheckPasswordHash_InvalidPassword(t *testing.T) {
	password := "mysecretpassword"
	invalidPassword := "wrongpassword"

	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	err = CheckPasswordHash(invalidPassword, hashedPassword)
	assert.Error(t, err)
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	validatedUserID, err := ValidateJWT(token, tokenSecret)
	assert.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	tokenSecret := "mysecret"
	invalidToken := "invalidtoken"

	_, err := ValidateJWT(invalidToken, tokenSecret)
	assert.Error(t, err)
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecret"
	expiresIn := -time.Hour // Token already expired

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	_, err = ValidateJWT(token, tokenSecret)
	assert.Error(t, err)
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer validtoken")

	token, err := GetBearerToken(headers)
	assert.NoError(t, err)
	assert.Equal(t, "validtoken", token)
}

func TestGetBearerToken_MissingAuthorizationHeader(t *testing.T) {
	headers := http.Header{}

	_, err := GetBearerToken(headers)
	assert.Error(t, err)
	assert.Equal(t, "authorization header missing", err.Error())
}

func TestGetBearerToken_WrongAuthorizationMethod(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Basic invalidtoken")

	_, err := GetBearerToken(headers)
	assert.Error(t, err)
	assert.Equal(t, "wrong authorization method", err.Error())
}

func TestGetBearerToken_MissingToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer ")

	_, err := GetBearerToken(headers)
	assert.Error(t, err)
	assert.Equal(t, "token missing", err.Error())
}

func TestMakeRefreshToken(t *testing.T) {
	token, err := MakeRefreshToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Len(t, token, 64) // 32 bytes * 2 (hex encoding)
}

func TestGetAPIKey(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "ApiKey validapikey")

	apiKey, err := GetAPIKey(headers)
	assert.NoError(t, err)
	assert.Equal(t, "validapikey", apiKey)
}

func TestGetAPIKey_MissingAuthorizationHeader(t *testing.T) {
	headers := http.Header{}

	_, err := GetAPIKey(headers)
	assert.Error(t, err)
	assert.Equal(t, "authorization header missing", err.Error())
}

func TestGetAPIKey_WrongAuthorizationMethod(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer invalidapikey")

	_, err := GetAPIKey(headers)
	assert.Error(t, err)
	assert.Equal(t, "wrong authorization method", err.Error())
}

func TestGetAPIKey_MissingAPIKey(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "ApiKey ")

	_, err := GetAPIKey(headers)
	assert.Error(t, err)
	assert.Equal(t, "api key missing", err.Error())
}
