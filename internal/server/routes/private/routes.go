package private

import (
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/server/handlers/private"

	"github.com/gorilla/mux"
)

func RegisterPrivateRoutes(r *mux.Router, db database.Service) {
	privateRouter := r.PathPrefix("/api/admin").Subrouter()
	privateRouter.HandleFunc("/autos", func(w http.ResponseWriter, r *http.Request) {
		private.CreateAutoHandler(w, r, db)
	}).Methods("POST")

	privateRouter.HandleFunc("/autos/{id}", func(w http.ResponseWriter, r *http.Request) {
		private.UpdateAutoHandler(w, r, db)
	}).Methods("PUT")

	privateRouter.HandleFunc("/autos/{id}", func(w http.ResponseWriter, r *http.Request) {
		private.DeleteAutoHandler(w, r, db)
	}).Methods("DELETE")

	privateRouter.HandleFunc("/autos/{stock_id}/destacado", func(w http.ResponseWriter, r *http.Request) {
		private.ToggleFeaturedHandler(w, r, db)
	}).Methods("POST")
}
