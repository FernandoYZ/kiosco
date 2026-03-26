package controllers

import (
	"fmt"
	"kiosco/internal/models"
	"kiosco/internal/utils"
	"kiosco/templates/pages"
	"log"
	"net/http"
	"strconv"
	"time"
)

// RegistrarPago procesa el formulario de registro de pago
func (m *Controlador) RegistrarPago(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
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
		fecha = time.Now()
	}

	fechaPago := fecha
	if fechaPagoStr := r.FormValue("fecha_pago"); fechaPagoStr != "" {
		if fp, err := time.Parse("2006-01-02", fechaPagoStr); err == nil {
			fechaPago = fp
		}
	}

	if err := m.servicio.RegistrarPagoDesdeFormulario(idEstudiante, monto, fechaPago); err != nil {
		log.Printf("Error al registrar pago: %v", err)
		http.Error(w, "Error al registrar pago: "+err.Error(), http.StatusInternalServerError)
		return
	}

	grado := r.FormValue("grado")
	redirect := r.FormValue("redirect")

	var urlRedireccion string
	if redirect == "editar-pagos" {
		urlRedireccion = fmt.Sprintf("/editar-pagos?id_estudiante=%d&fecha=%s", idEstudiante, fechaStr)
		if grado != "" {
			urlRedireccion += "&grado=" + grado
		}
	} else {
		urlRedireccion = "/?fecha=" + fechaStr
		if grado != "" {
			urlRedireccion += "&grado=" + grado
		}
	}

	http.Redirect(w, r, urlRedireccion, http.StatusSeeOther)
}

// EditarPagos muestra la vista para editar pagos de una semana
func (m *Controlador) EditarPagos(w http.ResponseWriter, r *http.Request) {
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

	pagos, err := m.servicio.Repo.ObtenerPagosSemanaDetalle(idEstudiante, fechaInicio, fechaFin)
	if err != nil {
		http.Error(w, "Error al obtener pagos", http.StatusInternalServerError)
		return
	}

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

	if err := pages.EditarPagos(datos).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar editar_pagos: %v", err)
	}
}

// EliminarPago elimina un pago; soporta respuesta HTMX o redirect normal
func (m *Controlador) EliminarPago(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	idPago, err := strconv.Atoi(r.FormValue("id_pago"))
	if err != nil {
		http.Error(w, "ID de pago inválido", http.StatusBadRequest)
		return
	}

	if err := m.servicio.Repo.EliminarPago(idPago); err != nil {
		log.Printf("Error al eliminar pago: %v", err)
		http.Error(w, "Error al eliminar pago", http.StatusInternalServerError)
		return
	}

	// HTMX: retornar respuesta vacía para que el elemento desaparezca del DOM
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	idEstudiante := r.FormValue("id_estudiante")
	fechaStr := r.FormValue("fecha")
	grado := r.FormValue("grado")

	urlRedireccion := fmt.Sprintf("/editar-pagos?id_estudiante=%s&fecha=%s", idEstudiante, fechaStr)
	if grado != "" {
		urlRedireccion += "&grado=" + grado
	}
	http.Redirect(w, r, urlRedireccion, http.StatusSeeOther)
}
