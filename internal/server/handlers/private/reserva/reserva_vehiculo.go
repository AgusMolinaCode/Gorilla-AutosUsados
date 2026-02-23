package reserva

import (
	"encoding/json"
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"
	publicReserva "go-gorilla-autos/internal/server/handlers/public/reserva"

	"github.com/gorilla/mux"
)

// writeJSONResponse escribe una respuesta JSON con el código de estado especificado
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Si falla la codificación, intentar escribir error simple
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// CrearReservaHandler maneja la creación de una nueva reserva
func CrearReservaHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	stockID := vars["stock_id"]

	var reserva models.Reserva
	if err := json.NewDecoder(r.Body).Decode(&reserva); err != nil {
		publicReserva.WriteErrorResponse(w, http.StatusBadRequest, "Error al decodificar el JSON")
		return
	}

	// Generar un ID aleatorio para la reserva
	reserva.ID = models.GenerarIDReserva()

	// Validaciones básicas
	if reserva.Nombre == "" || reserva.Apellido == "" {
		publicReserva.WriteErrorResponse(w, http.StatusBadRequest, "Nombre y apellido son requeridos")
		return
	}

	if reserva.FechaHora.IsZero() {
		publicReserva.WriteErrorResponse(w, http.StatusBadRequest, "Fecha y hora de reserva son requeridas")
		return
	}

	// Usar helper del paquete público para buscar auto
	result := publicReserva.FindAutoByStockID(r.Context(), db, stockID)
	if !result.Found {
		publicReserva.WriteNotFoundResponse(w, "Auto no encontrado")
		return
	}
	auto := result.Auto

	// Agregar la nueva reserva
	auto.Reservas = append(auto.Reservas, reserva)

	// Usar helper del paquete público para actualizar
	if err := publicReserva.UpdateAutoReservas(r.Context(), db, stockID, auto.Reservas); err != nil {
		publicReserva.WriteErrorResponse(w, http.StatusInternalServerError, "Error al actualizar las reservas")
		return
	}

	response := map[string]interface{}{
		"mensaje":        "Reserva creada exitosamente",
		"reserva":        reserva,
		"total_reservas": len(auto.Reservas),
	}
	writeJSONResponse(w, http.StatusCreated, response)
}

func EliminarReservaHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	stockID := vars["stock_id"]
	reservaID := vars["reserva_id"]

	// Usar helper del paquete público para buscar auto
	result := publicReserva.FindAutoByStockID(r.Context(), db, stockID)
	if !result.Found {
		publicReserva.WriteNotFoundResponse(w, "Auto no encontrado")
		return
	}
	auto := result.Auto

	// Buscar la reserva por ID
	var indexToDelete int
	found := false
	for i := range auto.Reservas {
		if auto.Reservas[i].ID == reservaID {
			indexToDelete = i
			found = true
			break
		}
	}

	if !found {
		publicReserva.WriteNotFoundResponse(w, "Reserva no encontrada")
		return
	}

	// Eliminar la reserva
	auto.Reservas = append(auto.Reservas[:indexToDelete], auto.Reservas[indexToDelete+1:]...)

	// Usar helper del paquete público para actualizar
	if err := publicReserva.UpdateAutoReservas(r.Context(), db, stockID, auto.Reservas); err != nil {
		publicReserva.WriteErrorResponse(w, http.StatusInternalServerError, "Error al actualizar las reservas")
		return
	}

	response := map[string]interface{}{
		"mensaje": "Reserva eliminada exitosamente",
	}
	writeJSONResponse(w, http.StatusOK, response)
}

func EditarReservaHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	stockID := vars["stock_id"]
	reservaID := vars["reserva_id"]

	var nuevaReserva models.Reserva
	if err := json.NewDecoder(r.Body).Decode(&nuevaReserva); err != nil {
		publicReserva.WriteErrorResponse(w, http.StatusBadRequest, "Error al decodificar el JSON")
		return
	}

	// Usar helper del paquete público para buscar auto
	result := publicReserva.FindAutoByStockID(r.Context(), db, stockID)
	if !result.Found {
		publicReserva.WriteNotFoundResponse(w, "Auto no encontrado")
		return
	}
	auto := result.Auto

	// Buscar la reserva por ID
	var reservaEncontrada *models.Reserva
	for i := range auto.Reservas {
		if auto.Reservas[i].ID == reservaID {
			reservaEncontrada = &auto.Reservas[i]
			break
		}
	}

	if reservaEncontrada == nil {
		publicReserva.WriteNotFoundResponse(w, "Reserva no encontrada")
		return
	}

	// Actualizar la reserva
	reservaEncontrada.Nombre = nuevaReserva.Nombre
	reservaEncontrada.Apellido = nuevaReserva.Apellido
	reservaEncontrada.Telefono = nuevaReserva.Telefono
	reservaEncontrada.Comentario = nuevaReserva.Comentario
	reservaEncontrada.FechaHora = nuevaReserva.FechaHora

	// Usar helper del paquete público para actualizar
	if err := publicReserva.UpdateAutoReservas(r.Context(), db, stockID, auto.Reservas); err != nil {
		publicReserva.WriteErrorResponse(w, http.StatusInternalServerError, "Error al actualizar las reservas")
		return
	}

	response := map[string]interface{}{
		"mensaje": "Reserva actualizada exitosamente",
		"reserva": reservaEncontrada,
	}
	writeJSONResponse(w, http.StatusOK, response)
}

// ObtenerReservasHandler obtiene todas las reservaciones de un auto
func ObtenerReservasHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener stock_id
	vars := mux.Vars(r)
	stockID := vars["stock_id"]

	// Validar formato de stock_id
	if err := models.ValidateStockID(stockID); err != nil {
		publicReserva.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Usar helper del paquete público para buscar auto
	result := publicReserva.FindAutoByStockID(r.Context(), db, stockID)
	if !result.Found {
		publicReserva.WriteNotFoundResponse(w, "Auto no encontrado")
		return
	}
	auto := result.Auto

	// Verificar si hay reservas
	if len(auto.Reservas) == 0 {
		response := map[string]interface{}{
			"mensaje":  "No hay reservas para este auto",
			"reservas": []models.Reserva{},
			"stock_id": stockID,
		}
		writeJSONResponse(w, http.StatusOK, response)
		return
	}

	// Preparar respuesta con reservas
	response := map[string]interface{}{
		"mensaje":  "Reservas obtenidas exitosamente",
		"reservas": auto.Reservas,
		"stock_id": stockID,
		"total":    len(auto.Reservas),
	}

	writeJSONResponse(w, http.StatusOK, response)
}
