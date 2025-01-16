package models

type Auto struct {
	Marca       string   `json:"marca" bson:"marca" binding:"required"`
	Modelo      string   `json:"modelo" bson:"modelo" binding:"required"`
	TipoVenta   string   `json:"tipo_venta" bson:"tipo_venta" binding:"required"` // "venta" o "alquiler"
	Año         int      `json:"año" bson:"año" binding:"required"`
	Kilometraje int      `json:"kilometraje" bson:"kilometraje" binding:"required"`
	Precio      float64  `json:"precio" bson:"precio" binding:"required"`
	Ciudad      string   `json:"ciudad" bson:"ciudad" binding:"required"`
	Transmision string   `json:"transmision" bson:"transmision" binding:"required"` // "Automático" o "Manual"
	Traccion    string   `json:"traccion" bson:"traccion" binding:"required"`       // "Delantera", "Trasera", "4x4"
	StockID     string   `json:"stock_id" bson:"stock_id" binding:"required"`
	Sucursal    string   `json:"sucursal" bson:"sucursal" binding:"required"`
	Imagenes    []string `json:"imagenes" bson:"imagenes"` // Array de URLs, no requerido
	Version     string   `json:"version" bson:"version" binding:"required"`
	Garantia    string   `json:"garantia" bson:"garantia" binding:"required"` // Ej: "3 meses de garantía mecánica"
}
