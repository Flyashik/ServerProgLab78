package api

import (
	"ApiService/internal/app/model"
	"ApiService/internal/app/storage"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"net/http"
	"time"
)

const hashSalt = "fjgh1rsZ04xzdsaka"

func (s *Server) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Error(err)
			return
		}
		defer r.Body.Close()
		var user model.User
		err = json.Unmarshal(body, &user)
		if err != nil {
			s.logger.Error(err)
			return
		}
		hashPassword := generatePasswordHash(user.Password)
		rows, err := s.storage.Query(storage.GetUserAuth, user.Username, hashPassword)
		if err != nil {
			http.Error(w, "Incorrect e-mail or password", http.StatusUnauthorized)
			return
		}
		defer rows.Close()
		if rows.Next() {
			err = rows.Scan(
				&user.Id,
				&user.Name,
				&user.Username,
				&user.Password,
				&user.Role)
			if err != nil {
				s.logger.Error(err)
				http.Error(w, "Can't fetch user", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Incorrect e-mail or password", http.StatusUnauthorized)
			s.logger.Error("Incorrect e-mail or password")
			return
		}

		token := Token{
			AccessToken:  s.generateAccessToken(user),
			RefreshToken: s.generateRefreshToken(user),
		}

		tokenCookie := http.Cookie{
			Name:     "access_token",
			Value:    token.AccessToken,
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24),
		}
		http.SetCookie(w, &tokenCookie)

		refreshTokenCookie := http.Cookie{
			Name:     "refresh_token",
			Value:    token.RefreshToken,
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24 * 7),
		}
		http.SetCookie(w, &refreshTokenCookie)

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie := &http.Cookie{
			Name:    "access_token",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		}
		http.SetCookie(w, accessTokenCookie)

		refreshTokenCookie := &http.Cookie{
			Name:    "refresh_token",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		}
		http.SetCookie(w, refreshTokenCookie)

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)
		if role != "admin" {
			s.logger.Info("User hasn't access")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Error(err)
			return
		}
		defer r.Body.Close()
		var user model.User
		err = json.Unmarshal(body, &user)
		if err != nil {
			s.logger.Error(err)
			return
		}

		if user.Name == "" || user.Username == "" || user.Password == "" || user.Role == "" {
			w.WriteHeader(http.StatusBadRequest)
			s.logger.Error("Missing fields")
			return
		}

		hashPassword := generatePasswordHash(user.Password)

		_, err = s.storage.Exec(storage.CreateUser, user.Name, user.Username, hashPassword, user.Role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Error("Can't create user: ", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) IsAuthorized(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			s.logger.Error(err)
			return
		}

		tokenString := cookie.Value

		token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})
		if err != nil || !token.Valid {
			refreshTokenCookie, err := r.Cookie("refresh_token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				s.logger.Error(err)
				return
			}

			refreshToken := refreshTokenCookie.Value

			refreshClaims, err := verifyRefreshToken(refreshToken)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				s.logger.Error(err)
				return
			}

			newAccessToken, err := s.refreshAccessToken(*refreshClaims)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				s.logger.Error(err)
			}

			cookie.Value = newAccessToken
			http.SetCookie(w, cookie)
		}

		claims, ok := token.Claims.(*model.Claims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			s.logger.Error("Invalid auth token: here")
			return
		}

		ctx := context.WithValue(r.Context(), "role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(hashSalt)))
}
