package private

import (
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/server/handlers/private"
	"go-gorilla-autos/internal/server/routes/middleware"

	"github.com/gorilla/mux"
)

func RegisterPrivateRoutes(r *mux.Router, db database.Service) {
	privateRouter := r.PathPrefix("/api/admin").Subrouter()

	// Aplicar middleware de autenticación
	privateRouter.Use(middleware.AuthMiddleware)

	privateRouter.HandleFunc("/autos", func(w http.ResponseWriter, r *http.Request) {
		private.CreateAutoHandler(w, r, db)
	}).Methods("POST")

	privateRouter.HandleFunc("/autos/{stock_id}", func(w http.ResponseWriter, r *http.Request) {
		private.UpdateAutoHandler(w, r, db)
	}).Methods("PUT")

	privateRouter.HandleFunc("/autos/{stock_id}", func(w http.ResponseWriter, r *http.Request) {
		private.DeleteAutoHandler(w, r, db)
	}).Methods("DELETE")

	privateRouter.HandleFunc("/autos/{stock_id}/featured", func(w http.ResponseWriter, r *http.Request) {
		private.ToggleFeaturedHandler(w, r, db)
	}).Methods("POST")

	privateRouter.HandleFunc("/autos/{stock_id}/status", func(w http.ResponseWriter, r *http.Request) {
		private.CambiarEstadoAutoHandler(w, r, db)
	}).Methods("POST")
}
