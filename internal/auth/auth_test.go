package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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