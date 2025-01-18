package middleware

import (
	"net/http"
	"os"
	"strings"
)

func AuthMiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir OPTIONS para CORS
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// Obtener el token de autorización
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Se requiere autorización", http.StatusUnauthorized)
			return
		}

		// Extraer el token (Bearer Token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato de autorización inválido", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		expectedToken := os.Getenv("AUTH_KEY")

		if token != expectedToken {
			http.Error(w, "Código de autorización inválido", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
