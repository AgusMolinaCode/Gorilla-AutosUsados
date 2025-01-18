package destacado

import (
	"context"
	"encoding/json"
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

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
