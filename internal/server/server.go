package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/server/handlers/private"
	"go-gorilla-autos/internal/server/handlers/public"
)

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	newServer := &Server{
		port: port,
		db:   database.New(),
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(s.corsMiddleware)

	// Ruta pública
	publicRouter := r.PathPrefix("/api").Subrouter()
	publicRouter.HandleFunc("/autos", func(w http.ResponseWriter, r *http.Request) {
		public.GetAutosHandler(w, r, s.db)
	}).Methods("GET")

	// Ruta privada
	privateRouter := r.PathPrefix("/api/admin").Subrouter()
	privateRouter.Use(s.authMiddleware)
	privateRouter.HandleFunc("/autos", func(w http.ResponseWriter, r *http.Request) {
		private.CreateAutoHandler(w, r, s.db)
	}).Methods("POST")

	return r
}

// Middlewares
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "false")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
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
