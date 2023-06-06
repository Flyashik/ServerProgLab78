package api

import (
	"ApiService/internal/app/model"
	"ApiService/internal/app/storage"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

var signingKey = []byte("very_secret_key")
var refreshSigningKey = []byte("refresh_secret_key")

type Token struct {
	AccessToken  string
	RefreshToken string
}

func (s *Server) generateAccessToken(user model.User) string {
	claims := &model.Claims{
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.FormatInt(user.Id, 10),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		s.logger.Error(err)
		return ""
	}

	return tokenString
}

func (s *Server) generateRefreshToken(user model.User) string {
	claims := &model.Claims{
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.FormatInt(user.Id, 10),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		s.logger.Error(err)
		return ""
	}

	return tokenString
}

func verifyRefreshToken(refreshToken string) (*model.Claims, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return refreshSigningKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	return claims, nil
}

func (s *Server) refreshAccessToken(claims model.Claims) (string, error) {
	userId, err := strconv.Atoi(claims.StandardClaims.Subject)
	if err != nil {
		return "", err
	}

	var user model.User
	row, err := s.storage.Query(storage.GetUserByID, userId)
	if err != nil {
		return "", err
	}
	err = row.Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.Password,
		&user.Role)
	if err != nil {
		return "", err
	}

	token := s.generateAccessToken(user)

	return token, nil
}
