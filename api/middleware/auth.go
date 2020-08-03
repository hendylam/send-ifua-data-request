package middleware

import (
	"context"
	"net/http"
	"strings"

	"tapera/util/constant"

	jwt "github.com/dgrijalva/jwt-go"
	"tapera.integrasi/api/util/httpx"
	"tapera.integrasi/api/util/identity"
)

// Auth func
func Auth(signingKey string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ext := httpx.New(w, r)
			token := ext.Header("Authorization")

			// bearear check
			b := "Bearer "
			if len(token) == 0 || !strings.Contains(token, b) {
				ext.JSONerr(http.StatusForbidden, "Your request is not authorized")
				return
			}

			t := strings.Split(token, b)
			if len(t) < 2 {
				ext.JSONerr(http.StatusForbidden, "Your request is not authorized")
				return
			}

			// Validate token
			var err error
			valid, err := jwt.Parse(t[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(signingKey), nil
			})

			if err != nil {
				ext.JSONerr(http.StatusForbidden, "Invalid authorization token")
				return
			}

			// https://github.com/dgrijalva/jwt-go/issues/186
			claims := valid.Claims.(jwt.MapClaims)
			userID := claims["id"].(string)
			name := claims["name"].(string)
			email := claims["email"].(string)
			role := int64(claims["role"].(float64))
			status := int64(claims["status"].(float64))

			user := identity.NewUser(userID, name, email, int(role), int(status))
			ctx := context.WithValue(r.Context(), constant.UserKey, user)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
