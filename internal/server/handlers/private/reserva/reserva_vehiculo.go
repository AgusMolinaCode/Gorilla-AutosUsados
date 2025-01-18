package reserva

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// Reserva representa la información de una reserva
type Reserva struct {
	Nombre     string    `json:"nombre"`
	Apellido   string    `json:"apellido"`
	Telefono   string    `json:"telefono"`
	Comentario string    `json:"comentario"`
	FechaHora  time.Time `json:"fecha_hora"`
}

// CrearReservaHandler maneja la creación de una nueva reserva
func CrearReservaHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	stockID := vars["stock_id"]

	var reserva models.Reserva
	if err := json.NewDecoder(r.Body).Decode(&reserva); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Validaciones básicas
	if reserva.Nombre == "" || reserva.Apellido == "" {
		http.Error(w, "Nombre y apellido son requeridos", http.StatusBadRequest)
		return
	}

	if reserva.FechaHora.IsZero() {
		http.Error(w, "Fecha y hora de reserva son requeridas", http.StatusBadRequest)
		return
	}

	collection := db.Collection("autos")
	var auto models.Auto
	filter := bson.M{"stock_id": stockID}
	err := collection.FindOne(context.Background(), filter).Decode(&auto)
	if err != nil {
		http.Error(w, "Auto no encontrado", http.StatusNotFound)
		return
	}

	// Agregar la nueva reserva sin restricciones
	auto.Reservas = append(auto.Reservas, reserva)

	update := bson.M{"$set": bson.M{"reservas": auto.Reservas}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error al actualizar las reservas", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"mensaje":        "Reserva creada exitosamente",
		"reserva":        reserva,
		"total_reservas": len(auto.Reservas),
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func EliminarReservaHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	stockID := vars["stock_id"]
	reservaIndex, _ := strconv.Atoi(vars["reserva_index"])

	collection := db.Collection("autos")
	var auto models.Auto
	filter := bson.M{"stock_id": stockID}
	err := collection.FindOne(context.Background(), filter).Decode(&auto)
	if err != nil {
		http.Error(w, "Auto no encontrado", http.StatusNotFound)
		return
	}

	if reservaIndex < 0 || reservaIndex >= len(auto.Reservas) {
		http.Error(w, "Índice de reserva inválido", http.StatusBadRequest)
		return
	}
	auto.Reservas = append(auto.Reservas[:reservaIndex], auto.Reservas[reservaIndex+1:]...)
	update := bson.M{"$set": bson.M{"reservas": auto.Reservas}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error al actualizar las reservas", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"mensaje": "Reserva eliminada exitosamente",
	}
	json.NewEncoder(w).Encode(response)
}

func EditarReservaHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	stockID := vars["stock_id"]
	reservaIndex, _ := strconv.Atoi(vars["reserva_index"])

	var nuevaReserva models.Reserva
	if err := json.NewDecoder(r.Body).Decode(&nuevaReserva); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	collection := db.Collection("autos")
	var auto models.Auto
	filter := bson.M{"stock_id": stockID}
	err := collection.FindOne(context.Background(), filter).Decode(&auto)
	if err != nil {
		http.Error(w, "Auto no encontrado", http.StatusNotFound)
		return
	}

	if reservaIndex < 0 || reservaIndex >= len(auto.Reservas) {
		http.Error(w, "Índice de reserva inválido", http.StatusBadRequest)
		return
	}
	auto.Reservas[reservaIndex] = nuevaReserva
	update := bson.M{"$set": bson.M{"reservas": auto.Reservas}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error al actualizar las reservas", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"mensaje": "Reserva actualizada exitosamente",
		"reserva": nuevaReserva,
	}
	json.NewEncoder(w).Encode(response)
}

// ObtenerReservasHandler obtiene todas las reservaciones de un auto
func ObtenerReservasHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
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

	// Verificar si hay reservas
	if len(auto.Reservas) == 0 {
		response := map[string]interface{}{
			"mensaje":  "No hay reservas para este auto",
			"reservas": []models.Reserva{},
			"stock_id": stockID,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Preparar respuesta con reservas
	response := map[string]interface{}{
		"mensaje":  "Reservas obtenidas exitosamente",
		"reservas": auto.Reservas,
		"stock_id": stockID,
		"total":    len(auto.Reservas),
	}

	json.NewEncoder(w).Encode(response)
}
