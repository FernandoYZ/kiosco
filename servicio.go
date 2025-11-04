package main

import (
	"fmt"
	"strings"
	"time"
)

// Servicio contiene la lógica de negocio
type Servicio struct {
	repo *Repositorio
}

// NuevoServicio crea una nueva instancia del servicio
func NuevoServicio(repo *Repositorio) *Servicio {
	return &Servicio{repo: repo}
}

// ObtenerGradosEstaticos retorna la lista estática de grados (sin "Todos")
func ObtenerGradosEstaticos() []InfoGrado {
	return []InfoGrado{
		{IdGrado: 1, Nombre: "5to Primaria"},
		{IdGrado: 2, Nombre: "6to Primaria"},
		{IdGrado: 3, Nombre: "1ro Secundaria"},
		{IdGrado: 4, Nombre: "2do Secundaria"},
		{IdGrado: 5, Nombre: "3ro Secundaria"},
		{IdGrado: 6, Nombre: "4to Secundaria"},
		{IdGrado: 7, Nombre: "5to Secundaria"},
	}
}

// ObtenerSemanaActual retorna el lunes y viernes de la semana actual
func ObtenerSemanaActual() (time.Time, time.Time) {
	ahora := time.Now()

	// Obtener el día de la semana (0 = Domingo, 1 = Lunes, ...)
	diaSemana := int(ahora.Weekday())

	// Ajustar para que Lunes sea 0
	if diaSemana == 0 {
		diaSemana = 7
	}
	diaSemana-- // Ahora Lunes = 0, Martes = 1, ..., Domingo = 6

	// Calcular el lunes de esta semana
	lunes := ahora.AddDate(0, 0, -diaSemana)
	lunes = time.Date(lunes.Year(), lunes.Month(), lunes.Day(), 0, 0, 0, 0, lunes.Location())

	// Calcular el sábado
	sabado := lunes.AddDate(0, 0, 5)
	sabado = time.Date(sabado.Year(), sabado.Month(), sabado.Day(), 23, 59, 59, 0, sabado.Location())

	return lunes, sabado
}

// ObtenerDiasHabiles retorna los días de lunes a sábado
func ObtenerDiasHabiles(inicio time.Time) []time.Time {
	dias := make([]time.Time, 6)
	for i := 0; i < 6; i++ {
		dias[i] = inicio.AddDate(0, 0, i)
	}
	return dias
}

// FormatearSemana retorna el texto de la semana (ej: "27 AL 31 DE OCTUBRE")
func FormatearSemana(inicio, fin time.Time) string {
	meses := []string{
		"ENERO", "FEBRERO", "MARZO", "ABRIL", "MAYO", "JUNIO",
		"JULIO", "AGOSTO", "SEPTIEMBRE", "OCTUBRE", "NOVIEMBRE", "DICIEMBRE",
	}

	mes := meses[inicio.Month()-1]

	// Si ambos están en el mismo mes
	if inicio.Month() == fin.Month() {
		return fmt.Sprintf("%d AL %d DE %s", inicio.Day(), fin.Day(), mes)
	}

	// Si están en meses diferentes
	mesInicio := meses[inicio.Month()-1]
	mesFin := meses[fin.Month()-1]
	return fmt.Sprintf("%d DE %s AL %d DE %s", inicio.Day(), mesInicio, fin.Day(), mesFin)
}

