package helpers

import (
	"encoding/json"
	"net/http"
)

// JSONResponse escribe una respuesta JSON con el código de estado especificado
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// JSONErrorResponse escribe una respuesta de error JSON
func JSONErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(w, statusCode, map[string]string{"error": message})
}

// JSONSuccessResponse escribe una respuesta exitosa JSON con mensaje
func JSONSuccessResponse(w http.ResponseWriter, statusCode int, message string, data map[string]interface{}) {
	response := map[string]interface{}{"mensaje": message}
	for key, value := range data {
		response[key] = value
	}
	JSONResponse(w, statusCode, response)
}
