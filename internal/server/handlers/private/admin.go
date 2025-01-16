package private

import (
	"context"
	"encoding/json"
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateAutoHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	var auto models.Auto
	if err := json.NewDecoder(r.Body).Decode(&auto); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if auto.Marca == "" || auto.Modelo == "" || auto.TipoVenta == "" ||
		auto.Año == 0 || auto.Kilometraje == 0 || auto.Precio == 0 ||
		auto.Ciudad == "" || auto.Transmision == "" || auto.Traccion == "" ||
		auto.StockID == "" || auto.Sucursal == "" || auto.Version == "" ||
		auto.Garantia == "" {
		http.Error(w, "Faltan campos requeridos", http.StatusBadRequest)
		return
	}

	// Validar stock_id
	if err := models.ValidateStockID(auto.StockID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := db.Collection("autos")
	_, err := collection.InsertOne(context.Background(), auto)
	if err != nil {
		http.Error(w, "Error al guardar el auto", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"mensaje": "Auto creado exitosamente",
		"auto":    auto,
	}

	json.NewEncoder(w).Encode(response)
}

// UpdateAutoHandler actualiza un auto existente
func UpdateAutoHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener ID del auto
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var auto models.Auto
	if err := json.NewDecoder(r.Body).Decode(&auto); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	collection := db.Collection("autos")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": auto}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error al actualizar el auto", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Auto no encontrado", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"mensaje": "Auto actualizado exitosamente",
		"auto":    auto,
	}

	json.NewEncoder(w).Encode(response)
}

// DeleteAutoHandler elimina un auto
func DeleteAutoHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener ID del auto
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	collection := db.Collection("autos")
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Error al eliminar el auto", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Auto no encontrado", http.StatusNotFound)
		return
	}

	response := map[string]string{
		"mensaje": "Auto eliminado exitosamente",
	}

	json.NewEncoder(w).Encode(response)
}

// ToggleFeaturedHandler marca o desmarca un auto como destacado
func ToggleFeaturedHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
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

	// Cambiar el estado featured
	update := bson.M{"$set": bson.M{"featured": !auto.Featured}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error al actualizar el auto", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"mensaje":  "Estado destacado actualizado exitosamente",
		"featured": !auto.Featured,
	}

	json.NewEncoder(w).Encode(response)
}
