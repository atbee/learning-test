package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/codeforpublic/morchana-static-qr-code-api/internal/jsonw"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/labstack/echo"
)

type contextKey int

const (
	contextKeyClaims contextKey = iota
)

var (
	ErrTokenMalformed    = errors.New("token is mulformed")
	ErrTokenExpired      = errors.New("token is expired")
	ErrInvalidSignMethod = errors.New("invalid signing method")
)

func Protect(secret []byte) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearer := bearerAuthorization(r)
			if bearer == "" {
				jsonw.Unauthorized(w, errors.New("authorization token required"))
				return
			}

			claims, err := bearerClaims(secret, bearer)
			if err != nil {
				switch err {
				case ErrTokenExpired, ErrTokenMalformed:
					jsonw.Unauthorized(w, err)
				default:
					jsonw.InternalServerError(w, err)
				}

				return
			}

			handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyClaims, claims)))
		})
	}
}

func bearerAuthorization(r *http.Request) string {
	header := r.Header.Get(echo.HeaderAuthorization)

	if header == "" {
		return ""
	}

	splitedHeader := strings.Split(header, "Bearer ")
	if len(splitedHeader) != 2 {
		return ""
	}

	return splitedHeader[1]
}

func bearerClaims(secret []byte, tokenString string) (*jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

	// }

	return &token.Claims, err
}
