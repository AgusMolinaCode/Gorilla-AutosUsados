package public

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// writeJSONResponse escribe una respuesta JSON con el código de estado especificado
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func GetAutosHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	// Usar el parser de query params
	qp := NewQueryParser(r)

	filter := bson.M{}

	// Filtrar por marca (case insensitive y coincidencias parciales)
	if marca := qp.GetString("marca"); marca != "" {
		filter["marca"] = bson.M{"$regex": "^" + marca, "$options": "i"}
	}

	// Filtrar por modelo (case insensitive y coincidencias parciales)
	if modelo := qp.GetString("modelo"); modelo != "" {
		filter["modelo"] = bson.M{"$regex": "^" + modelo, "$options": "i"}
	}

	// Filtrar por tipo de combustible
	if combustible := qp.GetString("combustible"); combustible != "" {
		filter["tipo_combustible"] = bson.M{"$regex": "^" + combustible, "$options": "i"}
	}

	// Filtrar por año
	if año := qp.GetInt("año"); año > 0 {
		filter["año"] = año
	}

	// Filtrar por kilometraje específico
	if km := qp.GetInt("kilometraje"); km > 0 {
		filter["kilometraje"] = km
	}

	// Filtrar por rango de kilometraje
	kmFilter := bson.M{}
	if kmMin := qp.GetInt("km_min"); kmMin > 0 {
		kmFilter["$gte"] = kmMin
	}
	if kmMax := qp.GetInt("km_max"); kmMax > 0 {
		kmFilter["$lte"] = kmMax
	}
	if len(kmFilter) > 0 {
		filter["kilometraje"] = kmFilter
	}

	// Filtrar por precio específico
	if precio := qp.GetFloat("precio"); precio > 0 {
		filter["precio"] = precio
	}

	// Filtrar por rango de precios
	precioFilter := bson.M{}
	if precioMin := qp.GetFloat("precio_min"); precioMin > 0 {
		precioFilter["$gte"] = precioMin
	}
	if precioMax := qp.GetFloat("precio_max"); precioMax > 0 {
		precioFilter["$lte"] = precioMax
	}
	if len(precioFilter) > 0 {
		filter["precio"] = precioFilter
	}

	// Filtrar por autos destacados
	if qp.GetString("destacado") == "true" {
		filter["featured"] = true
	}

	// Filtrar por autos con descuento
	if qp.GetString("descuento") == "true" {
		filter["descuento"] = bson.M{"$gt": 0}
	}

	// Configurar las opciones de ordenamiento
	opts := options.Find()

	// Ordenar por precio
	if sortOrder := qp.GetSortOrder("sort_precio", "asc", "desc"); sortOrder != 0 {
		opts.SetSort(bson.D{{Key: "precio", Value: sortOrder}})
	}

	// Ordenar por fecha de publicación
	if sortOrder := qp.GetSortOrder("sort_fecha", "viejo", "nuevo"); sortOrder != 0 {
		opts.SetSort(bson.D{{Key: "created_at", Value: -sortOrder}})
	}

	// Ordenar por kilometraje
	if sortOrder := qp.GetSortOrder("sort_km", "menor", "mayor"); sortOrder != 0 {
		opts.SetSort(bson.D{{Key: "kilometraje", Value: sortOrder}})
	}

	collection := db.Collection("autos")
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Printf("Error fetching autos from database: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Error al obtener los autos",
		})
		return
	}
	defer cursor.Close(context.Background())

	var autos []models.Auto
	if err = cursor.All(context.Background(), &autos); err != nil {
		log.Printf("Error decoding autos: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Error al decodificar los autos",
		})
		return
	}

	if len(autos) == 0 {
		writeJSONResponse(w, http.StatusNotFound, map[string]string{
			"mensaje": "No se encontraron autos con los filtros especificados",
		})
		return
	}

	writeJSONResponse(w, http.StatusOK, autos)
}

// GetFeaturedAutosHandler obtiene los autos marcados como destacados
func GetFeaturedAutosHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")

	// Filtrar solo autos destacados
	filter := bson.M{"featured": true}

	collection := db.Collection("autos")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Error fetching featured autos: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Error al obtener los autos destacados",
		})
		return
	}
	defer cursor.Close(context.Background())

	var autos []models.Auto
	if err = cursor.All(context.Background(), &autos); err != nil {
		log.Printf("Error decoding featured autos: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Error al decodificar los autos",
		})
		return
	}

	if len(autos) == 0 {
		writeJSONResponse(w, http.StatusNotFound, map[string]string{
			"mensaje": "No hay autos destacados disponibles",
		})
		return
	}

	writeJSONResponse(w, http.StatusOK, autos)
}
