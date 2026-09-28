package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"jobsite/helper"
	"jobsite/service"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type ErrorResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// JWTAuthMiddleware is an HTTP middleware that validates JWT tokens and (optionally) enforces a role.
func JWTAuthMiddleware(logger log.Logger, authSvc service.UserService, jwtSecret, requiredRole string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(ErrorResponse{Message: "missing Authorization header"})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid Authorization header"})
				return
			}

			tokenString := parts[1]
			claims, err := helper.ParseToken(tokenString, jwtSecret)
			if err != nil {
				level.Warn(logger).Log("msg", "invalid token", "err", err)
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid token"})
				return
			}

			if requiredRole != "" && string(claims.Role) != requiredRole {
				level.Warn(logger).Log("msg", "access denied", "role", claims.Role)
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(ErrorResponse{Message: "you don't have access for this operation"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
