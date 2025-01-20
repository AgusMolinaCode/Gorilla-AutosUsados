package models

import (
	"fmt"
	"regexp"
	"time"
)

// Definir constantes para los estados
const (
	EstadoDisponible = "disponible"
	EstadoReservado  = "reservado"
	EstadoVendido    = "vendido"
)

type Auto struct {
	Marca                          string             `json:"marca" bson:"marca" binding:"required"`
	Modelo                         string             `json:"modelo" bson:"modelo" binding:"required"`
	TipoVenta                      string             `json:"tipo_venta" bson:"tipo_venta" binding:"required"`
	Año                            int                `json:"año" bson:"año" binding:"required"`
	Kilometraje                    int                `json:"kilometraje" bson:"kilometraje" binding:"required"`
	Precio                         float64            `json:"precio" bson:"precio" binding:"required"`
	Ciudad                         string             `json:"ciudad" bson:"ciudad" binding:"required"`
	Transmision                    string             `json:"transmision" bson:"transmision" binding:"required"`
	Traccion                       string             `json:"traccion" bson:"traccion" binding:"required"`
	StockID                        string             `json:"stock_id" bson:"stock_id"`
	Sucursal                       string             `json:"sucursal" bson:"sucursal" binding:"required"`
	Imagenes                       []string           `json:"imagenes" bson:"imagenes"`
	Version                        string             `json:"version" bson:"version" binding:"required"`
	Garantia                       string             `json:"garantia" bson:"garantia" binding:"required"`
	Featured                       bool               `json:"featured" bson:"featured"`
	Estado                         string             `json:"estado" bson:"estado"`
	Descuento                      float64            `json:"descuento" bson:"descuento"`
	CreatedAt                      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt                      time.Time          `json:"updated_at" bson:"updated_at"`
	Imagen_Portada                 string             `json:"imagen_portada" bson:"imagen_portada"`
	Imagenes_Imperfecciones        []string           `json:"imagenes_imperfecciones" bson:"imagenes_imperfecciones"`
	EquipamientoDestacado          []string           `json:"equipamiento_destacado" bson:"equipamiento_destacado"`
	CaracteristicasGeneral         map[string]string  `json:"caracteristicas_general" bson:"caracteristicas_general"`
	CaracteristicasExterior        map[string]string  `json:"caracteristicas_exterior" bson:"caracteristicas_exterior"`
	CaracteristicasSeguridad       map[string]string  `json:"caracteristicas_seguridad" bson:"caracteristicas_seguridad"`
	CaracteristicasConfort         map[string]string  `json:"caracteristicas_confort" bson:"caracteristicas_confort"`
	CaracteristicasInterior        map[string]string  `json:"caracteristicas_interior" bson:"caracteristicas_interior"`
	CaracteristicasEntretenimiento map[string]string  `json:"caracteristicas_entretenimiento" bson:"caracteristicas_entretenimiento"`
	ReservadoPor                   *ReservadoInfo     `json:"reservado_por" bson:"reservado_por,omitempty"`
	VendidoPor                     *VendidoInfo       `json:"vendido_por" bson:"vendido_por,omitempty"`
	EnNegociacion                  *NegociacionInfo   `json:"en_negociacion" bson:"en_negociacion,omitempty"`
	EnMantenimiento                *MantenimientoInfo `json:"en_mantenimiento" bson:"en_mantenimiento,omitempty"`
	Reservas                       []Reserva          `json:"reservas" bson:"reservas"`
	TipoCombustible                string             `json:"tipo_combustible" bson:"tipo_combustible" binding:"required"`
	Moneda                         string             `json:"moneda" bson:"moneda" binding:"required"`
}

type ReservadoInfo struct {
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	Celular    string `json:"celular"`
	Comentario string `json:"comentario"`
}

type VendidoInfo struct {
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	Celular    string `json:"celular"`
	Comentario string `json:"comentario"`
}

type NegociacionInfo struct {
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	Celular    string `json:"celular"`
	Comentario string `json:"comentario"`
}

type MantenimientoInfo struct {
	Taller     string `json:"taller"`
	Mecanico   string `json:"mecanico"`
	Celular    string `json:"celular"`
	Comentario string `json:"comentario"`
}

type Reserva struct {
	ID         string    `json:"id" bson:"id"`
	Nombre     string    `json:"nombre"`
	Apellido   string    `json:"apellido"`
	Telefono   string    `json:"telefono"`
	Comentario string    `json:"comentario"`
	FechaHora  time.Time `json:"fecha_hora"`
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
	if a.TipoCombustible == "" {
		missingFields = append(missingFields, "tipo_combustible")
	}
	if a.Moneda == "" {
		missingFields = append(missingFields, "moneda")
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

	// Validar descuento
	if a.Descuento < 0 {
		missingFields = append(missingFields, "descuento no puede ser negativo")
	}

	// Validar que el descuento no supere el precio original
	if a.Descuento > a.Precio {
		missingFields = append(missingFields, "descuento no puede ser mayor al precio original")
	}

	// Validar que el precio con descuento sea positivo
	if a.Precio-a.Descuento <= 0 {
		missingFields = append(missingFields, "descuento no puede reducir el precio a cero o negativo")
	}

	// Validar que haya al menos un equipamiento destacado
	if len(a.EquipamientoDestacado) == 0 {
		missingFields = append(missingFields, "equipamiento_destacado")
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
