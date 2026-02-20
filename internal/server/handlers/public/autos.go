package public

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

// parseInt safely parses an integer from string, logging errors
func parseInt(s string, fieldName string) (int, bool) {
	val, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("Invalid %s value '%s': %v", fieldName, s, err)
		return 0, false
	}
	return val, true
}

// parseFloat safely parses a float from string, logging errors
// Returns the parsed value and true if successful, 0 and false otherwise
func parseFloat(s string, fieldName string) (float64, bool) {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("Invalid %s value '%s': %v", fieldName, s, err)
		return 0, false
	}
	return val, true
}

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
		if añoInt, ok := parseInt(año, "año"); ok {
			filter["año"] = añoInt
		}
	}

	// Filtrar por kilometraje específico
	if km := r.URL.Query().Get("kilometraje"); km != "" {
		if kmInt, ok := parseInt(km, "kilometraje"); ok {
			filter["kilometraje"] = kmInt
		}
	}

	// Filtrar por rango de kilometraje
	kmFilter := bson.M{}
	if kmMin := r.URL.Query().Get("km_min"); kmMin != "" {
		if kmMinInt, ok := parseInt(kmMin, "km_min"); ok {
			kmFilter["$gte"] = kmMinInt
		}
	}
	if kmMax := r.URL.Query().Get("km_max"); kmMax != "" {
		if kmMaxInt, ok := parseInt(kmMax, "km_max"); ok {
			kmFilter["$lte"] = kmMaxInt
		}
	}
	if len(kmFilter) > 0 {
		filter["kilometraje"] = kmFilter
	}

	// Filtrar por precio específico
	if precio := r.URL.Query().Get("precio"); precio != "" {
		if precioFloat, ok := parseFloat(precio, "precio"); ok {
			filter["precio"] = precioFloat
		}
	}

	// Filtrar por rango de precios
	precioFilter := bson.M{}
	if precioMin := r.URL.Query().Get("precio_min"); precioMin != "" {
		if precioMinFloat, ok := parseFloat(precioMin, "precio_min"); ok {
			precioFilter["$gte"] = precioMinFloat
		}
	}
	if precioMax := r.URL.Query().Get("precio_max"); precioMax != "" {
		if precioMaxFloat, ok := parseFloat(precioMax, "precio_max"); ok {
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

	// Configurar las opciones de ordenamiento
	opts := options.Find()

	// Ordenar por precio
	if sort := r.URL.Query().Get("sort_precio"); sort != "" {
		switch sort {
		case "asc":
			opts.SetSort(bson.D{{Key: "precio", Value: 1}}) // Orden ascendente
		case "desc":
			opts.SetSort(bson.D{{Key: "precio", Value: -1}}) // Orden descendente
		}
	}

	// Ordenar por fecha de publicación
	if sort := r.URL.Query().Get("sort_fecha"); sort != "" {
		switch sort {
		case "nuevo":
			opts.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Más nuevo primero
		case "viejo":
			opts.SetSort(bson.D{{Key: "created_at", Value: 1}}) // Más viejo primero
		}
	}

	// Ordenar por kilometraje
	if sort := r.URL.Query().Get("sort_km"); sort != "" {
		switch sort {
		case "mayor":
			opts.SetSort(bson.D{{Key: "kilometraje", Value: -1}}) // Más kilómetros primero
		case "menor":
			opts.SetSort(bson.D{{Key: "kilometraje", Value: 1}}) // Menos kilómetros primero
		}
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
