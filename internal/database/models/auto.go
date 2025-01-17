package models

import (
	"fmt"
	"regexp"
)

// Definir constantes para los estados
const (
	EstadoDisponible = "disponible"
	EstadoReservado  = "reservado"
	EstadoVendido    = "vendido"
)

type Auto struct {
	Marca       string   `json:"marca" bson:"marca" binding:"required"`
	Modelo      string   `json:"modelo" bson:"modelo" binding:"required"`
	TipoVenta   string   `json:"tipo_venta" bson:"tipo_venta" binding:"required"`
	Año         int      `json:"año" bson:"año" binding:"required"`
	Kilometraje int      `json:"kilometraje" bson:"kilometraje" binding:"required"`
	Precio      float64  `json:"precio" bson:"precio" binding:"required"`
	Ciudad      string   `json:"ciudad" bson:"ciudad" binding:"required"`
	Transmision string   `json:"transmision" bson:"transmision" binding:"required"`
	Traccion    string   `json:"traccion" bson:"traccion" binding:"required"`
	StockID     string   `json:"stock_id" bson:"stock_id"`
	Sucursal    string   `json:"sucursal" bson:"sucursal" binding:"required"`
	Imagenes    []string `json:"imagenes" bson:"imagenes"`
	Version     string   `json:"version" bson:"version" binding:"required"`
	Garantia    string   `json:"garantia" bson:"garantia" binding:"required"`
	Featured    bool     `json:"featured" bson:"featured"`
	Estado      string   `json:"estado" bson:"estado"`
}

func (a *Auto) ValidateRequired() ([]string, error) {
	missingFields := []string{}

	// Validar campos requeridos
	if a.Marca == "" {
		missingFields = append(missingFields, "marca")
	}
	if a.Modelo == "" {
		missingFields = append(missingFields, "modelo")
	}
	if a.TipoVenta == "" {
		missingFields = append(missingFields, "tipo_venta")
	}
	if a.Año == 0 {
		missingFields = append(missingFields, "año")
	}
	if a.Kilometraje < 0 {
		missingFields = append(missingFields, "kilometraje")
	}
	if a.Precio <= 0 {
		missingFields = append(missingFields, "precio")
	}
	if a.Ciudad == "" {
		missingFields = append(missingFields, "ciudad")
	}
	if a.Transmision == "" {
		missingFields = append(missingFields, "transmision")
	}
	if a.Traccion == "" {
		missingFields = append(missingFields, "traccion")
	}
	if a.Sucursal == "" {
		missingFields = append(missingFields, "sucursal")
	}
	if a.Version == "" {
		missingFields = append(missingFields, "version")
	}
	if a.Garantia == "" {
		missingFields = append(missingFields, "garantia")
	}

	// Validar estado
	if a.Estado == "" {
		a.Estado = EstadoDisponible
	} else if a.Estado != EstadoDisponible &&
		a.Estado != EstadoReservado &&
		a.Estado != EstadoVendido {
		missingFields = append(missingFields, "estado inválido")
	}

	if len(missingFields) > 0 {
		return missingFields, fmt.Errorf("faltan campos requeridos: %v", missingFields)
	}

	return nil, nil
}

func ValidateStockID(stockID string) error {
	pattern := `^[A-Z][0-9]{2}$`
	matched, _ := regexp.MatchString(pattern, stockID)
	if !matched {
		return fmt.Errorf("stock_id inválido: debe ser una letra mayúscula seguida de dos números")
	}
	return nil
}
