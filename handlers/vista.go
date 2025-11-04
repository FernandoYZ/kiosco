package handlers

import (
	"kiosco/models"
	"kiosco/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ManejadorInicio muestra la vista principal con la semana actual
func (m *Manejador) ManejadorInicio(w http.ResponseWriter, r *http.Request) {
	fechaInicio, fechaFin := utils.ObtenerSemanaActual()

	// Verificar si hay parámetros de fecha en la URL
	if fechaParam := r.URL.Query().Get("fecha"); fechaParam != "" {
		if fecha, err := time.Parse("2006-01-02", fechaParam); err == nil {
			fechaInicio, fechaFin = utils.CalcularSemanaDesdeFecha(fecha)
		}
	}

	// Verificar si hay filtro de grado (por defecto primer grado = 1)
	idGrado := 1
	if gradoParam := r.URL.Query().Get("grado"); gradoParam != "" {
		if grado, err := strconv.Atoi(gradoParam); err == nil {
			idGrado = grado
		}
	}

	// Obtener días deshabilitados del parámetro URL
	diasDeshabilitados := r.URL.Query().Get("dias_off")

	datos, err := m.servicio.ObtenerDatosVistaPrincipal(fechaInicio, fechaFin, idGrado, diasDeshabilitados)
	if err != nil {
		log.Printf("Error al obtener datos: %v", err)
		http.Error(w, "Error al cargar datos", http.StatusInternalServerError)
		return
	}

	err = m.templates.ExecuteTemplate(w, "index.tmpl", datos)
	if err != nil {
		log.Printf("Error al renderizar template: %v", err)
	}
}

// ManejadorVerConsumoSemanal muestra el comprobante de consumo semanal de un estudiante
func (m *Manejador) ManejadorVerConsumoSemanal(w http.ResponseWriter, r *http.Request) {
	idEstudiante, err := strconv.Atoi(r.URL.Query().Get("id_estudiante"))
	if err != nil {
		http.Error(w, "ID de estudiante inválido", http.StatusBadRequest)
		return
	}

	fechaStr := r.URL.Query().Get("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}

	idGrado := 0
	if gradoParam := r.URL.Query().Get("grado"); gradoParam != "" {
		if grado, err := strconv.Atoi(gradoParam); err == nil {
			idGrado = grado
		}
	}

	// Calcular la semana
	fechaInicio, fechaFin := utils.CalcularSemanaDesdeFecha(fecha)

	// Obtener datos del estudiante
	estudiantes, err := m.servicio.Repo.ObtenerEstudiantesPorGrado(0)
	if err != nil {
		http.Error(w, "Error al obtener estudiante", http.StatusInternalServerError)
		return
	}

	var nombreEstudiante string
	for _, est := range estudiantes {
		if est.IdEstudiante == idEstudiante {
			nombreEstudiante = est.Apellidos + ", " + est.Nombres
			break
		}
	}

	// Obtener consumos de la semana
	consumos, err := m.servicio.Repo.ObtenerConsumosSemana(fechaInicio, fechaFin)
	if err != nil {
		http.Error(w, "Error al obtener consumos", http.StatusInternalServerError)
		return
	}

	// Obtener productos
	productos, err := m.servicio.Repo.ObtenerProductosActivos()
	if err != nil {
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}

	// Crear mapa de productos para búsqueda rápida
	productosMap := make(map[int]models.Producto)
	for _, p := range productos {
		productosMap[p.IdProducto] = p
	}

	// Organizar consumos por día
	diasConsumo := make(map[string][]models.ConsumoProducto)
	totalesPorDia := make(map[string]float64)

	for _, c := range consumos {
		if c.IdEstudiante != idEstudiante {
			continue
		}

		fechaKey := c.FechaConsumo.Format("2006-01-02")
		producto := productosMap[c.IdProducto]

		consumoProducto := models.ConsumoProducto{
			Nombre:   producto.Nombre,
			Cantidad: c.Cantidad,
			Precio:   c.PrecioUnitarioVenta,
			Total:    c.TotalLinea,
		}

		diasConsumo[fechaKey] = append(diasConsumo[fechaKey], consumoProducto)
		totalesPorDia[fechaKey] += c.TotalLinea
	}

	// Crear lista de consumos diarios en orden
	var consumosPorDia []models.ConsumoDiario
	currentDate := fechaInicio
	subTotal := 0.0

	for currentDate.Before(fechaFin) || currentDate.Equal(fechaFin) {
		fechaKey := currentDate.Format("2006-01-02")

		consumoDia := models.ConsumoDiario{
			Fecha:     currentDate,
			Productos: diasConsumo[fechaKey],
			Total:     totalesPorDia[fechaKey],
		}

		// Solo agregar días con consumo
		if len(consumoDia.Productos) > 0 {
			consumosPorDia = append(consumosPorDia, consumoDia)
			subTotal += consumoDia.Total
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Obtener deuda anterior y pagos
	deudaAnterior, _ := m.servicio.Repo.ObtenerDeudaAnterior(idEstudiante, fechaInicio)
	pagos, _ := m.servicio.Repo.ObtenerPagosSemana(idEstudiante, fechaInicio, fechaFin)

	total := subTotal + deudaAnterior - pagos

	datos := models.DatosConsumoSemanal{
		IdEstudiante:      idEstudiante,
		NombreEstudiante:  nombreEstudiante,
		FechaInicio:       fechaInicio,
		FechaFin:          fechaFin,
		ConsumosPorDia:    consumosPorDia,
		SubTotal:          subTotal,
		DeudaAnterior:     deudaAnterior,
		Pagos:             pagos,
		Total:             total,
		GradoSeleccionado: idGrado,
	}

	err = m.templates.ExecuteTemplate(w, "ver_consumo_semanal.tmpl", datos)
	if err != nil {
		log.Printf("Error al renderizar template: %v", err)
	}
}
