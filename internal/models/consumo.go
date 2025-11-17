package models

import "time"

// Consumo representa una venta/consumo
type Consumo struct {
	IdConsumo           int64
	IdEstudiante        int
	IdProducto          int
	Cantidad            int
	PrecioUnitarioVenta float64
	TotalLinea          float64
	FechaConsumo        time.Time
}

// ConsumoDelDia agrupa los consumos por día y producto
type ConsumoDelDia struct {
	IdProducto int
	Cantidad   int
}

// ConsumoDiario representa el consumo de un día específico
type ConsumoDiario struct {
	Fecha     time.Time
	Productos []ConsumoProducto
	Total     float64
}

// ConsumoProducto representa un producto consumido con su cantidad y precio
type ConsumoProducto struct {
	Nombre   string
	Cantidad int
	Precio   float64
	Total    float64
}
