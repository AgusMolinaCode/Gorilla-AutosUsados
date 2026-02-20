package reserva

import (
	"context"
	"encoding/json"
	"net/http"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

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

	// Generar un ID aleatorio para la reserva
	reserva.ID = models.GenerarIDReserva()

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

	// Verificar si el cliente ya tiene una reserva para este auto
	for _, r := range auto.Reservas {
		if r.Nombre == reserva.Nombre && r.Apellido == reserva.Apellido {
			http.Error(w, "Ya existe una reserva activa para este cliente y vehículo", http.StatusConflict)
			return
		}
	}

	// Agregar la nueva reserva
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
