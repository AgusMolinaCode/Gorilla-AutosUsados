package estado

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// writeJSONResponse escribe una respuesta JSON con el código de estado especificado
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// TODO:"Disponible" "En negociación" "Reservado" "Vendido" "En mantenimiento"

// CambiarEstadoAutoHandler cambia el estado de un auto
func CambiarEstadoAutoHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener stock_id
	vars := mux.Vars(r)
	stockID := vars["stock_id"]

	// Validar formato de stock_id
	if err := models.ValidateStockID(stockID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decodificar el nuevo estado
	var estadoRequest struct {
		Estado          string                    `json:"estado"`
		ReservadoPor    *models.ReservadoInfo     `json:"reservado_por,omitempty"`
		VendidoPor      *models.VendidoInfo       `json:"vendido_por,omitempty"`
		EnNegociacion   *models.NegociacionInfo   `json:"en_negociacion,omitempty"`
		EnMantenimiento *models.MantenimientoInfo `json:"en_mantenimiento,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&estadoRequest); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Validar estado
	validStates := []string{"disponible", "reservado", "vendido", "en negociación", "en mantenimiento"}
	if !slices.Contains(validStates, estadoRequest.Estado) {
		http.Error(w, "Estado inválido", http.StatusBadRequest)
		return
	}

	collection := db.Collection("autos")

	// Buscar el auto por stock_id
	var auto models.Auto
	filter := bson.M{"stock_id": stockID}
	err := collection.FindOne(context.Background(), filter).Decode(&auto)
	if err != nil {
		http.Error(w, "Auto no encontrado", http.StatusNotFound)
		return
	}

	// Actualizar el estado y la información correspondiente
	update := bson.M{"$set": bson.M{"estado": estadoRequest.Estado}}

	switch estadoRequest.Estado {
	case "reservado":
		update["$set"].(bson.M)["reservado_por"] = estadoRequest.ReservadoPor
	case "vendido":
		update["$set"].(bson.M)["vendido_por"] = estadoRequest.VendidoPor
	case "en negociación":
		update["$set"].(bson.M)["en_negociacion"] = estadoRequest.EnNegociacion
	case "en mantenimiento":
		update["$set"].(bson.M)["en_mantenimiento"] = estadoRequest.EnMantenimiento
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error al actualizar el estado del auto", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"mensaje": "Estado del auto actualizado exitosamente",
		"estado":  estadoRequest.Estado,
	}

	writeJSONResponse(w, http.StatusOK, response)
}
