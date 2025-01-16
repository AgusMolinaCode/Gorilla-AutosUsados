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
