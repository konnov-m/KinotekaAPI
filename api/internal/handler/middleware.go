package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
)

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request is \"%s\". Method is %s", r.URL, r.Method)
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			newErrorResponse(w, errors.New("empty auth header"), "empty auth header", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			newErrorResponse(w, errors.New("invalid auth header"), "invalid auth header", http.StatusUnauthorized)
			return
		}

		if len(bearerToken[1]) == 0 {
			newErrorResponse(w, errors.New("token is empty"), "token is empty", http.StatusUnauthorized)
			return
		}

		userId, err := h.ser.User.ParseToken(bearerToken[1])
		if err != nil {
			newErrorResponse(w, err, "Can't parse token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
