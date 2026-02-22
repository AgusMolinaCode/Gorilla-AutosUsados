package descuentos

import (
	"context"
	"encoding/json"
	"net/http"

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

// AplicarDescuentoHandler aplica un descuento a un auto
func AplicarDescuentoHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener stock_id
	vars := mux.Vars(r)
	stockID := vars["stock_id"]

	// Validar formato de stock_id
	if err := models.ValidateStockID(stockID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decodificar el descuento
	var descuentoRequest struct {
		Descuento float64 `json:"descuento"`
	}
	if err := json.NewDecoder(r.Body).Decode(&descuentoRequest); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Validar descuento
	if descuentoRequest.Descuento < 0 {
		http.Error(w, "Descuento no puede ser negativo", http.StatusBadRequest)
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

	// Verificar si ya existe un descuento
	if auto.Descuento > 0 {
		http.Error(w, "Ya existe un descuento para este auto. Elimínelo primero para agregar uno nuevo", http.StatusBadRequest)
		return
	}

	// Validar que el descuento no supere el precio original
	if descuentoRequest.Descuento > auto.Precio {
		http.Error(w, "El descuento no puede ser mayor al precio original", http.StatusBadRequest)
		return
	}

	// Calcular precio con descuento
	precioOriginal := auto.Precio
	precioConDescuento := precioOriginal - descuentoRequest.Descuento

	// Validar que el precio con descuento sea positivo
	if precioConDescuento <= 0 {
		http.Error(w, "El descuento no puede reducir el precio a cero o negativo", http.StatusBadRequest)
		return
	}

	// Actualizar el descuento y el precio
	update := bson.M{
		"$set": bson.M{
			"descuento": descuentoRequest.Descuento,
			"precio":    precioConDescuento,
		},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error al aplicar el descuento", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"mensaje":              "Descuento aplicado exitosamente",
		"descuento":            descuentoRequest.Descuento,
		"precio_original":      precioOriginal,
		"precio_con_descuento": precioConDescuento,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// EliminarDescuentoHandler elimina el descuento de un auto
func EliminarDescuentoHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener stock_id
	vars := mux.Vars(r)
	stockID := vars["stock_id"]

	// Validar formato de stock_id
	if err := models.ValidateStockID(stockID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	// Verificar si hay descuento para eliminar
	if auto.Descuento == 0 {
		http.Error(w, "No hay descuento aplicado para este auto", http.StatusBadRequest)
		return
	}

	// Calcular precio original (sumando el descuento)
	precioOriginal := auto.Precio + auto.Descuento

	// Actualizar eliminando el descuento
	update := bson.M{
		"$set": bson.M{
			"descuento": 0,
			"precio":    precioOriginal,
		},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error al eliminar el descuento", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"mensaje":         "Descuento eliminado exitosamente",
		"precio_original": precioOriginal,
		"descuento":       0,
	}

	writeJSONResponse(w, http.StatusOK, response)
}
