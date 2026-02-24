package public

import (
	"net/http"
	"strconv"
)

// QueryParser encapsula la lógica de parseo de query parameters
type QueryParser struct {
	Query map[string]string
}

// NewQueryParser crea un nuevo parser a partir de la request
func NewQueryParser(r *http.Request) *QueryParser {
	query := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			query[key] = values[0]
		}
	}
	return &QueryParser{Query: query}
}

// GetString obtiene un string del query param, o empty string si no existe
func (qp *QueryParser) GetString(key string) string {
	return qp.Query[key]
}

// GetInt obtiene un int del query param, o 0 si no existe o es inválido
func (qp *QueryParser) GetInt(key string) int {
	val, err := strconv.Atoi(qp.Query[key])
	if err != nil {
		return 0
	}
	return val
}

// GetFloat obtiene un float64 del query param, o 0 si no existe o es inválido
func (qp *QueryParser) GetFloat(key string) float64 {
	val, err := strconv.ParseFloat(qp.Query[key], 64)
	if err != nil {
		return 0
	}
	return val
}

// Has verifica si existe el query param
func (qp *QueryParser) Has(key string) bool {
	_, exists := qp.Query[key]
	return exists && qp.Query[key] != ""
}

// GetSortOrder obtiene el orden de sort (1 para asc, -1 para desc, 0 si no es válido)
func (qp *QueryParser) GetSortOrder(key string, ascValue string, descValue string) int {
	val := qp.GetString(key)
	switch val {
	case ascValue:
		return 1
	case descValue:
		return -1
	default:
		return 0
	}
}
