package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type emailCtx struct{}

type jwtClaims struct {
	Email string `json:"email"`
	AppID string `json:"app_id"`
	jwt.RegisteredClaims
}

func parseJWT(tokenStr string, secret []byte, wantApp string) (*jwtClaims, error) {
	c := &jwtClaims{}
	tok, err := jwt.ParseWithClaims(tokenStr, c, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	}, jwt.WithExpirationRequired(), jwt.WithLeeway(30*time.Second))
	if err != nil || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	if c.AppID != wantApp {
		return nil, errors.New("wrong app_id")
	}
	return c, nil
}

func (s *server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
		if tokenStr == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
			return
		}
		claims, err := parseJWT(tokenStr, s.secret, s.appID)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
			return
		}
		ctx := context.WithValue(r.Context(), emailCtx{}, claims.Email)
		next(w, r.WithContext(ctx))
	}
}
