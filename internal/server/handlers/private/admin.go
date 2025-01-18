package private

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateAutoHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	var auto models.Auto
	if err := json.NewDecoder(r.Body).Decode(&auto); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if _, err := auto.ValidateRequired(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Establecer created_at y updated_at
	now := time.Now()
	auto.CreatedAt = now
	auto.UpdatedAt = now

	collection := db.Collection("autos")

	// Generar stock_id
	firstLetter := strings.ToUpper(string(auto.Marca[0]))

	// Buscar el último número usado para esta marca
	var lastAuto models.Auto
	opts := options.FindOne().SetSort(bson.M{"stock_id": -1})
	err := collection.FindOne(context.Background(),
		bson.M{"stock_id": bson.M{"$regex": "^" + firstLetter}},
		opts).Decode(&lastAuto)

	var num int
	if err == mongo.ErrNoDocuments {
		num = 1
	} else if err != nil {
		http.Error(w, "Error al generar stock_id", http.StatusInternalServerError)
		return
	} else {
		// Extraer el número del último stock_id
		numStr := lastAuto.StockID[1:]
		num, _ = strconv.Atoi(numStr)
		num++
	}

	// Generar nuevo stock_id
	auto.StockID = fmt.Sprintf("%s%02d", firstLetter, num)

	// Crear el auto con el stock_id generado
	_, err = collection.InsertOne(context.Background(), auto)
	if err != nil {
		http.Error(w, "Error al guardar el auto", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"mensaje": "Auto creado exitosamente",
		"auto":    auto,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateAutoHandler actualiza un auto existente usando stock_id
func UpdateAutoHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener stock_id
	vars := mux.Vars(r)
	stockID := vars["stock_id"]

	// Validar formato de stock_id
	if err := models.ValidateStockID(stockID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener auto existente
	collection := db.Collection("autos")
	var existingAuto models.Auto
	err := collection.FindOne(context.Background(), bson.M{"stock_id": stockID}).Decode(&existingAuto)
	if err != nil {
		http.Error(w, "Auto no encontrado", http.StatusNotFound)
		return
	}

	// Decodificar datos actualizados
	var updateData models.Auto
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if _, err := updateData.ValidateRequired(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mantener el stock_id original
	updateData.StockID = stockID

	// Establecer updated_at
	updateData.UpdatedAt = time.Now()

	// Actualizar el auto
	update := bson.M{"$set": updateData}
	result, err := collection.UpdateOne(context.Background(), bson.M{"stock_id": stockID}, update)
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
		"auto":    updateData,
	}

	json.NewEncoder(w).Encode(response)
}

// DeleteAutoHandler elimina un auto usando stock_id
func DeleteAutoHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
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
	filter := bson.M{"stock_id": stockID}

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
