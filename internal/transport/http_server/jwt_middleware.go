// the JWT auth middleware provides transparent HTTP command authentication using Bearer token
package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"internal/adapters/logger"
	"internal/app"
	"internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

// provides message authentication using JWT token
func HandleJWTAuth(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// check for corresponding JWT header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			logger.Error("jwt: authorization header is missing")
			http.Error(rw, "authorization header is missing", http.StatusUnauthorized)
			return
		}

		logger.Info("jwt: handling authorization")

		tokenString = tokenString[len("Bearer "):]

		token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				logger.Error("jwt: unexpected signing method")
				return nil, fmt.Errorf("jwt: unexpected signing method")
			}
			if app.Sc.SessionCryptoKey == nil {
				return nil, fmt.Errorf("jwt: session key is empty")
			}
			return &app.Sc.SessionCryptoKey.PublicKey, nil
		})

		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				logger.Error("jwt: malformed token")
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				logger.Error("jwt: invalid token signature")
			case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
				logger.Error("jwt: token time is not valid")
			default:
				logger.Error("jwt: error processing token: " + err.Error())
			}
			http.Error(rw, "invalid authorization token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			logger.Error("jwt: invalid authorization token")
			http.Error(rw, "invalid authorization token", http.StatusUnauthorized)
			return
		}

		//checking for username in custom claims
		//TODO: replace username with UID/hash/etc in JWT for better security
		//TODO: check for user ban/token revocation in production code
		claims, ok := token.Claims.(*domain.Claims)
		if !ok {
			logger.Error("jwt: unknown claims type, cannot proceed")
			http.Error(rw, "invalid authorization token", http.StatusUnauthorized)
			return
		}

		logger.Infof("jwt: user [%s] succesfully authorized", claims.Username)

		// extending context with authorized username value
		ctx := r.Context()
		ctx = context.WithValue(ctx, domain.CtxUsername, claims.Username)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
