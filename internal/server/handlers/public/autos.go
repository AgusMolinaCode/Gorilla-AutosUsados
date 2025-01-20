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

	// Filtrar por marca (case insensitive y coincidencias parciales)
	if marca := r.URL.Query().Get("marca"); marca != "" {
		filter["marca"] = bson.M{"$regex": "^" + marca, "$options": "i"}
	}

	// Filtrar por modelo (case insensitive y coincidencias parciales)
	if modelo := r.URL.Query().Get("modelo"); modelo != "" {
		filter["modelo"] = bson.M{"$regex": "^" + modelo, "$options": "i"}
	}

	// Filtrar por tipo de combustible
	if combustible := r.URL.Query().Get("combustible"); combustible != "" {
		filter["tipo_combustible"] = bson.M{"$regex": "^" + combustible, "$options": "i"}
	}

	// Filtrar por año
	if año := r.URL.Query().Get("año"); año != "" {
		añoInt, err := strconv.Atoi(año)
		if err == nil {
			filter["año"] = añoInt
		}
	}

	// Filtrar por kilometraje específico
	if km := r.URL.Query().Get("kilometraje"); km != "" {
		kmInt, err := strconv.Atoi(km)
		if err == nil {
			filter["kilometraje"] = kmInt
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

	// Filtrar por precio específico
	if precio := r.URL.Query().Get("precio"); precio != "" {
		precioInt, err := strconv.ParseFloat(precio, 64)
		if err == nil {
			filter["precio"] = precioInt
		}
	}

	// Filtrar por rango de precios
	precioFilter := bson.M{}
	if precioMin := r.URL.Query().Get("precio_min"); precioMin != "" {
		if precioMinFloat, err := strconv.ParseFloat(precioMin, 64); err == nil {
			precioFilter["$gte"] = precioMinFloat
		}
	}
	if precioMax := r.URL.Query().Get("precio_max"); precioMax != "" {
		if precioMaxFloat, err := strconv.ParseFloat(precioMax, 64); err == nil {
			precioFilter["$lte"] = precioMaxFloat
		}
	}
	if len(precioFilter) > 0 {
		filter["precio"] = precioFilter
	}

	// Filtrar por autos destacados
	if destacado := r.URL.Query().Get("destacado"); destacado == "true" {
		filter["featured"] = true
	}

	// Filtrar por autos con descuento
	if descuento := r.URL.Query().Get("descuento"); descuento == "true" {
		filter["descuento"] = bson.M{"$gt": 0} // Solo autos con descuento mayor a 0
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

// GetFeaturedAutosHandler obtiene los autos marcados como destacados
func GetFeaturedAutosHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	// Filtrar solo autos destacados
	filter := bson.M{"featured": true}

	collection := db.Collection("autos")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Error al obtener los autos destacados", http.StatusInternalServerError)
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
			"mensaje": "No hay autos destacados disponibles",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(autos)
}
