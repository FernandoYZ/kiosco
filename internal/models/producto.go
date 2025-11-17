package models

// Producto representa un producto del kiosco
type Producto struct {
	IdProducto     int
	Nombre         string
	PrecioUnitario float64
	EstaActivo     bool
}
