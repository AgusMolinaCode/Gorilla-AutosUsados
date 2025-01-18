package public

import (
	"encoding/json"
	"net/http"
)

func MensajeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"mensaje": "Bienvenido a la Agencia de Autos",
		"info":    "Para acceder al panel de administración",
	}

	json.NewEncoder(w).Encode(response)
}
