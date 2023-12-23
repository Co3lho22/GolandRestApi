package middleware

import (
	"GolandRestApi/pkg/config"
	"GolandRestApi/pkg/service"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// TODO: Add the logic to obligate the access token in every request
// unless the endpoint is /login, /register, /refreshToken
func Authenticate(logger *logrus.Logger, cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip middleware for specific endpoints
			if r.URL.Path == "/api/"+cfg.APIVersion+"/user/login" ||
				r.URL.Path == "/api/"+cfg.APIVersion+"/user/register" ||
				r.URL.Path == "/api/"+cfg.APIVersion+"/token/refresh" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			tokenString := bearerToken[1]

			// Verify token
			username, _, err := service.ExtractClaimsFromToken(logger, tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			err = service.VerifyToken(logger, cfg, username, tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Token is valid, proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
