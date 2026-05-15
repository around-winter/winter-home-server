// Package middleware 包含 HTTP 中间件
package middleware

import (
	"net/http"
	"os"
	"strings"
)

var validToken = os.Getenv("API_TOKEN")

func init() {
	if validToken == "" {
		validToken = "winter-secret-token-2024"
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "缺少认证信息", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "认证格式错误，请使用: Bearer <token>", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		if token != validToken {
			http.Error(w, "Token 无效", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
