package middleware

import (
	"GolandRestApi/pkg/config"
	"GolandRestApi/pkg/service"
	"GolandRestApi/pkg/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// Authenticate is a middleware function that enforces token-based authentication for all incoming requests
// except for specific endpoints like /login, /register, and /refreshToken. It verifies the validity of the access token
// provided in the Authorization header and ensures that it is correctly formatted as a Bearer token.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// cfg: A pointer to the config.Config struct that contains the API version and other configuration details.
//
// Returns a http.Handler that performs authentication checks before passing control to the next handler.
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
				//http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				service.HttpErrorResponse(logger,
					w,
					http.StatusUnauthorized,
					"/api/"+cfg.APIVersion,
					"Authorization header is required",
					nil,
					utils.LogTypeWarn,
					"not able to get the username")
				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				//http.Error(w, "Invalid token format", http.StatusUnauthorized)
				service.HttpErrorResponse(logger,
					w,
					http.StatusUnauthorized,
					"/api/"+cfg.APIVersion,
					"Invalid token format",
					nil,
					utils.LogTypeWarn,
					"not able to get the username")
				return
			}

			tokenString := bearerToken[1]

			// Verify token
			username, _, err := service.ExtractClaimsFromToken(logger, tokenString)
			if err != nil {
				//http.Error(w, "Invalid token", http.StatusUnauthorized)
				service.HttpErrorResponse(logger,
					w,
					http.StatusUnauthorized,
					"/api/"+cfg.APIVersion,
					"Invalid token",
					nil,
					utils.LogTypeWarn,
					"not able to get the username")
				return
			}

			err = service.VerifyToken(logger, cfg, username, tokenString)
			if err != nil {
				//http.Error(w, "Invalid token", http.StatusUnauthorized)
				service.HttpErrorResponse(logger,
					w,
					http.StatusUnauthorized,
					"/api/"+cfg.APIVersion,
					"Invalid token",
					nil,
					utils.LogTypeWarn,
					username)
				return
			}

			// Token is valid, proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
