package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/server/routes/private"
	"go-gorilla-autos/internal/server/routes/public"

	"go-gorilla-autos/internal/server/routes/middleware"

	"github.com/joho/godotenv"
)

func init() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env")
	}
}

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
	r.Use(middleware.CORSMiddleware)

	// Crear un subrouter para rutas privadas
	privateRouter := r.PathPrefix("/api/admin").Subrouter()

	// Aplicar middleware de autenticación solo a rutas privadas
	privateRouter.Use(middleware.AuthMiddlewareFunc)

	// Registrar rutas públicas
	public.RegisterPublicRoutes(r, s.db)

	// Registrar rutas privadas
	private.RegisterPrivateRoutes(privateRouter, s.db)

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
