package controllers

import (
	"kiosco/internal/models"
	"kiosco/internal/utils"
	"kiosco/templates/pages"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Inicio muestra la vista principal con la semana actual
func (m *Controlador) Inicio(w http.ResponseWriter, r *http.Request) {
	fechaInicio, fechaFin := utils.ObtenerSemanaActual()

	if fechaParam := r.URL.Query().Get("fecha"); fechaParam != "" {
		if fecha, err := time.Parse("2006-01-02", fechaParam); err == nil {
			fechaInicio, fechaFin = utils.CalcularSemanaDesdeFecha(fecha)
		}
	}

	idGrado := 1
	if gradoParam := r.URL.Query().Get("grado"); gradoParam != "" {
		if grado, err := strconv.Atoi(gradoParam); err == nil {
			idGrado = grado
		}
	}

	diasDeshabilitados := r.URL.Query().Get("dias_off")

	datos, err := m.servicio.ObtenerDatosVistaPrincipal(fechaInicio, fechaFin, idGrado, diasDeshabilitados)
	if err != nil {
		log.Printf("Error al obtener datos: %v", err)
		http.Error(w, "Error al cargar datos", http.StatusInternalServerError)
		return
	}

	// Si no hay estudiantes en absoluto, redirigir al setup
	if len(datos.EstudiantesConData) == 0 {
		todos, _ := m.servicio.Repo.ObtenerEstudiantesActivos()
		if len(todos) == 0 {
			http.Redirect(w, r, "/setup", http.StatusFound)
			return
		}
	}

	if err := pages.Inicio(*datos).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar inicio: %v", err)
	}
}

// VerConsumoSemanal muestra el comprobante de consumo semanal de un estudiante
func (m *Controlador) VerConsumoSemanal(w http.ResponseWriter, r *http.Request) {
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

	fechaInicio, fechaFin := utils.CalcularSemanaDesdeFecha(fecha)

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

	consumos, err := m.servicio.Repo.ObtenerConsumosSemana(fechaInicio, fechaFin)
	if err != nil {
		http.Error(w, "Error al obtener consumos", http.StatusInternalServerError)
		return
	}

	productos, err := m.servicio.Repo.ObtenerProductosActivos()
	if err != nil {
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}

	productosMap := make(map[int]models.Producto)
	for _, p := range productos {
		productosMap[p.IdProducto] = p
	}

	diasConsumo := make(map[string][]models.ConsumoProducto)
	totalesPorDia := make(map[string]float64)

	for _, c := range consumos {
		if c.IdEstudiante != idEstudiante {
			continue
		}
		fechaKey := c.FechaConsumo.Format("2006-01-02")
		producto := productosMap[c.IdProducto]
		diasConsumo[fechaKey] = append(diasConsumo[fechaKey], models.ConsumoProducto{
			Nombre:   producto.Nombre,
			Cantidad: c.Cantidad,
			Precio:   c.PrecioUnitarioVenta,
			Total:    c.TotalLinea,
		})
		totalesPorDia[fechaKey] += c.TotalLinea
	}

	var consumosPorDia []models.ConsumoDiario
	currentDate := fechaInicio
	subTotal := 0.0

	for !currentDate.After(fechaFin) {
		fechaKey := currentDate.Format("2006-01-02")
		consumoDia := models.ConsumoDiario{
			Fecha:     currentDate,
			Productos: diasConsumo[fechaKey],
			Total:     totalesPorDia[fechaKey],
		}
		if len(consumoDia.Productos) > 0 {
			consumosPorDia = append(consumosPorDia, consumoDia)
			subTotal += consumoDia.Total
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	deudaAnterior, _ := m.servicio.Repo.ObtenerDeudaAnterior(idEstudiante, fechaInicio)
	pagos, _ := m.servicio.Repo.ObtenerPagosSemana(idEstudiante, fechaInicio, fechaFin)

	datos := models.DatosConsumoSemanal{
		IdEstudiante:      idEstudiante,
		NombreEstudiante:  nombreEstudiante,
		FechaInicio:       fechaInicio,
		FechaFin:          fechaFin,
		ConsumosPorDia:    consumosPorDia,
		SubTotal:          subTotal,
		DeudaAnterior:     deudaAnterior,
		Pagos:             pagos,
		Total:             subTotal + deudaAnterior - pagos,
		GradoSeleccionado: idGrado,
	}

	if err := pages.VerConsumoSemanal(datos).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar comprobante: %v", err)
	}
}
