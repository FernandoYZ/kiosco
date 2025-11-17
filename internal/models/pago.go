package models

import "time"

// Pago representa un pago/descuento
type Pago struct {
	IdPago       int
	IdEstudiante int
	Monto        float64
	FechaPago    time.Time
}
