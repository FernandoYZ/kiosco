package models

import "time"

// DiaConEstado representa un día con su estado habilitado/deshabilitado
type DiaConEstado struct {
	Fecha   time.Time
	EsHabil bool
}

// DatosVistaPrincipal contiene todos los datos para la vista principal
type DatosVistaPrincipal struct {
	Semana             string
	FechaInicio        time.Time
	FechaFin           time.Time
	DiasConEstado      []DiaConEstado // Todos los días con su estado
	DiasHabiles        []time.Time    // Solo los días habilitados
	Productos          []Producto
	EstudiantesConData []EstudianteConDeuda
	ConsumosPorDia     map[int]map[string]map[int]int // [id_estudiante][fecha][id_producto]cantidad
	Grados             []InfoGrado
	GradoSeleccionado  int
	DiasDeshabilitados string // Parámetro URL con fechas separadas por comas
}

// DatosEditarConsumos contiene los datos para editar consumos de un día
type DatosEditarConsumos struct {
	IdEstudiante      int
	NombreEstudiante  string
	Fecha             time.Time
	Productos         []Producto
	Consumos          map[int]map[string]map[int]int
	GradoSeleccionado int
	Sector            string // "menor" | "mayor" | "" (empty = entrada desde grilla semanal)
	PuedeEditar       bool   // Permisos del usuario autenticado
}

// DatosEditarPagos contiene los datos para editar pagos de una semana
type DatosEditarPagos struct {
	IdEstudiante      int
	NombreEstudiante  string
	FechaInicio       time.Time
	FechaFin          time.Time
	Pagos             []Pago
	TotalPagos        float64
	DeudaActual       float64
	GradoSeleccionado int
}

// DatosConsumoSemanal contiene los datos para ver el consumo semanal de un estudiante
type DatosConsumoSemanal struct {
	IdEstudiante      int
	NombreEstudiante  string
	FechaInicio       time.Time
	FechaFin          time.Time
	ConsumosPorDia    []ConsumoDiario
	SubTotal          float64
	DeudaAnterior     float64
	Pagos             float64
	Total             float64
	GradoSeleccionado int
}

// DiaFecha representa un día de la semana con su fecha formateada
type DiaFecha struct {
	Nombre       string // "Lunes"
	Fecha        string // "2025-04-14" (para el value del select)
	FechaFormato string // "14 Abr"     (para mostrar al usuario)
	EsHoy        bool
}
