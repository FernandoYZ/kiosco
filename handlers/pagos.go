package handlers

import (
	"fmt"
	"kiosco/models"
	"kiosco/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ManejadorRegistrarPago procesa el formulario de registro de pago
func (m *Manejador) ManejadorRegistrarPago(w http.ResponseWriter, r *http.Request) {
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

	monto, err := strconv.ParseFloat(r.FormValue("monto"), 64)
	if err != nil {
		http.Error(w, "Monto inválido", http.StatusBadRequest)
		return
	}

	fechaStr := r.FormValue("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		fecha = time.Now() // Si no hay fecha, usar la actual
	}

	// Permitir especificar fecha_pago diferente (para la vista de editar pagos)
	fechaPagoStr := r.FormValue("fecha_pago")
	fechaPago := fecha // Por defecto usar la fecha de la semana
	if fechaPagoStr != "" {
		if fp, err := time.Parse("2006-01-02", fechaPagoStr); err == nil {
			fechaPago = fp
		}
	}

	// Registrar el pago con la fecha especificada
	err = m.servicio.RegistrarPagoDesdeFormulario(idEstudiante, monto, fechaPago)
	if err != nil {
		log.Printf("Error al registrar pago: %v", err)
		http.Error(w, "Error al registrar pago: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Obtener el grado para mantenerlo en la redirección
	grado := r.FormValue("grado")
	redirect := r.FormValue("redirect")

	var urlRedireccion string
	if redirect == "editar-pagos" {
		// Redirigir a la vista de editar pagos
		urlRedireccion = fmt.Sprintf("/editar-pagos?id_estudiante=%d&fecha=%s", idEstudiante, fechaStr)
		if grado != "" {
			urlRedireccion += "&grado=" + grado
		}
	} else {
		// Redirigir a la vista principal
		urlRedireccion = "/?fecha=" + fechaStr
		if grado != "" {
			urlRedireccion += "&grado=" + grado
		}
	}

	// Redireccionar de vuelta
	http.Redirect(w, r, urlRedireccion, http.StatusSeeOther)
}

// ManejadorEditarPagos muestra la vista para editar pagos de una semana
func (m *Manejador) ManejadorEditarPagos(w http.ResponseWriter, r *http.Request) {
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

	// Obtener pagos de la semana
	pagos, err := m.servicio.Repo.ObtenerPagosSemanaDetalle(idEstudiante, fechaInicio, fechaFin)
	if err != nil {
		http.Error(w, "Error al obtener pagos", http.StatusInternalServerError)
		return
	}

	// Calcular total de pagos
	totalPagos := 0.0
	for _, p := range pagos {
		totalPagos += p.Monto
	}

	datos := models.DatosEditarPagos{
		IdEstudiante:      idEstudiante,
		NombreEstudiante:  nombreEstudiante,
		FechaInicio:       fechaInicio,
		FechaFin:          fechaFin,
		Pagos:             pagos,
		TotalPagos:        totalPagos,
		GradoSeleccionado: idGrado,
	}

	err = m.templates.ExecuteTemplate(w, "editar_pagos.tmpl", datos)
	if err != nil {
		log.Printf("Error al renderizar template: %v", err)
	}
}

// ManejadorEliminarPago elimina un pago específico
func (m *Manejador) ManejadorEliminarPago(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	idPago, err := strconv.Atoi(r.FormValue("id_pago"))
	if err != nil {
		http.Error(w, "ID de pago inválido", http.StatusBadRequest)
		return
	}

	// Eliminar el pago
	err = m.servicio.Repo.EliminarPago(idPago)
	if err != nil {
		log.Printf("Error al eliminar pago: %v", err)
		http.Error(w, "Error al eliminar pago", http.StatusInternalServerError)
		return
	}

	// Redirigir de vuelta a la vista de editar pagos
	idEstudiante := r.FormValue("id_estudiante")
	fechaStr := r.FormValue("fecha")
	grado := r.FormValue("grado")

	urlRedireccion := fmt.Sprintf("/editar-pagos?id_estudiante=%s&fecha=%s", idEstudiante, fechaStr)
	if grado != "" {
		urlRedireccion += "&grado=" + grado
	}

	http.Redirect(w, r, urlRedireccion, http.StatusSeeOther)
}
