package middleware

import (
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Se requiere autorización", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != "123456" {
			http.Error(w, "Código de autorización inválido", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
