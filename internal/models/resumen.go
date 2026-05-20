package models

// ItemConsumo representa un producto y su cantidad en el resumen diario
type ItemConsumo struct {
	NombreProducto string
	Cantidad       int
}

// ResumenEstudiante agrupa los consumos del día para un estudiante
type ResumenEstudiante struct {
	IdEstudiante int
	Nombres      string
	Apellidos    string
	NombreGrado  string
	Items        []ItemConsumo
}

// DatosResumenSector contiene todos los datos para la página de resumen
type DatosResumenSector struct {
	Sector     string // "menor" | "mayor"
	Fecha      string // "2026-05-04" (para mostrar y param URL)
	Fechas     []DiaFecha // Semana — reutiliza DiaFecha de common.go
	Resumenes  []ResumenEstudiante
	TotalItems int // Suma de todos los items (badge en navbar)
}
