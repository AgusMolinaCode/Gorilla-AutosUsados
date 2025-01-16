package public

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAutosHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	filter := bson.M{}

	// Filtrar por marca (case insensitive)
	if marca := r.URL.Query().Get("marca"); marca != "" {
		filter["marca"] = bson.M{"$regex": marca, "$options": "i"}
	}

	// Filtrar por modelo
	if modelo := r.URL.Query().Get("modelo"); modelo != "" {
		filter["modelo"] = modelo
	}

	// Filtrar por año
	if año := r.URL.Query().Get("año"); año != "" {
		añoInt, err := strconv.Atoi(año)
		if err == nil {
			filter["año"] = añoInt
		}
	}

	// Filtrar por rango de kilometraje
	kmFilter := bson.M{}

	if kmMin := r.URL.Query().Get("km_min"); kmMin != "" {
		if kmMinInt, err := strconv.Atoi(kmMin); err == nil {
			kmFilter["$gte"] = kmMinInt
		}
	}

	if kmMax := r.URL.Query().Get("km_max"); kmMax != "" {
		if kmMaxInt, err := strconv.Atoi(kmMax); err == nil {
			kmFilter["$lte"] = kmMaxInt
		}
	}

	if len(kmFilter) > 0 {
		filter["kilometraje"] = kmFilter
	}

	collection := db.Collection("autos")
	cursor, err := collection.Find(context.Background(), filter)
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

	if len(autos) == 0 {
		response := map[string]string{
			"mensaje": "No se encontraron autos con los filtros especificados",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(autos)
}
