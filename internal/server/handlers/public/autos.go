package public

import (
	"context"
	"encoding/json"
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAutosHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	collection := db.Collection("autos")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Error al obtener los autos", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var autos []models.Auto
	if err = cursor.All(context.Background(), &autos); err != nil {
		http.Error(w, "Error al decodificar los autos", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(autos)
}