// ObtenerDatosVistaPrincipal prepara todos los datos para la vista principal
func (s *Servicio) ObtenerDatosVistaPrincipal(fechaInicio, fechaFin time.Time, idGrado int, diasDeshabilitados string) (*DatosVistaPrincipal, error) {
	// Obtener estudiantes activos filtrados por grado
	estudiantes, err := s.repo.ObtenerEstudiantesPorGrado(idGrado)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estudiantes: %v", err)
	}

	// Obtener productos activos
	productos, err := s.repo.ObtenerProductosActivos()
	if err != nil {
		return nil, fmt.Errorf("error al obtener productos: %v", err)
	}

	// Obtener consumos de la semana
	consumos, err := s.repo.ObtenerConsumosSemana(fechaInicio, fechaFin)
	if err != nil {
		return nil, fmt.Errorf("error al obtener consumos: %v", err)
	}

	// Crear mapas optimizados con capacidad pre-asignada
	numEstudiantes := len(estudiantes)
	consumosPorDia := make(map[int]map[string]map[int]int, numEstudiantes)
	subTotalesPorEstudiante := make(map[int]float64, numEstudiantes)

	for _, c := range consumos {
		// Mapa para template
		if consumosPorDia[c.IdEstudiante] == nil {
			consumosPorDia[c.IdEstudiante] = make(map[string]map[int]int)
		}
		fechaKey := c.FechaConsumo.Format("2006-01-02")
		if consumosPorDia[c.IdEstudiante][fechaKey] == nil {
			consumosPorDia[c.IdEstudiante][fechaKey] = make(map[int]int)
		}
		consumosPorDia[c.IdEstudiante][fechaKey][c.IdProducto] = c.Cantidad

		// Subtotales
		subTotalesPorEstudiante[c.IdEstudiante] += c.TotalLinea
	}

	// Obtener deudas y pagos en batch (2 queries en lugar de N*2)
	deudasAnteriores, err := s.repo.ObtenerDeudasAnterioresBatch(idGrado, fechaInicio)
	if err != nil {
		return nil, fmt.Errorf("error al obtener deudas anteriores: %v", err)
	}

	pagosSemana, err := s.repo.ObtenerPagosSemanaBatch(idGrado, fechaInicio, fechaFin)
	if err != nil {
		return nil, fmt.Errorf("error al obtener pagos: %v", err)
	}

	// Calcular datos de cada estudiante (sin queries adicionales)
	estudiantesConData := make([]EstudianteConDeuda, 0, len(estudiantes))
	for _, est := range estudiantes {
		subTotal := subTotalesPorEstudiante[est.IdEstudiante]
		deudaAnterior := deudasAnteriores[est.IdEstudiante]
		descuento := pagosSemana[est.IdEstudiante]
		total := subTotal + deudaAnterior - descuento

		estudiantesConData = append(estudiantesConData, EstudianteConDeuda{
			Estudiante:    est,
			SubTotal:      subTotal,
			DeudaAnterior: deudaAnterior,
			Descuento:     descuento,
			Total:         total,
		})
	}

	// Parsear días deshabilitados del parámetro URL (formato: "2025-01-05,2025-01-12")
	var mapaDiasDeshabilitados map[string]bool
	if diasDeshabilitados != "" {
		fechas := strings.Split(diasDeshabilitados, ",")
		mapaDiasDeshabilitados = make(map[string]bool, len(fechas))
		for _, f := range fechas {
			f = strings.TrimSpace(f)
			if f != "" {
				mapaDiasDeshabilitados[f] = true
			}
		}
	}

	// Obtener todos los días y marcar su estado
	todosDias := ObtenerDiasHabiles(fechaInicio)
	diasConEstado := make([]DiaConEstado, 0, len(todosDias))
	diasHabilesFiltrados := make([]time.Time, 0, len(todosDias))

	for _, dia := range todosDias {
		fechaKey := dia.Format("2006-01-02")
		esHabil := !mapaDiasDeshabilitados[fechaKey]

		diasConEstado = append(diasConEstado, DiaConEstado{
			Fecha:   dia,
			EsHabil: esHabil,
		})

		if esHabil {
			diasHabilesFiltrados = append(diasHabilesFiltrados, dia)
		}
	}

	return &DatosVistaPrincipal{
		Semana:             FormatearSemana(fechaInicio, fechaFin),
		FechaInicio:        fechaInicio,
		FechaFin:           fechaFin,
		DiasConEstado:      diasConEstado,
		DiasHabiles:        diasHabilesFiltrados,
		Productos:          productos,
		EstudiantesConData: estudiantesConData,
		ConsumosPorDia:     consumosPorDia,
		Grados:             ObtenerGradosEstaticos(),
		GradoSeleccionado:  idGrado,
		DiasDeshabilitados: diasDeshabilitados,
	}, nil
}

// RegistrarConsumoDesdeFormulario procesa el registro de un consumo desde el formulario
func (s *Servicio) RegistrarConsumoDesdeFormulario(idEstudiante, idProducto, cantidad int, fecha time.Time) error {
	// Obtener el precio actual del producto
	producto, err := s.repo.ObtenerProductoPorId(idProducto)
	if err != nil {
		return fmt.Errorf("producto no encontrado: %v", err)
	}

	// Actualizar o insertar el consumo
	return s.repo.ActualizarConsumo(idEstudiante, idProducto, fecha, cantidad, producto.PrecioUnitario)
}

// RegistrarPagoDesdeFormulario procesa el registro de un pago
func (s *Servicio) RegistrarPagoDesdeFormulario(idEstudiante int, monto float64, fecha time.Time) error {
	if monto <= 0 {
		return fmt.Errorf("el monto debe ser mayor a cero")
	}

	pago := Pago{
		IdEstudiante: idEstudiante,
		Monto:        monto,
		FechaPago:    fecha,
	}

	return s.repo.RegistrarPago(pago)
}
