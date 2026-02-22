package reserva

import (
	"encoding/json"
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"

	"github.com/gorilla/mux"
)

// CrearReservaHandler maneja la creación de una nueva reserva
func CrearReservaHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	stockID := vars["stock_id"]

	var reserva models.Reserva
	if err := json.NewDecoder(r.Body).Decode(&reserva); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Error al decodificar el JSON")
		return
	}

	// Generar un ID aleatorio para la reserva
	reserva.ID = models.GenerarIDReserva()

	// Validaciones básicas
	if reserva.Nombre == "" || reserva.Apellido == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Nombre y apellido son requeridos")
		return
	}

	if reserva.FechaHora.IsZero() {
		WriteErrorResponse(w, http.StatusBadRequest, "Fecha y hora de reserva son requeridas")
		return
	}

	// Buscar el auto usando el helper
	result := FindAutoByStockID(r.Context(), db, stockID)
	if !result.Found {
		WriteNotFoundResponse(w, "Auto no encontrado")
		return
	}
	auto := result.Auto

	// Verificar si el cliente ya tiene una reserva para este auto
	for _, existingReserva := range auto.Reservas {
		if existingReserva.Nombre == reserva.Nombre && existingReserva.Apellido == reserva.Apellido {
			WriteErrorResponse(w, http.StatusConflict, "Ya existe una reserva activa para este cliente y vehículo")
			return
		}
	}

	// Agregar la nueva reserva
	auto.Reservas = append(auto.Reservas, reserva)

	// Actualizar usando el helper
	if err := UpdateAutoReservas(r.Context(), db, stockID, auto.Reservas); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Error al actualizar las reservas")
		return
	}

	response := map[string]interface{}{
		"mensaje":        "Reserva creada exitosamente",
		"reserva":        reserva,
		"total_reservas": len(auto.Reservas),
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
