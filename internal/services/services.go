package services

import (
	"fmt"
	"kiosco/internal/models"
	"kiosco/internal/repositories"
	"kiosco/internal/utils"
	"strings"
	"time"
)

// Servicio contiene la lógica de negocio
type Servicio struct {
	Repo *repositories.Repositorio
}

// NuevoServicio crea una nueva instancia del servicio
func NuevoServicio(repo *repositories.Repositorio) *Servicio {
	return &Servicio{Repo: repo}
}

// ObtenerDatosVistaPrincipal prepara todos los datos para la vista principal
func (s *Servicio) ObtenerDatosVistaPrincipal(fechaInicio, fechaFin time.Time, idGrado int, diasDeshabilitados string) (*models.DatosVistaPrincipal, error) {
	// Obtener estudiantes activos filtrados por grado
	estudiantes, err := s.Repo.ObtenerEstudiantesPorGrado(idGrado)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estudiantes: %v", err)
	}

	// Obtener productos activos
	productos, err := s.Repo.ObtenerProductosActivos()
	if err != nil {
		return nil, fmt.Errorf("error al obtener productos: %v", err)
	}

	// Obtener consumos de la semana
	consumos, err := s.Repo.ObtenerConsumosSemana(fechaInicio, fechaFin)
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
	deudasAnteriores, err := s.Repo.ObtenerDeudasAnterioresBatch(idGrado, fechaInicio)
	if err != nil {
		return nil, fmt.Errorf("error al obtener deudas anteriores: %v", err)
	}

	pagosSemana, err := s.Repo.ObtenerPagosSemanaBatch(idGrado, fechaInicio, fechaFin)
	if err != nil {
		return nil, fmt.Errorf("error al obtener pagos: %v", err)
	}

	// Calcular datos de cada estudiante (sin queries adicionales)
	estudiantesConData := make([]models.EstudianteConDeuda, 0, len(estudiantes))
	for _, est := range estudiantes {
		subTotal := subTotalesPorEstudiante[est.IdEstudiante]
		deudaAnterior := deudasAnteriores[est.IdEstudiante]
		descuento := pagosSemana[est.IdEstudiante]
		total := subTotal + deudaAnterior - descuento

		estudiantesConData = append(estudiantesConData, models.EstudianteConDeuda{
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
	todosDias := utils.ObtenerDiasHabiles(fechaInicio)
	diasConEstado := make([]models.DiaConEstado, 0, len(todosDias))
	diasHabilesFiltrados := make([]time.Time, 0, len(todosDias))

	for _, dia := range todosDias {
		fechaKey := dia.Format("2006-01-02")
		esHabil := !mapaDiasDeshabilitados[fechaKey]

		diasConEstado = append(diasConEstado, models.DiaConEstado{
			Fecha:   dia,
			EsHabil: esHabil,
		})

		if esHabil {
			diasHabilesFiltrados = append(diasHabilesFiltrados, dia)
		}
	}

	return &models.DatosVistaPrincipal{
		Semana:             utils.FormatearSemana(fechaInicio, fechaFin),
		FechaInicio:        fechaInicio,
		FechaFin:           fechaFin,
		DiasConEstado:      diasConEstado,
		DiasHabiles:        diasHabilesFiltrados,
		Productos:          productos,
		EstudiantesConData: estudiantesConData,
		ConsumosPorDia:     consumosPorDia,
		Grados:             utils.ObtenerGradosEstaticos(),
		GradoSeleccionado:  idGrado,
		DiasDeshabilitados: diasDeshabilitados,
	}, nil
}

// RegistrarConsumoDesdeFormulario procesa el registro de un consumo desde el formulario
func (s *Servicio) RegistrarConsumoDesdeFormulario(idEstudiante, idProducto, cantidad int, fecha time.Time) error {
	// Obtener el precio actual del producto
	producto, err := s.Repo.ObtenerProductoPorId(idProducto)
	if err != nil {
		return fmt.Errorf("producto no encontrado: %v", err)
	}

	// Actualizar o insertar el consumo
	return s.Repo.ActualizarConsumo(idEstudiante, idProducto, fecha, cantidad, producto.PrecioUnitario)
}

// RegistrarPagoDesdeFormulario procesa el registro de un pago
func (s *Servicio) RegistrarPagoDesdeFormulario(idEstudiante int, monto float64, fecha time.Time) error {
	if monto <= 0 {
		return fmt.Errorf("el monto debe ser mayor a cero")
	}

	pago := models.Pago{
		IdEstudiante: idEstudiante,
		Monto:        monto,
		FechaPago:    fecha,
	}

	return s.Repo.RegistrarPago(pago)
}
