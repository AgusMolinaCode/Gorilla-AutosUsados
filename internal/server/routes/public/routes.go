package public

import (
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/server/handlers/public"

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
}
