package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Subject: userID.String(),
	});

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	});

	if err != nil {
		return uuid.Nil, err;
	}

	if !token.Valid {
		return uuid.Nil, errors.New("the token is invalid");
	}

	id, err := token.Claims.GetSubject()

	if err != nil {
		return uuid.Nil, err;
	}

	userUuid, err := uuid.Parse(id);

	if err != nil {
		return uuid.Nil, err;
	}

	return userUuid, nil;
}

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization");

	if bearerToken == "" {
		return "", errors.New("no authorization header found");
	}

	return strings.Replace(bearerToken, "Bearer ", "", 1), nil;
}

func MakeRefreshToken() (string, error) {
	refreshToken := make([]byte, 32);
	_, err := rand.Read(refreshToken);

	if err != nil {
		return "", err;
	}

	return hex.EncodeToString(refreshToken), nil;
}