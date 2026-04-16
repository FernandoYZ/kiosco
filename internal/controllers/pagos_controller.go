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

	// Obtener consumos de la semana
	consumos, err := m.servicio.Repo.ObtenerConsumosSemana(fechaInicio, fechaFin)
	subTotal := 0.0
	if err == nil {
		for _, c := range consumos {
			if c.IdEstudiante == idEstudiante {
				subTotal += c.TotalLinea
			}
		}
	}

	// Obtener deuda anterior (acumulada de semanas anteriores)
	deudaAnterior, err := m.servicio.Repo.ObtenerDeudaAnterior(idEstudiante, fechaInicio)
	if err != nil {
		deudaAnterior = 0.0
	}

	// Deuda actual = (Consumos semana + Deuda anterior) - Pagos realizados
	// Nota: Puede ser negativo si hay saldo a favor (cliente pagó más de lo debido)
	deudaActual := (subTotal + deudaAnterior) - totalPagos

	datos := models.DatosEditarPagos{
		IdEstudiante:      idEstudiante,
		NombreEstudiante:  nombreEstudiante,
		FechaInicio:       fechaInicio,
		FechaFin:          fechaFin,
		Pagos:             pagos,
		TotalPagos:        totalPagos,
		DeudaActual:       deudaActual,
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

	idEstudiante := r.FormValue("id_estudiante")
	fechaStr := r.FormValue("fecha")
	grado := r.FormValue("grado")

	if err := m.servicio.Repo.EliminarPago(idPago); err != nil {
		log.Printf("Error al eliminar pago: %v", err)
		http.Error(w, "Error al eliminar pago", http.StatusInternalServerError)
		return
	}

	// HTMX: retornar respuesta vacía para eliminar el pago + actualizar saldo con OOB (Out of Band)
	if r.Header.Get("HX-Request") == "true" {
		// Re-calcular datos actualizados después de eliminar el pago
		idEstudianteInt, _ := strconv.Atoi(idEstudiante)
		fecha, _ := time.Parse("2006-01-02", fechaStr)
		fechaInicio, fechaFin := utils.CalcularSemanaDesdeFecha(fecha)

		pagos, _ := m.servicio.Repo.ObtenerPagosSemanaDetalle(idEstudianteInt, fechaInicio, fechaFin)
		totalPagos := 0.0
		for _, p := range pagos {
			totalPagos += p.Monto
		}

		consumos, _ := m.servicio.Repo.ObtenerConsumosSemana(fechaInicio, fechaFin)
		subTotal := 0.0
		for _, c := range consumos {
			if c.IdEstudiante == idEstudianteInt {
				subTotal += c.TotalLinea
			}
		}

		deudaAnterior, _ := m.servicio.Repo.ObtenerDeudaAnterior(idEstudianteInt, fechaInicio)
		deudaActual := (subTotal + deudaAnterior) - totalPagos

		// Retornar respuesta vacía para eliminar el pago + OOB para actualizar saldo
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// El elemento vacío elimina el pago del DOM
		fmt.Fprint(w, "")

		// Out-of-Band Swap para actualizar el saldo sin cambiar el target
		fmt.Fprintf(w, `<div id="saldo-info" hx-swap-oob="outerHTML" class="mt-6 p-5 bg-white rounded-3xl border border-gray-200 hidden lg:block">
			<div class="flex items-center justify-between">
				<span class="text-[15px] font-medium text-gray-500">Saldo actual</span>`)

		if deudaActual > 0 {
			fmt.Fprintf(w, `<span class="text-[20px] font-black text-[#FF3B30] tabular-nums">S/ %s</span>`, utils.FormatearMoneda(deudaActual))
		} else {
			fmt.Fprintf(w, `<span class="text-[20px] font-black text-[#34C759] tabular-nums">S/ %s</span>`, utils.FormatearMoneda(deudaActual))
		}

		fmt.Fprint(w, `</div>
		</div>`)
		return
	}

	urlRedireccion := fmt.Sprintf("/editar-pagos?id_estudiante=%s&fecha=%s", idEstudiante, fechaStr)
	if grado != "" {
		urlRedireccion += "&grado=" + grado
	}
	http.Redirect(w, r, urlRedireccion, http.StatusSeeOther)
}
