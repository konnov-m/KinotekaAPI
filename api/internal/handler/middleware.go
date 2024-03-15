package handler

import (
	"KinotekaAPI/internal/service"
	"context"
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

func userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := service.ParseToken(bearerToken[1])
		if err != nil {
			newErrorResponse(w, err, "Can't parse token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
