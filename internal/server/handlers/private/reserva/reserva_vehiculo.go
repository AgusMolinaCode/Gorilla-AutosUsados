package reserva

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"go-gorilla-autos/internal/database"
	"go-gorilla-autos/internal/database/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// Reserva representa la información de una reserva
type Reserva struct {
	ID         string    `json:"id" bson:"id"`
	Nombre     string    `json:"nombre"`
	Apellido   string    `json:"apellido"`
	Telefono   string    `json:"telefono"`
	Comentario string    `json:"comentario"`
	FechaHora  time.Time `json:"fecha_hora"`
}

// GenerarIDReserva genera un ID aleatorio para la reserva
func GenerarIDReserva() string {
	rand.Seed(time.Now().UnixNano())
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"

	id := string(letters[rand.Intn(len(letters))]) // Letra aleatoria
	for i := 0; i < 3; i++ {
		id += string(numbers[rand.Intn(len(numbers))]) // Tres números aleatorios
	}
	return id
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

	// Generar un ID aleatorio para la reserva
	reserva.ID = GenerarIDReserva()

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

func EliminarReservaHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	stockID := vars["stock_id"]
	reservaID := vars["reserva_id"]

	collection := db.Collection("autos")
	var auto models.Auto
	filter := bson.M{"stock_id": stockID}
	err := collection.FindOne(context.Background(), filter).Decode(&auto)
	if err != nil {
		http.Error(w, "Auto no encontrado", http.StatusNotFound)
		return
	}

	// Buscar la reserva por ID
	var indexToDelete int
	found := false
	for i := range auto.Reservas {
		if auto.Reservas[i].ID == reservaID {
			indexToDelete = i
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		return
	}

	// Eliminar la reserva
	auto.Reservas = append(auto.Reservas[:indexToDelete], auto.Reservas[indexToDelete+1:]...)

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
	reservaID := vars["reserva_id"]

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

	// Buscar la reserva por ID
	var reservaEncontrada *models.Reserva
	for i := range auto.Reservas {
		if auto.Reservas[i].ID == reservaID {
			reservaEncontrada = &auto.Reservas[i]
			break
		}
	}

	if reservaEncontrada == nil {
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		return
	}

	// Actualizar la reserva
	reservaEncontrada.Nombre = nuevaReserva.Nombre
	reservaEncontrada.Apellido = nuevaReserva.Apellido
	reservaEncontrada.Telefono = nuevaReserva.Telefono
	reservaEncontrada.Comentario = nuevaReserva.Comentario
	reservaEncontrada.FechaHora = nuevaReserva.FechaHora

	update := bson.M{"$set": bson.M{"reservas": auto.Reservas}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error al actualizar las reservas", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"mensaje": "Reserva actualizada exitosamente",
		"reserva": reservaEncontrada,
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
