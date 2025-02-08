package middleware

import (
	"auth_backend/domain"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const claimsContextKey = contextKey("claims")

func JwtMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Отсутствует заголовок Authorization", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Неверный формат заголовка Authorization", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Неверный или просроченный токен", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Неверные данные токена", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetClaims(r.Context())
		if !ok {
			http.Error(w, "Нет информации о пользователе", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(float64)
		if !ok || domain.Role(int(role)) != domain.AdminRole {
			http.Error(w, "Доступ разрешён только администраторам", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetClaims(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(jwt.MapClaims)
	return claims, ok
}
