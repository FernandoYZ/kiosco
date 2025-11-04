package main

import "time"

// Grado representa un nivel educativo
type Grado struct {
	IdGrado     int
	NivelGrado  string // Primaria, Secundaria
	NombreGrado string // Primero, Segundo, etc.
	AnioGrado   string // 1ro, 2do, etc.
}

// Estudiante representa un estudiante
type Estudiante struct {
	IdEstudiante int
	Nombres      string
	Apellidos    string
	IdGrado      int
	EstaActivo   bool
	NombreGrado  string // Para mostrar en la vista
}

// Producto representa un producto del kiosco
type Producto struct {
	IdProducto      int
	Nombre          string
	PrecioUnitario  float64
	EstaActivo      bool
}

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

// Pago representa un pago/descuento
type Pago struct {
	IdPago       int
	IdEstudiante int
	Monto        float64
	FechaPago    time.Time
}

// EstudianteConDeuda contiene los datos del estudiante y sus cálculos
type EstudianteConDeuda struct {
	Estudiante
	SubTotal      float64 // Total de consumos de la semana
	DeudaAnterior float64 // Deuda de semanas anteriores
	Descuento     float64 // Pagos realizados
	Total         float64 // SubTotal + DeudaAnterior - Descuento
}

// ConsumoDelDia agrupa los consumos por día y producto
type ConsumoDelDia struct {
	IdProducto int
	Cantidad   int
}

// DiaConEstado representa un día con su estado habilitado/deshabilitado
type DiaConEstado struct {
	Fecha    time.Time
	EsHabil  bool
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

// InfoGrado contiene información básica de un grado
type InfoGrado struct {
	IdGrado int
	Nombre  string
}

// DatosEditarConsumos contiene los datos para editar consumos de un día
type DatosEditarConsumos struct {
	IdEstudiante      int
	NombreEstudiante  string
	Fecha             time.Time
	Productos         []Producto
	Consumos          map[int]map[string]map[int]int
	GradoSeleccionado int
}

// DatosEditarPagos contiene los datos para editar pagos de una semana
type DatosEditarPagos struct {
	IdEstudiante      int
	NombreEstudiante  string
	FechaInicio       time.Time
	FechaFin          time.Time
	Pagos             []Pago
	TotalPagos        float64
	GradoSeleccionado int
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
