package models

import (
	"fmt"
	"regexp"
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
	StockID     string   `json:"stock_id" bson:"stock_id" binding:"required"`
	Sucursal    string   `json:"sucursal" bson:"sucursal" binding:"required"`
	Imagenes    []string `json:"imagenes" bson:"imagenes"`
	Version     string   `json:"version" bson:"version" binding:"required"`
	Garantia    string   `json:"garantia" bson:"garantia" binding:"required"`
	Featured    bool     `json:"featured" bson:"featured"`
}

func ValidateStockID(stockID string) error {
	pattern := `^[A-Z][0-9]{2}$`
	matched, _ := regexp.MatchString(pattern, stockID)
	if !matched {
		return fmt.Errorf("stock_id inválido: debe ser una letra mayúscula seguida de dos números")
	}
	return nil
}
