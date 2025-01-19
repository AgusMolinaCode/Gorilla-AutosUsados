package public

import (
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/server/handlers/public"
	"go-gorilla-autos/internal/server/handlers/public/reserva"

	"github.com/gorilla/mux"
)

func RegisterPublicRoutes(r *mux.Router, db database.Service) {
	publicRouter := r.PathPrefix("/api").Subrouter()
	publicRouter.HandleFunc("/autos", func(w http.ResponseWriter, r *http.Request) {
		public.GetAutosHandler(w, r, db)
	}).Methods("GET")

	publicRouter.HandleFunc("/autos/destacados", func(w http.ResponseWriter, r *http.Request) {
		public.GetFeaturedAutosHandler(w, r, db)
	}).Methods("GET")

	publicRouter.HandleFunc("/autos/{stock_id}/reservations", func(w http.ResponseWriter, r *http.Request) {
		reserva.CrearReservaHandler(w, r, db)
	}).Methods("POST")
}
