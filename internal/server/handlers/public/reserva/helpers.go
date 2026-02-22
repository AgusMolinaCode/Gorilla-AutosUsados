package reserva

import (
	"context"
	"encoding/json"
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"

	"go.mongodb.org/mongo-driver/bson"
)

// AutoResult contiene el resultado de buscar un auto por stock_id
type AutoResult struct {
	Auto  models.Auto
	Found bool
	Error error
}

// FindAutoByStockID busca un auto por su stock_id y devuelve el resultado encapsulado
func FindAutoByStockID(ctx context.Context, db database.Service, stockID string) AutoResult {
	collection := db.Collection("autos")
	var auto models.Auto
	filter := bson.M{"stock_id": stockID}
	err := collection.FindOne(ctx, filter).Decode(&auto)
	if err != nil {
		return AutoResult{Found: false, Error: err}
	}
	return AutoResult{Auto: auto, Found: true}
}

// WriteNotFoundResponse escribe una respuesta 404 estándar
func WriteNotFoundResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// WriteErrorResponse escribe una respuesta de error con código de estado
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// UpdateAutoReservas actualiza las reservas de un auto en la base de datos
func UpdateAutoReservas(ctx context.Context, db database.Service, stockID string, reservas []models.Reserva) error {
	collection := db.Collection("autos")
	filter := bson.M{"stock_id": stockID}
	update := bson.M{"$set": bson.M{"reservas": reservas}}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
