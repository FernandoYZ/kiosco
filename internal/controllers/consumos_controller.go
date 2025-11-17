package controllers

import (
	"fmt"
	"kiosco/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ManejadorRegistrarConsumo procesa el formulario de registro de consumo
func (m *Controlador) RegistrarConsumo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear datos del formulario (soporta form-urlencoded y multipart)
	err := r.ParseForm()
	if err != nil {
		// Intentar con multipart si falla
		err = r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			log.Printf("Error al parsear formulario: %v", err)
			http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
			return
		}
	}

	idEstudianteStr := r.FormValue("id_estudiante")
	idEstudiante, err := strconv.Atoi(idEstudianteStr)
	if err != nil {
		log.Printf("Error al convertir id_estudiante '%s': %v", idEstudianteStr, err)
		http.Error(w, "ID de estudiante inválido", http.StatusBadRequest)
		return
	}

	idProducto, err := strconv.Atoi(r.FormValue("id_producto"))
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}

	cantidad, err := strconv.Atoi(r.FormValue("cantidad"))
	if err != nil {
		http.Error(w, "Cantidad inválida", http.StatusBadRequest)
		return
	}

	fechaStr := r.FormValue("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}

	// Registrar el consumo
	err = m.servicio.RegistrarConsumoDesdeFormulario(idEstudiante, idProducto, cantidad, fecha)
	if err != nil {
		log.Printf("Error al registrar consumo: %v", err)
		http.Error(w, "Error al registrar consumo", http.StatusInternalServerError)
		return
	}

	// Obtener el grado para mantenerlo en la redirección
	grado := r.FormValue("grado")
	urlRedireccion := "/?fecha=" + fechaStr
	if grado != "" {
		urlRedireccion += "&grado=" + grado
	}

	// Redireccionar de vuelta a la vista principal
	http.Redirect(w, r, urlRedireccion, http.StatusSeeOther)
}

// EditarConsumos muestra la vista para editar consumos de un día
func (m *Controlador) EditarConsumos(w http.ResponseWriter, r *http.Request) {
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

	// Obtener datos del estudiante
	estudiantes, err := m.servicio.Repo.ObtenerEstudiantesPorGrado(0) // Todos para buscar
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

	// Obtener productos
	productos, err := m.servicio.Repo.ObtenerProductosActivos()
	if err != nil {
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}

	// Obtener consumos existentes
	fechaInicio := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, fecha.Location())
	fechaFin := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 23, 59, 59, 0, fecha.Location())
	consumos, err := m.servicio.Repo.ObtenerConsumosSemana(fechaInicio, fechaFin)
	if err != nil {
		http.Error(w, "Error al obtener consumos", http.StatusInternalServerError)
		return
	}

	// Crear mapa de consumos
	consumosPorDia := make(map[int]map[string]map[int]int)
	for _, c := range consumos {
		if consumosPorDia[c.IdEstudiante] == nil {
			consumosPorDia[c.IdEstudiante] = make(map[string]map[int]int)
		}
		fechaKey := c.FechaConsumo.Format("2006-01-02")
		if consumosPorDia[c.IdEstudiante][fechaKey] == nil {
			consumosPorDia[c.IdEstudiante][fechaKey] = make(map[int]int)
		}
		consumosPorDia[c.IdEstudiante][fechaKey][c.IdProducto] = c.Cantidad
	}

	datos := models.DatosEditarConsumos{
		IdEstudiante:      idEstudiante,
		NombreEstudiante:  nombreEstudiante,
		Fecha:             fecha,
		Productos:         productos,
		Consumos:          consumosPorDia,
		GradoSeleccionado: idGrado,
	}

	m.renderizar(w, "editar_consumos.tmpl", datos)
}

// ManejadorGuardarConsumosDia guarda todos los consumos de un día
func (m *Controlador) GuardarConsumosDia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	idEstudiante, err := strconv.Atoi(r.FormValue("id_estudiante"))
	if err != nil {
		http.Error(w, "ID de estudiante inválido", http.StatusBadRequest)
		return
	}

	fechaStr := r.FormValue("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}

	grado := r.FormValue("grado")

	// Obtener productos para procesar
	productos, err := m.servicio.Repo.ObtenerProductosActivos()
	if err != nil {
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}

	// Procesar cada producto
	for _, producto := range productos {
		cantidadStr := r.FormValue(fmt.Sprintf("cantidad_%d", producto.IdProducto))
		cantidad, err := strconv.Atoi(cantidadStr)
		if err != nil {
			cantidad = 0
		}

		// Actualizar consumo
		err = m.servicio.RegistrarConsumoDesdeFormulario(idEstudiante, producto.IdProducto, cantidad, fecha)
		if err != nil {
			log.Printf("Error al registrar consumo producto %d: %v", producto.IdProducto, err)
		}
	}

	// Redireccionar
	urlRedireccion := "/?fecha=" + fechaStr
	if grado != "" {
		urlRedireccion += "&grado=" + grado
	}

	http.Redirect(w, r, urlRedireccion, http.StatusSeeOther)
}
