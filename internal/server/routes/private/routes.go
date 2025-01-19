package private

import (
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/server/handlers/private"
	"go-gorilla-autos/internal/server/handlers/private/descuentos"
	"go-gorilla-autos/internal/server/handlers/private/destacado"
	"go-gorilla-autos/internal/server/handlers/private/estado"
	"go-gorilla-autos/internal/server/handlers/private/reserva"

	"github.com/gorilla/mux"
)

func RegisterPrivateRoutes(privateRouter *mux.Router, db database.Service) {
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
		destacado.ToggleFeaturedHandler(w, r, db)
	}).Methods("POST")

	privateRouter.HandleFunc("/autos/{stock_id}/status", func(w http.ResponseWriter, r *http.Request) {
		estado.CambiarEstadoAutoHandler(w, r, db)
	}).Methods("POST")

	privateRouter.HandleFunc("/autos/{stock_id}/discount", func(w http.ResponseWriter, r *http.Request) {
		descuentos.AplicarDescuentoHandler(w, r, db)
	}).Methods("POST")

	privateRouter.HandleFunc("/autos/{stock_id}/discount", func(w http.ResponseWriter, r *http.Request) {
		descuentos.EliminarDescuentoHandler(w, r, db)
	}).Methods("DELETE")

	// Ruta para obtener reservas de un auto
	privateRouter.HandleFunc("/autos/{stock_id}/reservations", func(w http.ResponseWriter, r *http.Request) {
		reserva.ObtenerReservasHandler(w, r, db)
	}).Methods("GET")

	privateRouter.HandleFunc("/autos/{stock_id}/reservations", func(w http.ResponseWriter, r *http.Request) {
		reserva.CrearReservaHandler(w, r, db)
	}).Methods("POST")

	privateRouter.HandleFunc("/autos/{stock_id}/reservations/{reserva_id}", func(w http.ResponseWriter, r *http.Request) {
		reserva.EditarReservaHandler(w, r, db)
	}).Methods("PUT")

	privateRouter.HandleFunc("/autos/{stock_id}/reservations/{reserva_id}", func(w http.ResponseWriter, r *http.Request) {
		reserva.EliminarReservaHandler(w, r, db)
	}).Methods("DELETE")

}
