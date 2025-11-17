package models

// Estudiante representa un estudiante
type Estudiante struct {
	IdEstudiante int
	Nombres      string
	Apellidos    string
	IdGrado      int
	EstaActivo   bool
	NombreGrado  string // Para mostrar en la vista
}

// EstudianteConDeuda contiene los datos del estudiante y sus cálculos
type EstudianteConDeuda struct {
	Estudiante
	SubTotal      float64 // Total de consumos de la semana
	DeudaAnterior float64 // Deuda de semanas anteriores
	Descuento     float64 // Pagos realizados
	Total         float64 // SubTotal + DeudaAnterior - Descuento
}